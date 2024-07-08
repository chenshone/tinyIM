package dao

import (
	"time"
	"tinyIM/data/db"
)

var dbIns = db.GetDb()

type User struct {
	Id         int `gorm:"primary_key"`
	Username   string
	Password   string
	CreateTime time.Time
}
