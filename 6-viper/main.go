package main

import (
	"bytes"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	/*
		Viper会按照下面的优先级。每个项目的优先级都高于它下面的项目:

		显示调用Set设置值
		命令行参数（flag）
		环境变量
		配置文件
		key/value存储
		默认值

		目前Viper配置的键（Key）是大小写不敏感的。目前正在讨论是否将这一选项设为可选。
	*/

	// 设置默认值
	viper.SetDefault("filepath", ".")
	viper.SetDefault("filename", "haha")

	// 读取配置文件，首先需要知道在哪里查找配置文件，viper默认不配置搜索路径
	viper.SetConfigFile("./config.yaml")
	// 也可以分开配置：
	viper.SetConfigName("config")		// 配置文件名称，没有扩展名
	viper.SetConfigType("yaml")			// 如果配置文件名称中没有扩展名，就需要配置这个

	// 设置查找配置文件所在的路径
	viper.AddConfigPath("/etc/appname")
	viper.AddConfigPath("￥HOME/.appname")
	viper.AddConfigPath(".")

	// 查找并读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _,ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到错误
		} else {
			// 找到了配置文件，但是还有另外的错误
		}
	}

	// 代码执行到这里，就正确读取了配置了

	/*
		接下来就是要写入配置文件了

		从配置文件中读取配置文件是有用的，但是有时你想要存储在运行时所做的所有修改。
		为此，可以使用下面一组命令

		WriteConfig - 将当前的viper配置写入预定义的路径并覆盖（如果存在的话）。如果没有预定义的路径，则报错。
		SafeWriteConfig - 将当前的viper配置写入预定义的路径。如果没有预定义的路径，则报错。如果存在，将不会覆盖当前的配置文件。
		WriteConfigAs - 将当前的viper配置写入给定的文件路径。将覆盖给定的文件(如果它存在的话)。
		SafeWriteConfigAs - 将当前的viper配置写入给定的文件路径。不会覆盖给定的文件(如果它存在的话)。

		根据经验，标记为safe的所有方法都不会覆盖任何文件，而是直接创建（如果不存在），而默认行为是创建或截断。
	*/
	viper.WriteConfig()			// 将当前配置写入 viper.AddConfigPath() 和 viper.SetConfigName() 里面去
	viper.SafeWriteConfig()
	viper.WriteConfigAs("/path/to/my/.config")
	viper.SafeWriteConfigAs("/path/to/my/.other_config")

	/*
		下一个功能是实时监控配置文件的变化

		需要重新启动服务器以使配置生效的日子已经一去不复返了，viper驱动的应用程序
		可以在运行时读取配置文件的更新，而不会错过任何消息。

		只需告诉viper实例watchConfig。可选地，你可以为Viper提供一个回调函数，
		以便在每次发生更改时运行。
	*/
	// 确保在调用WatchConfig()之前添加了所有的配置路径!!!
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变化后会调用的回调函数
		fmt.Println("Config file changed:", e.Name)
	})

	/*
		在gin里面配置使用viper

		这里的例子是访问/version，会返回给前端版本号的信息，然后动态修改一下配置文件，
		再重新访问/version，就会发现返回了一个修改过后的值。这就是热加载
	*/
	r := gin.Default()
	r.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, viper.GetString("version"))
	})
	// r.Run()

	/*
		viper 也可以从io.Reader里面读取配置
		Viper预先定义了许多配置源，如文件、环境变量、标志和远程K/V存储，但你不受其约束。
		你还可以实现自己所需的配置源并将其提供给viper。
	*/
	// 任何需要将此配置添加到程序中的方法。
	var yamlExample = []byte(`
		Hacker: true
		name: steve
		hobbies:
		- skateboarding
		- snowboarding
		- go
		clothing:
		  jacket: leather
		  trousers: denim
		age: 35
		eyes : brown
		beard: true
		`)

	viper.ReadConfig(bytes.NewBuffer(yamlExample))

	viper.Get("name") // 这里会得到 "steve"

	/*
		覆盖设置
		这些可能来自命令行标志，也可能来自于自己的应用程序逻辑
	*/
	viper.Set("filename", "hehe")
	viper.Set("logfile", "logfile")

	/*
		注册和使用别名
		别名允许单个建引用单个值
	*/
	viper.RegisterAlias("loud", "verbose")		// 这里就是 loud 和 verbose 建立了别名
	viper.Set("verbose", true)		// 结果和下一行相同
	viper.Set("loud", true)			// 结果和上一行相同

	viper.GetBool("loud")
	viper.GetBool("verbose")

	/*
		从环境变量中去读取配置，这个先放下，暂时用不到
	*/

	/*
		从命令行参数中去读取配置，这个暂时也用不到，先放下
	*/

	/*
		从远程 key/value 中去读取配置，这个暂时也用不到，先放下
	*/


	/*
		开始从配置文件中读取配置的值，viper 提供了如下的方法：

		Get(key string) : interface{}
		GetBool(key string) : bool
		GetFloat64(key string) : float64
		GetInt(key string) : int
		GetIntSlice(key string) : []int
		GetString(key string) : string
		GetStringMap(key string) : map[string]interface{}
		GetStringMapString(key string) : map[string]string
		GetStringSlice(key string) : []string
		GetTime(key string) : time.Time
		GetDuration(key string) : time.Duration
		IsSet(key string) : bool
		AllSettings() : map[string]interface{}

		需要认识到的一件重要事情是，每一个Get方法在找不到值的时候都会返回零值。
		为了检查给定的键是否存在，提供了IsSet()方法。
	*/
	viper.GetString("port")
	if viper.GetBool("verbose") {
		fmt.Println("verbose enabled")
	}

	// 也可以访问嵌套的键
	/*
		对于下面的JSON文件：
		{
			"host": {
				"address": "localhost",
				"port": 5799
			},
			"datastore": {
				"metric": {
					"host": "127.0.0.1",
						"port": 3099
				},
				"warehouse": {
					"host": "198.0.0.1",
						"port": 2112
				}
			}
		}
	*/
	// viper 可以通过传入 . 分隔的路径来访问嵌套的配置文件
	viper.GetString("datastore.metric.host")		// 返回 "127.0.0.1"

}







