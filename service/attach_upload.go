package service

import (
	"fmt"
	"gim/common"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

//Image 图片语音上传并返回url
func Image(ctx *gin.Context) {
	w := ctx.Writer
	req := ctx.Request
	//获取文件
	srcFile, head, err := req.FormFile("file")
	if err != nil {
		common.RespFail(w, err.Error())
		return
	}

	//检查文件后缀
	suffix := ".png"
	ofilName := head.Filename
	tem := strings.Split(ofilName, ".")
	if len(tem) > 1 {
		suffix = "." + tem[len(tem)-1]
	}

	//保存文件
	fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)
	dstFile, err := os.Create("./asset/upload/" + fileName)
	if err != nil {
		common.RespFail(w, err.Error())
		return
	}
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		common.RespFail(w, err.Error())
	}
	url := "./asset/upload/" + fileName
	common.RespOK(w, url, "发送成功")
}
