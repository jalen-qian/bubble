package router

import (
	"github.com/bubble/controller"
	"github.com/gin-gonic/gin"
)

func SetupRouters() *gin.Engine {
	eng := gin.Default()
	//加载静态文件
	eng.Static("/static", "./static")
	eng.Delims("[{", "}]").LoadHTMLGlob("./templates/*")

	eng.GET("/", controller.IndexHandler)
	//通过路由组来管理
	v1Group := eng.Group("v1")
	{
		//获取所有的清单
		v1Group.GET("/todo", controller.GetTodoList)
		//创建清单
		v1Group.POST("/todo", controller.CreateATodo)
		//改变清单状态
		v1Group.PUT("/todo/:id", controller.UpdateATodo)
		//删除清单
		v1Group.DELETE("/todo/:id", controller.DeleteATodo)
	}
	return eng
}
