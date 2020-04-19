/**
 * @version: 1.0.0
 * @author: zhangguodong:general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: Gogland
 * @file: consulwatcher.go
 * @time: 2017/8/11 10:29
 */
package consul

import (
	"errors"
	"net/url"
	"sync"
	"time"

	`github.com/astaxie/beego/logs`
	"github.com/hashicorp/consul/api/watch"

	"github.com/generalzgd/comm-libs/event"
)

const (
	WatchType_Key      = "key"
	WatchType_Services = "services"
	WatchType_Service  = "service"
	WatchType_Event    = "event"
)

type ConsulWatcher struct {
	sync.Mutex
	planMap map[string]*watch.Plan

	eventChan common.EventChannel
	address   string
}

// http://localhost:8500
func NewWatcher(address string) *ConsulWatcher {
	if len(address) == 0 {
		address = "http://127.0.0.1:8500"
	}
	uri, err := url.Parse(address)
	if err != nil {
		return nil
	}
	p := &ConsulWatcher{planMap: map[string]*watch.Plan{}}
	p.address = uri.Host
	return p
}

func NewDefaultWatcher() *ConsulWatcher {
	return NewWatcher("http://127.0.0.1:8500")
}

// add svr-frame event chan for watcher,
// if consul kv/service/services/consulevent change, user can get svr-frame event from chan
func (p *ConsulWatcher) AddEventChan(eventChan event.EventChannel) {
	p.eventChan = eventChan
}

func (p *ConsulWatcher) getPlan(watchType string, val string) (*watch.Plan, bool) {
	plan, ok := p.planMap[watchType+"_"+val]
	return plan, ok
}

func (p *ConsulWatcher) setPlan(watchType string, val string, plan *watch.Plan) {
	p.Lock()
	defer p.Unlock()

	mapKey := watchType + "_" + val
	p.planMap[mapKey] = plan
}

func (p *ConsulWatcher) stopPlan(watchType string, val string) {
	p.Lock()
	defer p.Unlock()

	mapKey := watchType + "_" + val
	if plan, ok := p.planMap[mapKey]; ok {
		if !plan.IsStopped() {
			plan.Stop()
		}
		delete(p.planMap, mapKey)
	}
}

func (p *ConsulWatcher) runPlan(plan *watch.Plan) error {
	errChan := make(chan error, 1)
	go func() {
		errChan <- plan.Run(p.address)
	}()

	select {
	case err := <-errChan:
		logs.Info("watch plan run fail: ", err)
		return err
	case <-time.After(1 * time.Second):
		logs.Info("watch plan is running")
	}
	return nil
}

func (p *ConsulWatcher) StopKeyWatch(key string) {
	p.stopPlan(WatchType_Key, key)
}

// listen kv change from consul
// one watch target one plan and one gorutine run
func (p *ConsulWatcher) AddKeyWatch(key string) error {
	watchType := WatchType_Key
	plan, ok := p.getPlan(watchType, key)
	if ok {
		if plan.IsStopped() {
			return errors.New("watch plan is stopped: " + watchType + "=" + key)
		} else {
			return errors.New("watch plan is running: " + watchType + "=" + key)
		}
	}

	params := make(map[string]interface{})
	params["type"] = watchType
	params["key"] = key

	plan, err := watch.Parse(params)
	if err != nil {
		return err
	}

	plan.Handler = func(idx uint64, val interface{}) {
		if val != nil {
			// logs.Info(fmt.Sprintf("%s", reflect.TypeOf(val).Elem().Name()))
			// switch d := val.(type) {
			// case *consulapi.KVPair:
			p.dispatchEvent(event.ConsulEvent, &ConsulKeyChange{Idx: idx, Value: val})
			// default:
			// }
		} else {
			// logs.Info(fmt.Sprintf("on %s plan handler: idx=%d, nil", watchType, idx))
			// p.dispatchEvent(event.ConsulEvent, &ConsulKeyChange{Idx: idx, Value: nil})
		}
	}

	if err := p.runPlan(plan); err != nil {
		return err
	}

	p.setPlan(watchType, key, plan)
	return nil
}

func (p *ConsulWatcher) StopServiceWatch(svrName string) {
	p.stopPlan(WatchType_Service, svrName)
}

// listen single service change from consul
func (p *ConsulWatcher) AddServiceWatch(svrName string, tag string, passingOnly bool) error {
	watchType := WatchType_Service
	plan, ok := p.getPlan(watchType, svrName)
	if ok {
		if plan.IsStopped() {
			return errors.New("watch plan is stopped: " + watchType + "=" + svrName)
		} else {
			return errors.New("watch plan is running: " + watchType + "=" + svrName)
		}
	}

	params := make(map[string]interface{})
	params["type"] = watchType
	params["service"] = svrName
	if len(tag) > 0 {
		params["tag"] = tag
	}
	params["passingonly"] = passingOnly

	plan, err := watch.Parse(params)
	if err != nil {
		return err
	}

	plan.Handler = func(idx uint64, val interface{}) {
		if val != nil {
			// logs.Info(fmt.Sprintf("%s", reflect.TypeOf(val).Elem().Name()))
			// switch d := val.(type) {
			// case []*consulapi.ServiceEntry:
			p.dispatchEvent(event.ConsulEvent, &ConsulServiceChange{Idx: idx, Value: val})
			// default:
			// }
		} else {
			// logs.Info("on "+watchType+" plan handler: ", idx, val)
			// p.dispatchEvent(event.ConsulEvent, &ConsulServiceChange{Idx: idx, Value: nil})
		}
	}

	if err := p.runPlan(plan); err != nil {
		return err
	}

	p.setPlan(watchType, svrName, plan)
	return nil
}

// cancel watch event from consul
func (p *ConsulWatcher) StopConsulEventWatch(eventName string) {
	p.stopPlan(WatchType_Event, eventName)
}

// listen event from consul by event fire
func (p *ConsulWatcher) AddConsulEventWatch(eventName string) error {
	watchType := WatchType_Event
	plan, ok := p.getPlan(watchType, eventName)
	if ok {
		if plan.IsStopped() {
			return errors.New("watch plan is stopped: " + watchType + "=" + eventName)
		} else {
			return errors.New("watch plan is running: " + watchType + "=" + eventName)
		}
	}

	params := make(map[string]interface{})
	params["type"] = watchType
	params["name"] = eventName

	plan, err := watch.Parse(params)
	if err != nil {
		return err
	}

	plan.Handler = func(idx uint64, val interface{}) {
		if val != nil {
			// logs.Info(fmt.Sprintf("%s", reflect.TypeOf(val).Elem().Name()))
			// switch d := val.(type) {
			// case []*consulapi.UserEvent:
			p.dispatchEvent(event.ConsulEvent, &ConsulEventChange{Idx: idx, Value: val})
			// default:
			// }
		} else {
			// logs.Info("on "+watchType+" plan handler: ", idx, val)
			// p.dispatchEvent(event.ConsulEvent, &ConsulEventChange{Idx: idx, Value: nil})
		}

	}

	if err := p.runPlan(plan); err != nil {
		return err
	}

	p.setPlan(watchType, eventName, plan)
	return nil
}

func (p *ConsulWatcher) StopServicesWatch() {
	p.stopPlan(WatchType_Services, "")
}

// listen all services from consul
func (p *ConsulWatcher) AddServicesWatch(stale bool) error {
	watchType := WatchType_Services
	plan, ok := p.getPlan(watchType, "")
	if ok {
		if plan.IsStopped() {
			return errors.New("watch plan is stopped: " + watchType)
		} else {
			return errors.New("watch plan is running: " + watchType)
		}
	}
	params := map[string]interface{}{}
	params["type"] = watchType
	params["stale"] = stale

	plan, err := watch.Parse(params)
	if err != nil {
		return err
	}

	plan.Handler = func(idx uint64, val interface{}) {
		if val != nil {
			p.dispatchEvent(common.ConsulEvent, &ConsulServicesChange{Idx: idx, Value: val})
		}
	}

	if err := p.runPlan(plan); err != nil {
		return err
	}

	p.setPlan(watchType, "", plan)
	return nil
}

func (p *ConsulWatcher) dispatchEvent(eType int, eData interface{}) {
	if p.eventChan != nil {
		p.eventChan <- event.NewEventObj(eType, eData)
	}

}
