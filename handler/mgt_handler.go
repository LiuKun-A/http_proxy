package handler

import (
	"encoding/json"
	"fmt"
	"http_proxy/config"
	"http_proxy/model"
	"log"
	"net/http"
)

type configHandlerManagement struct {
	mapping map[string]func(http.ResponseWriter, *http.Request)
}

func (chm *configHandlerManagement) initMapping() {
	chm.mapping["/mgt/service/update"] = mgtUpdate
	chm.mapping["/mgt/service/list"] = mgtList
	chm.mapping["/mgt/domain/update"] = domainUpdate
	chm.mapping["/mgt/domain/list"] = domainMappingList
}

func domainUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "<html><body>GET方法不支持， 只支持POST提交</body></html>\n")
	} else {
		param := model.DomainUpdateParams{}
		err := json.NewDecoder(r.Body).Decode(&param)
		if err != nil {
			log.Println("解析域名对应服务更新请求报错", err)
		} else {
			log.Println("接收到域名对应服务更新请求: ", param)
			if len(param.Data) > 0 {
				config.GetProxyConfig().UpdateMapping(param.Data)
			}
		}
	}
}

// 获取当前域名对应服务列表
func domainMappingList(w http.ResponseWriter, r *http.Request) {
	data := config.GetProxyConfig().ListMapping()
	byte, err := json.Marshal(data)
	if err != nil {
		log.Println("生成域名对应服务列表数据报错", err)
	} else {
		w.Header().Set("content-type", "text/json")
		fmt.Fprintf(w, string(byte))
	}
}

// 处理服务对应ip地址更新 只支持post
func mgtUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "<html><body>GET方法不支持， 只支持POST提交</body></html>\n")
	} else {
		param := model.MgtParam{}
		err := json.NewDecoder(r.Body).Decode(&param)
		if err != nil {
			log.Println("解析服务更新请求报错", err)
		} else {
			log.Println("接收到管理请求: ", param)
			if param.Action == 1 {
				config.GetServiceManagement().AddItem(param.ServiceName, param.Address)
			} else {
				config.GetServiceManagement().DelItem(param.ServiceName, param.Address)
			}
		}
	}
}

// 获取当前服务列表
func mgtList(w http.ResponseWriter, r *http.Request) {
	data := config.GetServiceManagement().ListItems()
	byte, err := json.Marshal(data)
	if err != nil {
		log.Println("生成服务列表数据报错", err)
	} else {
		w.Header().Set("content-type", "text/json")
		fmt.Fprintf(w, string(byte))
	}
}

func initConfigHandlerManagement() {
	chm := configHandlerManagement{
		mapping: make(map[string]func(http.ResponseWriter, *http.Request)),
	}
	chm.initMapping()
	if len(chm.mapping) > 0 {
		for url, handleFun := range chm.mapping {
			http.HandleFunc(url, handleFun)
		}
	}
}
