package model

type MgtParam struct {
	// 服务名称
	ServiceName string `json:"serviceName"`
	// 要操作的动作，1：增加， 0：删除
	Action int `json:"action"`
	// 要操作的地址
	Address []string `json:"address"`
}

type DomainUpdateParam struct {
	// 域名
	Domain string `json:"domain"`
	// 服务名
	ServiceName string `json:"serviceName"`
}

type DomainUpdateParams struct {
	Data []DomainUpdateParam `json:"data"`
}
