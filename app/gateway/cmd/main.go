package main

import (
	"context"
	"douyin/app/gateway/middleware"
	"douyin/app/gateway/router"
	"douyin/app/gateway/rpc"
	"douyin/config"
	"douyin/pkg/logger"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config.InitConfig()
	rpc.Init()
	middleware.EsInit()

	go startListen()
	{
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
		s := <-osSignals
		fmt.Println("exit! ", s)
	}
	fmt.Println("gateway listen on: http://192.168.1.7:8080")
}

func startListen() {
	r := router.NewRouter()
	server := &http.Server{
		Addr:           config.Conf.Server.Port,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("getway启动失败, err", err)
	}

	// 关闭
	go ShutDown(server)
}

func ShutDown(server *http.Server) {
	// 创建系统信号接收器接收关闭信号
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	logger.LogrusObj.Println("closing http server gracefully ...")

	if err := server.Shutdown(context.Background()); err != nil {
		logger.LogrusObj.Fatalln("closing http server gracefully failed: ", err)
	}
}
