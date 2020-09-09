package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDb() (err error) {
	dsn := "root:123456@tcp(127.0.0.1)/bubble?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db.Debug()
	sqlDb, _ := DB.DB()
	return sqlDb.Ping()
}
