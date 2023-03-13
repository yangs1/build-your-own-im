package models

import (
	"context"
	"fmt"
	"gim/global"
	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"gopkg.in/fatih/set.v0"
	"net/http"
	"strconv"
	"sync"
)

type Message struct {
	Model
	FormId   int64  `json:"userId"`   //信息发送者
	TargetId int64  `json:"targetId"` //信息接收者
	Type     int    //聊天类型：群聊 私聊 广播
	Media    int    //信息类型：文字 图片 音频
	Content  string //消息内容
	Pic      string `json:"url"` //图片相关
	Url      string //文件相关
	Desc     string //文件描述
	Amount   int    //其他数据大小
}

//MsgTableName 生成指定数据表名
func (m *Message) MsgTableName() string {
	return "message"
}

//Node 构造连接
type Node struct {
	Conn      *websocket.Conn //连接
	Addr      string          //客户端地址
	DataQueue chan []byte     //消息
	GroupSets set.Interface   //好友 / 群
}

//映射关系
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

//读写锁，绑定node时需要线程安全
var rwLocker sync.RWMutex

//Chat	需要 ：发送者ID ，接受者ID ，消息类型，发送的内容，发送类型
func Chat(w http.ResponseWriter, r *http.Request) {
	//1.  获取参数 并 检验 token 等合法性
	query := r.URL.Query()
	fmt.Println("handle:", query)
	Id := query.Get("userId")
	//token := query.Get("token")

	userId, err := strconv.ParseInt(Id, 10, 64)
	if err != nil {
		zap.S().Info("类型转换失败", err)
		return
	}

	//升级为socket
	var isvalida = true
	conn, err := (&websocket.Upgrader{
		//token 校验
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	//获取socket连接,构造消息节点
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	//用户关系

	//将userId和Node绑定
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()

	fmt.Println("uid", userId)

	//发送接收消息
	//发送消息
	go sendProc(node)
	//接收消息
	go recProc(node)
	sendMsg(userId, []byte("欢迎进入聊天系统"))
}

func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				zap.S().Info("写入消息失败", err)
				return
			}
			fmt.Println("数据发送socket成功")
		}

	}
}

//recProc 从websocket中将消息体拿出，然后进行解析，再进行信息类型判断， 最后将消息发送至目的用户的node中
func recProc(node *Node) {
	for {
		//获取信息
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			zap.S().Info("读取消息失败", err)
			return
		}

		// 处理发生逻辑
		dispatch(data)

	}
}

//dispatch 解析消息，聊天类型判断
func dispatch(data []byte) {
	//解析消息
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		zap.S().Info("消息解析失败", err)
		return
	}

	//判断消息类型
	switch msg.Type {
	case 1: //私聊
		sendMsg(msg.TargetId, data)
	case 2: //群发
		sendGroupMsg(uint(msg.FormId), uint(msg.TargetId), data)
	}
}

//sendMs 向用户发送消息
func sendMsg(id int64, msg []byte) {
	rwLocker.Lock()
	node, ok := clientMap[id]
	rwLocker.Unlock()

	if !ok {
		zap.S().Info("userID没有对应的node")
		return
	}

	zap.S().Info("targetID:", id, "node:", node)
	if ok {
		node.DataQueue <- msg
	}
}

//sendGroupMsg 群发
func sendGroupMsg(formId, target uint, data []byte) (int, error) {
	//群发的逻辑：1获取到群里所有用户，然后向除开自己的每一位用户发送消息
	userIDs, err := FindUsers(target)
	if err != nil {
		return -1, err
	}

	for _, userId := range *userIDs {
		if formId != userId { //不能给当前发送消息的成员进行转发
			//调用单聊的函数，群聊变成了多次单聊
			sendMsgAndSave(int64(userId), data)
		}
	}
	return 0, nil
}

//sendMsgTest 发送消息 并存储聊天记录到redis
func sendMsgAndSave(userId int64, msg []byte) {

	rwLocker.RLock()              //保证线程安全，上锁
	node, ok := clientMap[userId] //对方是否在线
	rwLocker.RUnlock()            //解锁

	jsonMsg := Message{}
	json.Unmarshal(msg, &jsonMsg)
	ctx := context.Background()
	targetIdStr := strconv.Itoa(int(userId))
	userIdStr := strconv.Itoa(int(jsonMsg.FormId))

	if ok {
		//如果当前用户在线，将消息转发到当前用户的websocket连接中，然后进行存储
		node.DataQueue <- msg
	}

	//userIdStr和targetIdStr进行拼接唯一key
	var key string
	if userId > jsonMsg.FormId {
		key = "msg_" + userIdStr + "_" + targetIdStr
	} else {
		key = "msg_" + targetIdStr + "_" + userIdStr
	}

	//创建记录
	res, err := global.RedisDB.ZRevRange(ctx, key, 0, -1).Result()
	if err != nil {
		fmt.Println(err)
		return
	}

	//将聊天记录写入redis缓存中
	score := float64(cap(res)) + 1
	ress, e := global.RedisDB.ZAdd(ctx, key, &redis.Z{score, msg}).Result() //jsonMsg
	//res, e := utils.Red.Do(ctx, "zadd", key, 1, jsonMsg).Result() //备用 后续拓展 记录完整msg
	if e != nil {
		fmt.Println(e)
		return
	}
	fmt.Println(ress)
}
