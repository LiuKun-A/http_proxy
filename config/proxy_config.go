package config

import (
	"http_proxy/model"
	"log"
	"sync"
)

var proxyOnce sync.Once
var proxyConfigInstance *proxyConfig

type proxyConfig struct {
	HealthCheckUrl string
	// 域名和后端服务名映射
	domainMapping map[string]string
}

func (p *proxyConfig) GetService(domain string) serviceItem {
	serviceName, ok := p.domainMapping[domain]
	if !ok {
		log.Println("未能找到", domain, "对应的servie配置， 请检查配置是否正确")
		return NewServiceItem(serviceName)
	}
	return GetServiceManagement().GetItem(serviceName)
}

func (p *proxyConfig) UpdateMapping(params []model.DomainUpdateParam) {
	for _, param := range params {
		p.domainMapping[param.Domain] = param.ServiceName
	}
}

func (p *proxyConfig) ListMapping() map[string]string {
	return p.domainMapping
}

func initProxyConfig(instance *proxyConfig) {
	instance.domainMapping["localhost:8080"] = "nginx-service"
}

func GetProxyConfig() *proxyConfig {
	proxyOnce.Do(func() {
		proxyConfigInstance = &proxyConfig{
			HealthCheckUrl: "/",
			domainMapping:  make(map[string]string),
		}
		initProxyConfig(proxyConfigInstance)
	})
	return proxyConfigInstance
}
