package config

import (
	"github.com/spf13/viper"
	"os"
	"runtime"
	"strings"
	"sync"
)

var once sync.Once
var realPath string
var Conf *Config

type Config struct {
	Common Common
	Api    ApiConfig
}

func init() {
	Init()
}

func Init() {
	once.Do(func() {
		env := GetMode()
		realPath = getCurrentDir()
		cfgFilePath := realPath + "/" + env + "/"
		configNames := []string{"/common", "/api"}
		loadConfig(cfgFilePath, configNames)

		Conf = &Config{}
		if err := viper.Unmarshal(&Conf.Common); err != nil {
			panic(err)
		}
		if err := viper.Unmarshal(&Conf.Api); err != nil {
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

type Common struct {
	CommonEtcd CommonEtcd `mapstructure:"common-etcd"`
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

type ApiConfig struct {
	ApiBase ApiBase `mapstructure:"api-base"`
}

type ApiBase struct {
	ListenPort int `mapstructure:"listenPort"`
}
