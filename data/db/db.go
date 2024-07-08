package db

import (
	"fmt"
	"tinyIM/config"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&&loc=Asia/Shanghai", config.Conf.Common.CommonMysql.UserName, config.Conf.Common.CommonMysql.Password, config.Conf.Common.CommonMysql.Host, config.Conf.Common.CommonMysql.Port, config.Conf.Common.CommonMysql.Database)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Errorf("connect db fail:%s\n", err.Error())
	}
}

func GetDb() *gorm.DB {
	return db
}