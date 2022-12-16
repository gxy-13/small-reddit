package main

import (
	"awesomeProject/controller"
	"awesomeProject/dao/mysql"
	"awesomeProject/dao/redis"
	"awesomeProject/logger"
	"awesomeProject/routers"
	"awesomeProject/settings"
	"awesomeProject/utils/snowflake"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	// 1. 加载配置文件
	if err := settings.Init(); err != nil {
		fmt.Printf("setting init failed, err:%v\n", err)
		return
	}
	// 2. 加载日志
	if err := logger.Init(settings.Conf.LogConf, settings.Conf.Mode); err != nil {
		fmt.Printf("logger init failed, err:%v\n", err)
		return
	}
	// 3. 链接mysql
	if err := mysql.Init(settings.Conf.MySQLConf); err != nil {
		fmt.Printf("mysql init failed, err:%v\n", err)
		return
	}
	defer mysql.Close()
	// 4. 链接redis
	if err := redis.Init(settings.Conf.RedisConf); err != nil {
		fmt.Printf("redis init failed, err:%v\n", err)
		return
	}
	defer redis.Close()
	// 初始化雪花算法,通过配置文件读取
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}
	// 初始化翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init translator failed, err:%v\n", err)
		return
	}
	// 5. 注册路由
	r := routers.Setup(settings.Conf.Mode)
	// 6. 项目启动，优雅关机
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("port")),
		Handler: r,
	}
	// 使用goroutine 保证代码可以继续往下执行
	go func() {
		// 开启一个goroutine 启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// 等待终端信号来优雅关闭服务器，为关闭服务器操作设置一个5s的超时
	quit := make(chan os.Signal, 1) //创建一个接收信号的通道
	// kill 默认发送，syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号， 我们常用的ctrl + c 就是输出系统的SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号， 但是不能被捕获，所以不添加他
	// signal.Notify 把收到的 syscall.SIGINT或 SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5s超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5s内优雅关闭服务，将来处理万的请求处理完在关闭服务，超过5s就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown:", zap.Error(err))
	}
	zap.L().Info("Server exiting")
}
