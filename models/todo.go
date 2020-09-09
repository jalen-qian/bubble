package models

import (
	"github.com/bubble/dao"
	"gorm.io/gorm"
)

/**
Todo 这个model的所有增删改查操作都放在这个文件
*/

type Todo struct {
	ID uint `json:"id" gorm:"primarykey"`
	gorm.Model
	Title  string `json:"title" gorm:"type:varchar(100);"`
	Status bool   `json:"status"`
}

func CreateATodo(todo *Todo) (err error) {
	err = dao.DB.Create(todo).Error
	return
}

func GetTodoList() (todoList *[]Todo, err error) {
	todoList = new([]Todo) //注意这里要分配内存
	err = dao.DB.Find(todoList).Error
	return
}

func GetFirst(todo *Todo) (err error) {
	err = dao.DB.First(todo).Error
	return
}

func UpdateATodo(todo *Todo) (err error) {
	err = dao.DB.Save(todo).Error
	return

}

func DeleteTodo(todo *Todo) (err error) {
	err = dao.DB.Delete(todo).Error
	return
}
