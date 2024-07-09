package dao

import (
	"errors"
	"time"
	"tinyIM/data/db"

	"gorm.io/gorm"
)

var dbIns = db.GetDb()

type User struct {
	Id         int `gorm:"primary_key"`
	Username   string
	Password   string
	CreateTime time.Time
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) Add(dbInstance *gorm.DB) (userId int, err error) {
	if dbInstance == nil {
		dbInstance = dbIns
	}
	if u.Username == "" || u.Password == "" {
		return 0, errors.New("username or password empty")
	}
	user, ok := u.CheckByUsername(u.Username)
	if ok {
		return user.Id, nil
	}
	u.CreateTime = time.Now()
	if err = dbInstance.Table(u.TableName()).Create(&u).Error; err != nil {
		return 0, err
	}
	return u.Id, nil
}

func (u *User) CheckByUsername(username string) (data User, ok bool) {
	result := dbIns.Table(u.TableName()).Where("username = ?", username).Take(&data)
	return data, result.RowsAffected > 0
}

func (u *User) GetUserNameByUserId(userId int) (userName string) {
	var data User
	dbIns.Table(u.TableName()).Where("id=?", userId).Take(&data)
	return data.Username
}
