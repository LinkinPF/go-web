package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type config struct {
	Port int		`mapstructure:"port"`
	Version string  `mapstructure:"version"`
	Mysql			`mapstructure:"mysql"`
}

type Mysql struct {
	Dbname string	`mapstructure:"dbname"`
	Host string		`mapstructure:"host"`
	Port int		`mapstructure:"port"`
}

var C config

func main() {
	viper.SetConfigFile("./config.yaml")
	viper.AddConfigPath(".")

	// 查找并读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _,ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到错误
		} else {
			// 找到了配置文件，但是还有另外的错误
		}
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变化后会调用的回调函数
		fmt.Println("Config file changed:", e.Name)
	})

	/*
		viper 可以把选定的或者所有的值解析到结构体、map中去

		有两个方法可以做：
		Unmarshal(rawVal interface{}) : error
		UnmarshalKey(key string, rawVal interface{}) : error
	*/
	// 反序列化，开始读取配置信息
	err := viper.Unmarshal(&C)
	if err != nil {
		fmt.Println("viper.Unmarshal failed : ",err)
	}
	fmt.Printf("c:%#v\n", C)

	// 序列化成字符串
	c := viper.AllSettings()
	bs, err := yaml.Marshal(c)
	if err != nil {
		fmt.Printf("unable to marshal config to YAML: %v", err)
	}
	fmt.Println(string(bs))
}









