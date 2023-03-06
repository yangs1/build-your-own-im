package test

import (
	"encoding/json"
	"gim/dao"
	"gim/global"
	"gim/initialize"
	"gim/models"
	"log"
	"testing"
)

func TestUserDB(t *testing.T) {
	initialize.InitDB()
	err := global.DB.AutoMigrate(&models.UserBasic{})

	if err != nil {
		panic(err)
	}
}

func TestAddUser(t *testing.T) {
	initialize.InitDB()

	str := `[ 
  {
        "ID": 8,
        "CreatedAt": "2022-12-22T19:17:16.365+08:00",
        "UpdatedAt": "2022-12-22T19:17:16.365+08:00",
        "DeletedAt": null,
        "Name": "ice_moss",
        "PassWord": "d41d8cd98f00b204e9800998ecf8427e$1298498081",
        "Avatar": "",
        "Gender": "male",
        "Phone": "",
        "Email": "",
        "Identity": "",
        "ClientIp": "",
        "ClientPort": "",
        "Salt": "1298498081",
        "LoginTime": "2022-12-22T19:17:16.363+08:00",
        "HeartBeatTime": "2022-12-22T19:17:16.363+08:00",
        "LoginOutTime": "2022-12-22T19:17:16.363+08:00",
        "IsLoginOut": false,
        "DeviceInfo": ""
    },
    {
        "ID": 9,
        "CreatedAt": "2022-12-22T19:30:34.893+08:00",
        "UpdatedAt": "2022-12-22T19:30:34.893+08:00",
        "DeletedAt": null,
        "Name": "ice_moss1",
        "PassWord": "d41d8cd98f00b204e9800998ecf8427e$1298498081",
        "Avatar": "",
        "Gender": "male",
        "Phone": "",
        "Email": "",
        "Identity": "",
        "ClientIp": "",
        "ClientPort": "",
        "Salt": "1298498081",
        "LoginTime": "2022-12-22T19:30:34.892+08:00",
        "HeartBeatTime": "2022-12-22T19:30:34.892+08:00",
        "LoginOutTime": "2022-12-22T19:30:34.892+08:00",
        "IsLoginOut": false,
        "DeviceInfo": ""
    },
    {
        "ID": 10,
        "CreatedAt": "2022-12-22T19:37:19.508+08:00",
        "UpdatedAt": "2022-12-24T16:38:56.717+08:00",
        "DeletedAt": null,
        "Name": "ice_moss2",
        "PassWord": "0192023a7bbd73250516f069df18b500$1298498081",
        "Avatar": "https://mxshopfiles.oss-cn-shanghai.aliyuncs.com/work/103800kbdgbv2zdv1vnnrd.jpeg",
        "Gender": "male",
        "Phone": "",
        "Email": "",
        "Identity": "9fce97499eea554562d27d086da558e3",
        "ClientIp": "",
        "ClientPort": "",
        "Salt": "1298498081",
        "LoginTime": "2022-12-22T19:37:19.507+08:00",
        "HeartBeatTime": "2022-12-22T19:37:19.507+08:00",
        "LoginOutTime": "2022-12-22T19:37:19.507+08:00",
        "IsLoginOut": false,
        "DeviceInfo": ""
    },
    {
        "ID": 11,
        "CreatedAt": "2022-12-24T16:51:53.418+08:00",
        "UpdatedAt": "2022-12-24T18:26:06.611+08:00",
        "DeletedAt": null,
        "Name": "ice_moss4",
        "PassWord": "0192023a7bbd73250516f069df18b500$1298498081",
        "Avatar": "",
        "Gender": "male",
        "Phone": "",
        "Email": "",
        "Identity": "5993bcfc7b16b8a84e22aefc6b42a528",
        "ClientIp": "",
        "ClientPort": "",
        "Salt": "1298498081",
        "LoginTime": "2022-12-24T16:51:53.417+08:00",
        "HeartBeatTime": "2022-12-24T16:51:53.417+08:00",
        "LoginOutTime": "2022-12-24T16:51:53.417+08:00",
        "IsLoginOut": false,
        "DeviceInfo": ""
    }
 ]`
	var list []*models.UserBasic

	json.Unmarshal([]byte(str), &list)

	log.Printf("slice s: %v(MISSING)\n", list)
	for k, v := range list {
		dao.CreateUser(v)
		log.Printf("v%d : %v\n", k, *v)
	}

}

func TestAddRelation(t *testing.T) {

	initialize.InitDB()
	global.DB.AutoMigrate(&models.Relation{})

	str := `{
        "ID": 11,
        "CreatedAt": "2022-12-24T16:51:53.418+08:00",
        "UpdatedAt": "2022-12-24T18:26:06.611+08:00",
        "DeletedAt": null,
        "OwnerId": 8,
        "TargetID": 14,
		"Type" : 1,
		"Desc" : "-----"
    }`

	relation := models.Relation{}
	json.Unmarshal([]byte(str), &relation)

	log.Printf("%v\n", relation)

	dao.AddFriend(8, 14)
	dao.AddFriend(8, 12)
}
