package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

var Cfg Config
var Db *gorm.DB
var RedisCli *redis.Client

func InitConfig(cfgFile string, debug bool) {
	v := viper.New()
	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		home, _ := os.Getwd()
		v.AddConfigPath(home)
		v.SetConfigName("")
	}
	v.SetConfigType("yaml")
	//dir := util.AbsPath()
	//configPath := filepath.Join(dir, "../config", fmt.Sprintf("%s.yaml", node))
	v.SetConfigFile(cfgFile)
	e := v.ReadInConfig()
	if e != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", e))
	}
	v.WatchConfig()
	if err := v.Unmarshal(&Cfg); err != nil {
		fmt.Println(err)
		panic("read config error:" + err.Error())
	}
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e)
		if err := v.Unmarshal(&Cfg); err != nil {
			fmt.Println(err)
		}
	})
	//debug mode
	if debug {
		Cfg.Zap.ConsoleOut = true
		Cfg.Zap.Level = "debug"
	}
	fmt.Printf("Using config file: %s, debug mode: %v \n", v.ConfigFileUsed(), debug)
	InitMysql()
	err := InitLogger()
	if err != nil {
		panic("InitLogger error:" + err.Error())
	}
	//InitRedis()
}

func ValidateEngineConfig(notify bool) error {
	if notify {
		if Cfg.DApp.GroupConversationId == "" {
			return fmt.Errorf("when notify is true, GroupConversationId must be not empty")
		}
	}
	return nil
}

func ValidateApiConfig() error {
	if Cfg.Deribit.ClientId == "" || Cfg.Deribit.Secret == "" {
		return fmt.Errorf("deribit clientID and Secret must be not empty")
	}
	return nil
}

func InitTestConfig(path string) {
	if path == "dev" {
		path = "../../config/od_dev_config/node1.yaml"
	}
	if path == "beta" {
		path = "../../config/od_beta_config/node1.yaml"
	}
	if path == "uat" {
		path = "../../config/od_uat_config/node1.yaml"
	}
	wd, _ := os.Getwd()
	join := filepath.Join(wd, path)
	InitConfig(join, true)
}

func InitRedis() {
	r := Cfg.Redis
	RedisCli = redis.NewClient(&redis.Options{
		Addr:     r.Address,
		Password: r.Password,
		DB:       r.Db,
	})
}
