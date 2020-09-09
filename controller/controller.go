package controller

import (
	"errors"
	"github.com/bubble/models"
	"github.com/bubble/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func GetTodoList(c *gin.Context) {
	if todoList, err := models.GetTodoList(); err != nil {
		//有错误，返回错误信息
		c.JSON(http.StatusOK, service.Result{
			Code:    1000,
			Message: "获取待办列表失败，err:" + err.Error(),
			Data:    todoList,
		})
	} else {
		//无错误，返回正常数据
		c.JSON(http.StatusOK, service.Result{
			Code:    1000,
			Message: "获取待办列表成功",
			Data:    todoList,
		})
	}

}

func CreateATodo(c *gin.Context) {
	var todo models.Todo
	err := c.BindJSON(&todo)
	if err != nil {
		c.JSON(http.StatusBadRequest, service.Result{Code: 1001, Message: "数据不合法，请传JSON格式数据"})
		return
	}
	//将数据保存到数据库
	if err := models.CreateATodo(&todo); err != nil {
		c.JSON(http.StatusOK, service.Result{Code: 1001, Message: err.Error()})
	} else {
		c.JSON(http.StatusOK, service.Result{Code: 1000, Message: "添加待办事项成功"})
	}
}

/**
更新一个Todo的状态
*/
func UpdateATodo(c *gin.Context) {
	//获取路由参数中的id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, service.Result{
			Code:    1001,
			Message: "参数不正确",
		})
		return
	}
	//参数正确，判断参数是否大于0，如果参数为0，也不合法
	if id <= 0 {
		c.JSON(http.StatusOK, service.Result{
			Code:    1001,
			Message: "参数不能小于0",
		})
		return
	}
	var todo = models.Todo{ID: uint(id)}
	//找到这条记录
	if err := models.GetFirst(&todo); err != nil {
		//找不到，分为1-数据库没有这条记录
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, service.Result{
				Code:    1001,
				Message: "未找到这条记录",
			})
		} else { //或者数据库连接不上了
			c.JSON(http.StatusOK, service.Result{
				Code:    1001,
				Message: "数据访问异常：" + err.Error(),
			})
		}
		return
	} else {
		//能找到这条记录，那么更新状态
		log.Println(todo)
		//将状态设置为与原来相反
		todo.Status = !todo.Status
		//更新
		if err := models.UpdateATodo(&todo); err != nil {
			c.JSON(http.StatusOK, service.Result{
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
			c.JSON(http.StatusOK, service.Result{
				Code:    1000,
				Message: "事项【" + todo.Title + "】设置为" + statusName,
			})
		}
	}
}

func DeleteATodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, service.Result{
			Code:    1001,
			Message: "参数不正确",
		})
		return
	}
	//参数正确，判断参数是否大于0，如果参数为0，也不合法
	if id <= 0 {
		c.JSON(http.StatusOK, service.Result{
			Code:    1001,
			Message: "参数不能小于0",
		})
		return
	}
	if err := models.DeleteTodo(&models.Todo{ID: uint(id)}); err != nil {
		c.JSON(http.StatusOK, service.Result{
			Code:    1001,
			Message: "删除失败",
		})
		return
	} else {
		c.JSON(http.StatusOK, service.Result{
			Code:    1000,
			Message: "删除成功",
		})
	}
}
