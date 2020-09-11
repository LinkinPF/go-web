package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	//"os"

	"github.com/natefinch/lumberjack"
)

/*
	改造zap日志库的第一个需求就是把日志保存到文件之中，而不是打印到应用程序控制台

	使用zap.New(…)方法来手动传递所有配置，而不是使用像zap.NewProduction()这样的预置方法来创建logger。
	func New(core zapcore.Core, options ...Option) *Logger

	zapcore.Core需要三个配置——Encoder，WriteSyncer，LogLevel。

*/

var logger *zap.Logger
var sugarLogger *zap.SugaredLogger

func InitLogger() {
	// 首先得到需要的三个配置
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	// 这里的打印级别有：
	/*
		DebugLevel Level = iota - 1
		// InfoLevel is the default logging priority.
		InfoLevel
		// WarnLevel logs are more important than Info, but don't need individual
		// human review.
		WarnLevel
		// ErrorLevel logs are high-priority. If an application is running smoothly,
		// it shouldn't generate any error-level logs.
		ErrorLevel
		// DPanicLevel logs are particularly important errors. In development the
		// logger panics after writing the message.
		DPanicLevel
		// PanicLevel logs a message, then panics.
		PanicLevel
		// FatalLevel logs a message, then calls os.Exit(1).
		FatalLevel
	*/
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	// 加上 zap.AddCaller() 来添加上函数调用的信息，也就是会出来在哪一个包里面哪一行调用的这个
	logger = zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder{
	// 获取 JSON 格式的一个日志格式
	// 使用开箱即用的 NewJSONEncoder()，并使用预先设置的NewProductionEncoderConfig()
	//return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	// zap.NewProductionEncoderConfig()这个的默认配置也是定义了下面的结构体，所以我们自己也可以定义使用
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
	//file, _ := os.Create("./test.log")
	// 也可以使用追加的方式
	//file, _ := os.OpenFile("./test.log", os.O_CREATE | os.O_APPEND | os.O_RDWR, 0744)

	// 终极改造：使用第三方库，来实现日志切割
	lumberJackLogger := &lumberjack.Logger{
		Filename : "./test.log",
		MaxSize : 10,				// Mb
		MaxBackups : 5,
		MaxAge : 30,				// days
		Compress : false,
	}
	return zapcore.AddSync(lumberJackLogger)
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
	defer logger.Sync()
	simpleHttpGet("www.google.com")
	simpleHttpGet("http://www.google.com")
}
