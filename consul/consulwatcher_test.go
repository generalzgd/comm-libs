/**
 * @version: 1.0.0
 * @author: zhangguodong:general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: Gogland
 * @file: consulwatcher_test.go
 * @time: 2017/8/11 10:58
 */
package consul

import (
	"testing"

	`github.com/astaxie/beego/logs`
	"github.com/hashicorp/consul/api"

	"github.com/generalzgd/comm-libs/event"
)

func readEventChan(t *testing.T, eventChan event.EventChannel) {
	for {
		select {
		case eObj := <-eventChan:
			switch d := eObj.Data.(type) {
			case *ConsulKeyChange:
				if d.Value != nil {
					kv := d.Value.(*api.KVPair)
					logs.Info("ConsulKeyChange: ", d.Idx, kv.Key, string(kv.Value))
				} else {
					logs.Info("ConsulKeyChange: ", d.Idx, nil)
				}
			case *ConsulServiceChange:
				if d.Value != nil {
					sli := d.Value.([]*api.ServiceEntry)
					for _, v := range sli {
						logs.Info("ConsulServiceChange: ", v.Service.ID, v.Service.Service, v.Service.Address, v.Service.Port)
					}
				} else {
					logs.Info("ConsulServiceChange: ", d.Idx, nil)
				}

			case *ConsulServicesChange:
				if d.Value != nil {
					svrs := d.Value.(map[string][]string)
					for n := range svrs {
						logs.Info("ConsulServicesChange: ", n)
						TestConsulAgent_CheckService(t)
					}
				} else {
					logs.Info("ConsulServicesChange: ", d.Idx, nil)
				}

			case *ConsulEventChange:
				if d.Value != nil {
					sli := d.Value.([]*api.UserEvent)
					for _, v := range sli {
						logs.Info("ConsulEventChange: ", v.ID, v.Name)
					}
				} else {
					logs.Info("ConsulEventChange: ", d.Idx, nil)
				}
			}
		}
	}
}

func TestConsulWatcher_AddKeyWatch(t *testing.T) {
	uri := "http://192.168.163.184:8500"
	w := NewWatcher(uri)
	eventChan := make(event.EventChannel, 100)
	w.AddEventChan(eventChan)

	go readEventChan(t, eventChan)

	key := "urlAddr/userInfo"
	err := w.AddKeyWatch(key)
	if err != nil {
		logs.Error("key watch fail: ", err)
	}

	//quit := make(chan bool)
	//<-quit

	// w.StopKeyWatch(key)
}

func TestConsulWatcher_AddServiceWatch(t *testing.T) {
	uri := "http://127.0.0.1:8500"
	w := NewWatcher(uri)
	eventChan := make(event.EventChannel, 100)
	w.AddEventChan(eventChan)

	go readEventChan(t, eventChan)

	svrName := "TestApp1.0"

	err := w.AddServiceWatch(svrName, "", false)
	if err != nil {
		logs.Error("service watch fail: ", err)
	}

	//quit := make(chan bool)
	//<-quit

	w.StopServiceWatch(svrName)
}

func TestConsulWatcher_AddServicesWatch(t *testing.T) {
	uri := "http://127.0.0.1:8500"
	w := NewWatcher(uri)
	eventChan := make(event.EventChannel, 100)
	w.AddEventChan(eventChan)

	go readEventChan(t, eventChan)

	err := w.AddServicesWatch(false)
	if err != nil {
		logs.Error("service watch fail: ", err)
	}

	//quit := make(chan bool)
	//<-quit

	w.StopServicesWatch()
}

func TestConsulWatcher_AddEventWatch(t *testing.T) {
	uri := "http://127.0.0.1:8500"
	w := NewWatcher(uri)
	eventChan := make(event.EventChannel, 100)
	w.AddEventChan(eventChan)

	go readEventChan(t, eventChan)

	eventName := "TestEvent"
	err := w.AddConsulEventWatch(eventName)
	if err != nil {
		logs.Error("service watch fail: ", err)
	}

	//quit := make(chan bool)
	//<-quit

	w.StopConsulEventWatch(eventName)
}
