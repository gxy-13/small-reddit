package logger

import (
	"awesomeProject/settings"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Init 初始化zap，通过配置文件中的mode，控制日志的输出地点
func Init(cfg *settings.LogConf, mode string) (err error) {
	writeSyncer := getLogWriter(cfg.FileName, cfg.MaxAge, cfg.MaxSize, cfg.MaxBackUps)
	encoder := getEncoder()
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return
	}
	var core zapcore.Core
	// 当 "dev" 模式时
	if mode == "dev" {
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		// 用zap.NewTee() 配置多个日志输出的地方, zapcore.Lock(os.Stdout) 将标准输出转换为WriteSyncer 类型
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writeSyncer, l),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
	} else {
		core = zapcore.NewCore(encoder, writeSyncer, l)
	}
	// zap.New() 手动配置zap，而不是直接通过zap.NewProduction() 这样的预制logger
	// New() 有两个参数 zapCore.Core Option
	// zapCore.Core 有三个参数，encoder, WriteSyncer, LevelEnabler
	// encoder 编码器 使用开箱即用的NewJsonEncoder() 并且可以使用预先设置的newProductionEncoderConfig(), writeSyncer 将配置文件写到哪里去，使用zapcore.AddSync()设置地址，
	// levelEnabler 表示哪种级别的日志将被写入
	logger := zap.New(core, zap.AddCaller())
	// 替换zap库中全局变量，可以直接通过zap.L()访问
	zap.ReplaceGlobals(logger)
	return
}

func getLogWriter(filename string, maxAge, maxSize, maxBackups int) zapcore.WriteSyncer {
	// os.Create() 创建日志存放文件，
	//file, _ := os.Create("./bluebell.log")
	//return zapcore.AddSync(file)

	// 需要对日志进行归档和切割，使用lumberjack
	lumberjack := &lumberjack.Logger{
		Filename:   filename,   // 日志文件名
		MaxAge:     maxAge,     // 保留旧文件的最大天数
		MaxSize:    maxSize,    // 切割之前，日志文件的最大大小 MB
		MaxBackups: maxBackups, // 保留旧文件的最大个数
	}
	return zapcore.AddSync(lumberjack)
}

func getEncoder() zapcore.Encoder {
	// 预先设置好的NewProductionEncoderConfig()的日志时间为时间戳，我们修改时间字段并且将日志中的隔离级别用大写字母显示
	encoderConfig := zap.NewProductionEncoderConfig()
	// 标准时间
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 大写字母表示
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		// 先去执行其余的中间件
		c.Next()

		cost := time.Since(start)
		zap.L().Info(path,
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
func GinRecovery(stack bool) gin.HandlerFunc {
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
					zap.L().Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					zap.L().Error("[Recovery from panic]",
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
