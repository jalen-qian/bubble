package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//返回所有的清单列表
func GetAllToDo(c *gin.Context) {
	c.JSON(http.StatusOK,todoList)
}
