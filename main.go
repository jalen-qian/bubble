package main

import (
	"github.com/bubble/dao"
	"github.com/bubble/models"
	"github.com/bubble/router"
	"log"
)

func main() {
	//创建数据库
	// sql: CREATE DATABASE bubble
	//连接数据库
	err := dao.InitDb()
	if err != nil {
		log.Fatalf("connect to database failed,err:%v\n", err)
	}
	//数据库表与模型对应
	_ = dao.DB.AutoMigrate(&models.Todo{})

	//启动服务
	err = router.SetupRouters().Run()
	if err != nil {
		log.Fatalf("gin run failed,err:%v\n", err)
	}
}
