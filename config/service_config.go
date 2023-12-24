package config

import (
	"sync"
	"time"
)

type serviceItem struct {
	Name  string
	Items map[string]time.Time
	addr  []string
}

func (si *serviceItem) AddItems(addr ...string) {
	for _, addrName := range addr {
		si.Items[addrName] = time.Now()
	}
	si.initAddr()
}

func (si *serviceItem) initAddr() {
	tmpAddrSlice := make([]string, 10)
	for ip := range si.Items {
		tmpAddrSlice = append(si.addr, ip)
	}
	si.addr = tmpAddrSlice
}

func (si *serviceItem) GetAddrList() []string {
	return si.addr
}

func (si *serviceItem) DeleteItems(addr ...string) {
	for _, addrName := range addr {
		delete(si.Items, addrName)
	}
	si.initAddr()
}

var serviceManagementInstance *serviceManagement
var once sync.Once

func NewServiceItem(name string) serviceItem {
	return serviceItem{
		Name:  name,
		Items: make(map[string]time.Time),
		addr:  make([]string, 0),
	}
}

type serviceManagement struct {
	serviceMap map[string]serviceItem
}

func (sm *serviceManagement) GetItem(serviceName string) serviceItem {
	service, ok := sm.serviceMap[serviceName]
	if !ok {
		service = NewServiceItem(serviceName)
	}
	return service
}

func (sm *serviceManagement) AddItem(serviceName string, address []string) {
	service, ok := sm.serviceMap[serviceName]
	if !ok {
		service = NewServiceItem(serviceName)
	}
	service.AddItems(address...)
	sm.serviceMap[serviceName] = service
}

func (sm *serviceManagement) DelItem(serviceName string, address []string) {
	service, ok := sm.serviceMap[serviceName]
	if ok {
		service.DeleteItems(address...)
	}
}

func (sm *serviceManagement) ListItems() map[string]serviceItem {
	return sm.serviceMap
}

func GetServiceManagement() *serviceManagement {
	once.Do(func() {
		serviceManagementInstance = &serviceManagement{
			serviceMap: make(map[string]serviceItem),
		}
	})
	return serviceManagementInstance
}
