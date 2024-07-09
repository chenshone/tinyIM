package config

import (
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var once sync.Once
var realPath string
var Conf *Config

const (
	SuccessReplyCode      = 0
	FailReplyCode         = 1
	SuccessReplyMsg       = "success"
	QueueName             = "tinyim_queue"
	RedisBaseValidTime    = 86400
	RedisPrefix           = "tinyim_"
	RedisRoomPrefix       = "tinyim_room_"
	RedisRoomOnlinePrefix = "tinyim_room_online_count_"
	MsgVersion            = 1
	OpSingleSend          = 2 // single user
	OpRoomSend            = 3 // send to room
	OpRoomCountSend       = 4 // get online user count
	OpRoomInfoSend        = 5 // send info to room
	OpBuildTcpConn        = 6 // build tcp conn
)

type Config struct {
	Common Common
	Api    ApiConfig
	Logic  LogicConfig
}

func init() {
	Init()
}

func Init() {
	once.Do(func() {
		env := GetMode()
		realPath = getCurrentDir()
		cfgFilePath := realPath + "/" + env + "/"
		configNames := []string{"/common", "/api", "/logic"}
		loadConfig(cfgFilePath, configNames)

		Conf = &Config{}
		if err := viper.Unmarshal(&Conf.Common); err != nil {
			panic(err)
		}
		if err := viper.Unmarshal(&Conf.Api); err != nil {
			panic(err)
		}
		if err := viper.Unmarshal(&Conf.Logic); err != nil {
			panic(err)
		}
	})
}

func loadConfig(cfgFilePath string, configNames []string) {
	viper.SetConfigType("toml")
	viper.AddConfigPath(cfgFilePath)

	for _, name := range configNames {
		viper.SetConfigName(name)
		if err := viper.MergeInConfig(); err != nil {
			panic(err)
		}
	}
}

func getCurrentDir() string {
	_, file, _, _ := runtime.Caller(1)
	path := strings.Split(file, "/")
	dir := strings.Join(path[:len(path)-1], "/")
	return dir
}

func GetMode() string {
	env := os.Getenv("RUN_MODE")
	if env == "" {
		env = "dev"
	}
	return env
}

func GetGinRunMode() string {
	env := GetMode()
	//gin have debug,test,release mode
	if env == "dev" {
		return "debug"
	}
	if env == "test" {
		return "debug"
	}
	if env == "prod" {
		return "release"
	}
	return "release"
}

type LogicConfig struct {
	LogicBase LogicBase `mapstructure:"logic-base"`
}

type LogicBase struct {
	ServerId   string `mapstructure:"serverId"`
	CpuNum     int    `mapstructure:"cpuNum"`
	RpcAddress string `mapstructure:"rpcAddress"`
	CertPath   string `mapstructure:"certPath"`
	KeyPath    string `mapstructure:"keyPath"`
}

type Common struct {
	CommonEtcd  CommonEtcd  `mapstructure:"common-etcd"`
	CommonMysql CommonMysql `mapstructure:"common-mysql"`
	CommonRedis CommonRedis `mapstructure:"common-redis"`
}

type CommonEtcd struct {
	Host              string `mapstructure:"host"`
	BasePath          string `mapstructure:"basePath"`
	ServerPathLogic   string `mapstructure:"serverPathLogic"`
	ServerPathConnect string `mapstructure:"serverPathConnect"`
	UserName          string `mapstructure:"userName"`
	Password          string `mapstructure:"password"`
	ConnectionTimeout int    `mapstructure:"connectionTimeout"`
}

type CommonMysql struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	UserName string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type CommonRedis struct {
	RedisAddress  string `mapstructure:"redisAddress"`
	RedisPassword string `mapstructure:"redisPassword"`
	Db            int    `mapstructure:"db"`
}

type ApiConfig struct {
	ApiBase ApiBase `mapstructure:"api-base"`
}

type ApiBase struct {
	ListenPort int `mapstructure:"listenPort"`
}
