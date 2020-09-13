package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/*

	我们编写的Web项目部署之后，经常会因为需要进行配置变更或功能迭代而重启服务，
	单纯的kill -9 pid的方式会强制关闭进程，这样就会导致服务端当前正
	在处理的请求失败

	什么是优雅关机：
	优雅关机就是服务端关机命令发出后不是立即关机，而是等待当前还在处理的请求
	全部处理完毕后再退出程序，是一种对客户端友好的关机方式。而执行Ctrl+C关闭服务端时，
	会强制结束进程导致正在访问的请求出现问题。
*/

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           router,
	}

	// 这里为什么要单独开启一个goroutine？因为如果不开启的话，在ListenAndServe()之后，
	// 就会阻塞在这里了，一直等待监听，代码就不会往下执行了，所以包装成一个单独的goroutine
	go func() {
		// 开启一个 goroutine 启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed{
			fmt.Println("listen faied :", err)
		}
	}()

	// 等待中断信号来优雅的关闭本地服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1)		// 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)		// 这里不会阻塞
	<- quit		// 在这里阻塞，当接收到上面两种信号的时候才会往下进行
	fmt.Println("shutdown server ...")
	// 创建一个5s超时的context，就是等待5s之后，无论如何，都会把服务器关掉 shutdown 掉
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel( )

	// net 包默认支持优雅关机了
	// 5s 内优雅关闭服务，也就是把未处理完的请求处理完以后再关闭服务，超过5s就超时直接退出
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("Server shutdown", err)
	}

	fmt.Println("Server exiting")
}
