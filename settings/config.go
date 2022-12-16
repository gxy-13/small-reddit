package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 构建结构体用来存放配置信息

var Conf = new(AppConf)

type AppConf struct {
	Name       string `mapstructure:"name"`
	Version    string `mapstructure:"version"`
	Mode       string `mapstructure:"mode"`
	StartTime  string `mapstructure:"startTime"`
	MachineID  int64  `mapstructure:"machineID"`
	Port       int    `mapstructure:"port"`
	*RedisConf `mapstructure:"redis"`
	*MySQLConf `mapstructure:"mysql"`
	*LogConf   `mapstructure:"log"`
}

type RedisConf struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
}

type MySQLConf struct {
	Host     string `mapstructure:"host"'`
	Password string `mapstructure:"password"`
	User     string `mapstructure:"user"`
	DB       string `mapstructure:"dbname"`
	Port     int    `mapstructure:"port"`
}

type LogConf struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackUps int    `mapstructure:"max_backups"`
}

func Init() (err error) {
	// 通过viper读取配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadInConfig failed,err:%v\n", err)
		return
	}
	if viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Ummarshal failed, err:%v\n", err)
		return
	}
	// 监视配置文件，实现热加载
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了")
		if viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper Unmarshal failed,%v\n", err)
			return
		}
	})
	return
}
