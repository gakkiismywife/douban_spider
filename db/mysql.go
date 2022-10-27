package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var Db *gorm.DB

type BaseModel struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func init() {
	Db = SetupDb()
}

func SetupDb() *gorm.DB {
	dsn := "root:qweqwe@tcp(120.78.67.238:3306)/spider?charset=utf8mb4&parseTime=True&loc=Local"
	Db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	db, err := Db.DB()
	if err != nil {
		log.Println("mysql connect error", err)
		return nil
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(2)
	db.SetConnMaxIdleTime(time.Minute)

	return Db
}
