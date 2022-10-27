package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"spider_douban/config"
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
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%s&loc=Local",
		config.Database.Username,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.DbName,
		config.Database.Charset,
		config.Database.ParseTime,
	)
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
