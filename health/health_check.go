package health

import (
	"http_proxy/config"
	"http_proxy/model"
	"log"
	"net/http"
	"sync"
	"time"
)

type healthCheckJob struct {
	CheckChannel chan model.MgtParam
}

func (h *healthCheckJob) check() {
	log.Println("正在进行health check....")
	serviceMap := config.GetServiceManagement().ListItems()
	for _, serviceItem := range serviceMap {
		for address := range serviceItem.Items {
			go func(serviceName string, address string) {
				healthUrl := "http://" + address + config.GetProxyConfig().HealthCheckUrl
				resp, err := http.Get(healthUrl)
				isSuccess := true
				if err != nil || resp.StatusCode != 200 {
					isSuccess = false
				}
				if !isSuccess {
					items := make([]string, 1)
					items = append(items, address)
					data := model.MgtParam{
						ServiceName: serviceName,
						Action:      0,
						Address:     items,
					}
					h.CheckChannel <- data
				}

			}(serviceItem.Name, address)
		}
	}
}

var healthCheckJobInstance *healthCheckJob
var once sync.Once

func GetHealthCheckJob() *healthCheckJob {
	once.Do(func() {
		healthCheckJobInstance = &healthCheckJob{
			CheckChannel: make(chan model.MgtParam, 100),
		}
	})
	return healthCheckJobInstance
}

func StartHealthCheck() *time.Ticker {
	myticker := time.NewTicker(5 * time.Second)
	go func(t *time.Ticker) {
		for {
			<-t.C
			GetHealthCheckJob().check()
		}
	}(myticker)
	go func() {
		for {
			data := <-GetHealthCheckJob().CheckChannel
			log.Println("health check failure, del: ", data)
			config.GetServiceManagement().DelItem(data.ServiceName, data.Address)
		}
	}()
	return myticker
}
