package main

import (
	handler2 "http_proxy/handler"
	"http_proxy/health"
	"log"
	"net/http"
)

func main() {
	var httpServer http.Server
	handler2.InitHandler()

	httpServer.Addr = ":8080"
	myTicker := health.StartHealthCheck()
	defer myTicker.Stop()
	log.Println("启动http代理服务，地址：", httpServer.Addr)
	go log.Println(httpServer.ListenAndServe())
}
