package main

import (
	"go.uber.org/zap"
	"net/http"
)

/*
	Zap提供了两种类型的日志记录器—Sugared Logger和Logger。

	在性能很好但不是很关键的上下文中，使用SugaredLogger。
	它比其他结构化日志记录包快4-10倍，并且支持结构化和printf风格的日志记录。

	在每一微秒和每一次内存分配都很重要的上下文中，使用Logger。
	它甚至比SugaredLogger更快，内存分配次数也更少，但它只支持强类型的结构化日志记录。
*/

/*
	这个代码记录Sugared Logger的使用，实现的功能和上一个logger实现的功能一样

	惟一的区别是，我们通过调用主logger的. Sugar()方法来获取一个SugaredLogger。
	然后使用SugaredLogger以printf格式记录语句

	其实可以发现 logger 和 SugaredLogger 是可以替换的
*/

var sugarLogger *zap.SugaredLogger

func InitLogger() {
	// 把 zap.NewProduction() 换成 zap.Development() 输出的就是不是json格式了，而是终端上显示的格式
	logger, _ := zap.NewProduction()
	sugarLogger = logger.Sugar()
}

func simpleHttpGet(url string) {
	sugarLogger.Debugf("Trying to hit GET request for %s",url)
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Errorf("Error fetching URL %s : Error = %s",
			url, err)
	} else {
		sugarLogger.Infof("Success! statusCode = %s for URL %s",
			resp.Status,
			url)
		resp.Body.Close()
	}
}

func main() {
	InitLogger()
	defer sugarLogger.Sync()
	simpleHttpGet("www.google.com")
	simpleHttpGet("http://www.google.com")
}

