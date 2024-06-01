package main

import (
	"FeishuCozeRobot/config"
	"FeishuCozeRobot/handlers"
	"FeishuCozeRobot/routers"
	"FeishuCozeRobot/utils"
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
	"time"
)

//go:embed config/app.json
var configjson string

func main() {

	conf, err := config.ChangeConfig(configjson)
	if err != nil {
		log.Printf("读取配置文件失败: %v", err.Error())
	}
	handlers.InitRedisUtil(conf.RedisConfig)

	handlers.InitHandlers(conf)

	logger := enableLog()
	defer utils.CloseLogger(logger)
	// 获取当前运行目录
	currentDir, err := os.Getwd()
	if err != nil {

		panic("获取当前运行目录失败，" + err.Error())
	}
	fmt.Println("当前运行目录：", currentDir)
	// 判断文件夹是否存在，不存在则创建
	folderPath := currentDir + "\\tempfile\\"
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.Mkdir(folderPath, 0755)
		if err != nil {
			panic("创建文件夹失败，" + err.Error())
		}
		fmt.Println("文件夹已创建：", folderPath)
	} else {
		fmt.Println("文件夹已存在：", folderPath)
	}

	// 注册处理器 默认开启日志打印
	//g := gin.Default()
	// 注册处理器 默认关闭日志打印
	g := gin.New()
	//设置日志级别为 gin.DebugLevelNone，不打印请求路径日志
	g.Use(utils.CustomMiddleware())

	g.GET("/ping", func(c *gin.Context) {

		c.Header("Server", "Go-Gin-Server")
		c.JSON(200, gin.H{
			"message":   "pong",
			"code":      200,
			"success":   true,
			"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
		})
	})

	g.POST("/ping", func(c *gin.Context) {
		plainEventJsonStr, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(400, "Error reading request body")
			return
		}
		c.Header("Server", "Go-Gin-Server")
		c.JSON(200, gin.H{
			"message":   "pong",
			"code":      200,
			"result":    string(plainEventJsonStr),
			"success":   true,
			"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
		})
	})
	//添加路由
	routers.RegisterRouter(g)
	// 启动WEB服务
	err = g.Run(":" + conf.AppPort)
	if err != nil {
		log.Printf("failed to start server: %v", err)
	}

}

func enableLog() *lumberjack.Logger {
	// Set up the logger
	var logger *lumberjack.Logger
	logger = &lumberjack.Logger{
		Filename: "logs/coze_robot.log",
		MaxSize:  100,      // megabytes
		MaxAge:   365 * 10, // days
	}

	fmt.Printf("logger %T\n", logger)
	// Set up the logger to write to both file and console
	log.SetOutput(io.MultiWriter(logger, os.Stdout))
	log.SetFlags(log.Ldate | log.Ltime)
	// Write some log messages
	log.Println("Starting application...")

	return logger
}
