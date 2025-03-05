package main

import (
	"github/lhh-gh/gift/database"
	"github/lhh-gh/gift/handler"
	"github/lhh-gh/gift/util"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

var (
	writeOrderFinish bool
)

func listenSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM) //注册信号2和15。Ctrl+C对应SIGINT信号

	for {
		sig := <-c //阻塞，直到信号的到来
		if writeOrderFinish {
			util.LogRus.Infof("receive signal %s, exit", sig.String())
			os.Exit(0) //进程退出
		} else {
			util.LogRus.Infof("receive signal %s, but not exit", sig.String())
		}
	}
}

func Init() {
	util.InitLog("log")
	database.InitGiftInventory()
	if err := database.ClearOrders(); err != nil {
		panic(err)
	} else {
		util.LogRus.Info("clear table orders")
	}

	handler.InitChannel()
	go func() {
		handler.TakeOrder() //把channel里的订单信息写入Mysql
		writeOrderFinish = true
	}()
	go listenSignal()

	// handler.InitMQ()
}

func main() {
	Init()

	//GIN自带logger和recover中间件
	//[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached

	gin.SetMode(gin.ReleaseMode)   //GIN线上发布模式
	gin.DefaultWriter = io.Discard //禁止GIN的输出
	// 修改静态资源不需要重启GIN，刷新页面即可
	router := gin.Default()

	router.Static("/js", "views/js") //在url是访问目录/js相当于访问文件系统中的views/js目录
	router.Static("/img", "views/img")
	router.StaticFile("/favicon.ico", "views/img/dqq.png") //在url中访问文件/favicon.ico，相当于访问文件系统中的views/img/dqq.png文件
	router.LoadHTMLFiles("views/lottery.html")             //使用这些.html文件时就不需要加路径了

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "lottery.html", nil)
	})
	router.GET("/gifts", handler.GetAllGifts) //获取所有奖品信息
	router.GET("/lucky", handler.Lottery)     //点击抽奖按钮

	router.Run("localhost:5678")
}

// go run ./main.go
