package router

import (
	"gim/middlewear"
	"gim/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Router() *gin.Engine {
	//初始化路由
	router := gin.Default()

	//v1版本
	v1 := router.Group("v1")

	//用户模块，后续有个用户的api就放置其中
	user := v1.Group("user").Use(middlewear.JWY())
	{
		user.GET("/list", service.List)
		user.POST("/new", service.NewUser)
		user.DELETE("/delete", service.DeleteUser)
		user.POST("/updata", service.UpdataUser)
	}

	v1.POST("/login_pw", service.LoginByNameAndPassWord)

	//好友关系
	relation := v1.Group("relation").Use(middlewear.JWY())
	{
		relation.POST("/list", service.FriendList)
		relation.POST("/add", service.AddFriendByName)
	}

	router.GET("/:name/:id", func(c *gin.Context) {
		//使用porsen对数据进行解组
		var porsen Porsen
		if err := c.ShouldBindUri(&porsen); err != nil {
			c.Status(404)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"name": porsen.Name,
			"id":   porsen.ID,
		})
	})
	return router
}

//结构体声明，并做一些约束
type Porsen struct {
	ID   int    `uri:"id" binding:"required"`   //uri指在client中的名字为id，binding:"required指必填
	Name string `uri:"name" binding:"required"` //同理
}
