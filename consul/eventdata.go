/**
 * @version: 1.0.0
 * @author: zhangguodong:general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: Gogland
 * @file: eventdata.go
 * @time: 2017/8/11 12:45
 */
package consul

// import "github.com/hashicorp/consul/api"

type ConsulKeyChange struct {
	Idx   uint64
	Value interface{} // *api.KVPair
}

type ConsulServiceChange struct {
	Idx   uint64
	Value interface{} // []*api.ServiceEntry
}

type ConsulServicesChange struct {
	Idx   uint64
	Value interface{} // map[string][]string
}

type ConsulEventChange struct {
	Idx   uint64
	Value interface{} // []*api.UserEvent
}
