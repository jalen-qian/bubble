package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

var db *gorm.DB

func initDb() (err error) {
	dsn := "root:123456@tcp(192.168.86.128)/bubble?charset=utf8mb4&parseTime=True&loc=Local"
	mdb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	db = mdb.Debug()
	return nil
}

type Todo struct {
	ID uint `json:"id" gorm:"primarykey"`
	gorm.Model
	Title  string `json:"title" gorm:"type:varchar(100);"`
	Status bool   `json:"status"`
}

type result struct {
	Code    uint        `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func main() {
	//创建数据库
	// sql: CREATE DATABASE bubble
	//连接数据库
	err := initDb()
	if err != nil {
		log.Fatalf("connect to database failed,err:%v\n", err)
	}
	//数据库表与模型对应
	_ = db.AutoMigrate(&Todo{})

	eng := gin.Default()
	//加载静态文件
	eng.Static("/static", "./static")
	eng.Delims("[{", "}]").LoadHTMLGlob("./templates/*")
	eng.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	//通过路由组来管理
	v1Group := eng.Group("v1")
	{
		//获取所有的清单
		v1Group.GET("/todo", func(c *gin.Context) {
			var todoList []Todo
			if err := db.Find(&todoList).Error; err != nil {
				//有错误，返回错误信息
				c.JSON(http.StatusOK, result{
					Code:    1000,
					Message: "获取待办列表失败，err:" + err.Error(),
					Data:    todoList,
				})
			} else {
				//无错误，返回正常数据
				c.JSON(http.StatusOK, result{
					Code:    1000,
					Message: "获取待办列表成功",
					Data:    todoList,
				})
			}

		})
		//获取单个清单
		v1Group.GET("/todo/:id", func(c *gin.Context) {

		})
		//创建清单
		v1Group.POST("/todo", func(c *gin.Context) {
			var todo Todo
			err := c.BindJSON(&todo)
			if err != nil {
				c.JSON(http.StatusBadRequest, result{Code: 1001, Message: "数据不合法，请传JSON格式数据"})
				return
			}
			//将数据保存到数据库
			if err := db.Create(&todo).Error; err != nil {
				c.JSON(http.StatusOK, result{Code: 1001, Message: err.Error()})
			} else {
				c.JSON(http.StatusOK, result{Code: 1000, Message: "添加待办事项成功"})
			}
		})
		//改变清单状态
		v1Group.PUT("/todo/:id", func(c *gin.Context) {
			//获取路由参数中的id
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusOK, result{
					Code:    1001,
					Message: "参数不正确",
				})
				return
			}
			//参数正确，判断参数是否大于0，如果参数为0，也不合法
			if id <= 0 {
				c.JSON(http.StatusOK, result{
					Code:    1001,
					Message: "参数不能小于0",
				})
				return
			}
			var todo = Todo{ID: uint(id),}
			//找到这个条记录
			if err := db.First(&todo).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					c.JSON(http.StatusOK, result{
						Code:    1001,
						Message: "未找到这条记录",
					})
				} else {
					c.JSON(http.StatusOK, result{
						Code:    1001,
						Message: "数据访问异常：" + err.Error(),
					})
				}
				return
			} else {
				log.Println(todo)
				//更新状态
				todo.Status = !todo.Status
				/*if err := db.Save(todo).Error;err != nil {*/

				//下面这种写法更好，上面是所有字段都更新，这个是只更新status这个字段
				if err := db.Model(&todo).Update("status", todo.Status).Error; err != nil {

					c.JSON(http.StatusOK, result{
						Code:    1001,
						Message: "设置失败！",
					})
				} else {
					statusName := ""
					if todo.Status {
						statusName = "已完成"
					} else {
						statusName = "未完成"
					}
					c.JSON(http.StatusOK, result{
						Code:    1000,
						Message: "事项【" + todo.Title + "】设置为" + statusName,
					})
				}
			}
		})
		//删除清单
		v1Group.DELETE("/todo/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusOK, result{
					Code:    1001,
					Message: "参数不正确",
				})
				return
			}
			//参数正确，判断参数是否大于0，如果参数为0，也不合法
			if id <= 0 {
				c.JSON(http.StatusOK, result{
					Code:    1001,
					Message: "参数不能小于0",
				})
				return
			}
			if err := db.Delete(&Todo{}, id).Error; err != nil {
				c.JSON(http.StatusOK, result{
					Code:    1001,
					Message: "删除失败",
				})
				return
			} else {
				c.JSON(http.StatusOK, result{
					Code:    1000,
					Message: "删除成功",
				})
			}
		})
	}

	//启动服务
	err = eng.Run()
	if err != nil {
		log.Fatalf("gin run failed,err:%v\n", err)
	}
}
