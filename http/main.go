package main

import (
	"GiftCode2/redis"
	"GiftCode2/router"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

// listenAddr 保存命令行选项 `--listen` 输入的监听地址
var listenAddr string

// redisAddr 保存命令行选项 `--redis` 输入的 redis 服务器地址
var redisAddr string

// rootCmd 表示一个命令行应用
var rootCmd = &cobra.Command{
	Use:   "GifCode",
	Short: "礼品码生成",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Infof("初始化 redis 客户端，服务器连接地址: %s", redisAddr)
		redis.InitRedis(redisAddr)

		logrus.Infof("启动 Web 服务器，监听地址: http://%s", listenAddr)
		err := router.InitRouter().Run(listenAddr)
		if err != nil {
			logrus.Fatalf("路由器运行失败: %v", err)
		}
	},
}

// main 是主执行方法
func main() {

	// 初始化日志格式
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 为命令行增加 flag
	// --listen 指定服务器的监听地址，默认 0.0.0.0:8080
	rootCmd.PersistentFlags().StringVarP(&listenAddr, "listen", "l", "0.0.0.0:8080", "服务器监听端口")
	// --redis 指定 redis 服务器地址，默认 127.0.0.1:6379
	rootCmd.PersistentFlags().StringVarP(&redisAddr, "redis", "r", "127.0.0.1:6379", "redis 服务器地址")

	// 执行主命令
	err := rootCmd.Execute()
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
