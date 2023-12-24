package loadbalance

import (
	"http_proxy/config"
	"log"
	"sync"
	"sync/atomic"
)

var rbMap sync.Map

type roundRobinCounter struct {
	count int32
}

func (r *roundRobinCounter) GetNextValue() int32 {
	nextVal := atomic.AddInt32(&r.count, 1)
	if nextVal > 100000000 {
		nextVal = 0
	}
	return nextVal
}

func GetNextAddr(host string) (string, bool) {
	serviceItem := config.GetProxyConfig().GetService(host)
	addr := serviceItem.Items
	if len(addr) == 0 {
		log.Println("domain: ", host, "对应服务: ", serviceItem.Name, " 地址为空")
		return "", false
	}
	counter, _ := rbMap.LoadOrStore(host, &roundRobinCounter{
		count: 0,
	})
	robinCounter := counter.(*roundRobinCounter)
	count := robinCounter.GetNextValue()
	var idx = count % int32(len(addr))

	return serviceItem.GetAddrList()[idx], true
}
