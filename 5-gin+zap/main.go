package main

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"github.com/natefinch/lumberjack"
)

var logger *zap.Logger
var sugarLogger *zap.SugaredLogger

func InitLogger() {
	// 首先得到需要的三个配置
	writeSyncer := getLogWriter()
	encoder := getEncoder()

	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	// 加上 zap.AddCaller() 来添加上函数调用的信息，也就是会出来在哪一个包里面哪一行调用的这个
	logger = zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder{
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:       "msg",
		LevelKey:         "level",
		TimeKey:          "ts",
		NameKey:          "logger",
		CallerKey:        "caller",
		StacktraceKey:    "stacktrace",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      zapcore.LowercaseLevelEncoder,
		EncodeTime:       zapcore.ISO8601TimeEncoder,			// 这里改用人类可以识别的时间格式
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
	}
	// 像普通的console的形式来打印日志格式, 还是保存在文件中的
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename : "./test.log",
		MaxSize : 10,				// Mb
		MaxBackups : 5,
		MaxAge : 30,				// days
		Compress : false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func GinLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}


// GinRecovery recover掉项目可能出现的panic
func GinRecovery(logger *zap.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

func main() {
	/*
		func Default() *Engine {
			debugPrintWARNINGDefault()
			engine := New()
			engine.Use(Logger(), Recovery())
			return engine
		}

		Default()方法默认使用Logger(), Recovery()两个中间件；
		其中Logger()是把gin框架本身的日志输出到标准输出（我们本地开发调试时在
		终端输出的那些日志就是它的功劳），而Recovery()是在程序出现panic的时候
		恢复现场并写入500响应的。

		所以我们要模仿gin里面中间件的写法，自己写一个中间件进去，eg:

		func Logger() HandlerFunc {
			return LoggerWithConfig(LoggerConfig{})
		}
	*/
	//r := gin.Default()
	InitLogger()
	defer logger.Sync()
	r := gin.New()
	r.Use(GinLogger(logger), GinRecovery(logger,true))
	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK,"zcy")
	})
	r.Run()
}
