//package main
//
//import (
//	"go.uber.org/zap"
//	"net/http"
//)
//
///*
//	Zap提供了两种类型的日志记录器—Sugared Logger和Logger。
//
//	在性能很好但不是很关键的上下文中，使用SugaredLogger。
//	它比其他结构化日志记录包快4-10倍，并且支持结构化和printf风格的日志记录。
//
//	在每一微秒和每一次内存分配都很重要的上下文中，使用Logger。
//	它甚至比SugaredLogger更快，内存分配次数也更少，但它只支持强类型的结构化日志记录。
//*/
//
///*
//	这个代码记录logger的使用
//*/
//
//var logger *zap.Logger
//
//func Initlogger() {
//	/*
//		通过调用zap.NewProduction()/zap.NewDevelopment()
//		或者zap.Example()创建一个Logger。
//
//		上面的每一个函数都将创建一个logger。唯一的区别在于它将记录的信息不同。
//		例如production logger默认记录调用函数信息、日期和时间等。
//	*/
//	// zap.NewProduction() 返回的是json格式的日志
//	logger, _ = zap.NewProduction()
//}
//
//func simpleHttpGet(url string) {
//	resp, err := http.Get(url)
//	if err != nil {
//		logger.Error(
//			"error fetching url...",
//			zap.String("url",url),
//			zap.Error(err))
//	} else {
//		logger.Info("success",
//			zap.String("stastuscode", resp.Status,),
//			zap.String("url",url))
//		resp.Body.Close()
//	}
//}
//
//func main() {
//	Initlogger()
//	// 在退出的时候把日志都刷新到磁盘上面
//	defer logger.Sync()
//
//	simpleHttpGet("www.google.com")
//	simpleHttpGet("http://www.google.cn")
//}
