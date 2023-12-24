package handler

import (
	"http_proxy/loadbalance"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func initProxyHandlerManagement() {
	http.HandleFunc("/", ForwardHandler)
}

func ForwardHandler(writer http.ResponseWriter, request *http.Request) {
	serviceIp, ok := loadbalance.GetNextAddr(request.Host)
	if !ok {
		log.Println("获取", request.Host, "下的服务地址失败，请求被忽略")
		return
	}
	u := &url.URL{
		Scheme: "http",
		Host:   serviceIp,
	}

	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(writer, request)
}
