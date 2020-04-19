/**
 * @version: 1.0.0
 * @author: zhangguodong:general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: Gogland
 * @file: consulagent_test.go
 * @time: 2017/8/11 13:58
 */
package consul

import (
	"testing"

	`github.com/astaxie/beego/logs`
	"github.com/hashicorp/consul/api"
)

func TestConsulAgent_Register(t *testing.T) {
	// logs.Info(util.UUIDv4().ToFullString())

	// CreateCheckResponser_Http("http://127.0.0.1:12580/check")
	// CreateCheckResponser_Tcp("tcp://127.0.0.1:12580")

	uri := "http://127.0.0.1:8500"
	agentToken := "4c868079-3816-a681-b875-bbc7ae64be69"
	agent := GetConsulRemoteInst(uri, agentToken)
	if agent == nil {
		logs.Info("consul agent get fail.")
		return
	}

	service := &api.AgentServiceRegistration{
		ID:      "ActGuessSvr_Post_27217",
		Name:    "TestApp",
		Address: "127.0.0.1",
		Tags:    []string{},
		Port:    12580}

	// service.Check = &api.AgentServiceCheck{
	// 	HTTP: fmt.Sprintf("http://%s:%d%s", service.Address, service.Port, "/check"),
	// 	//TCP: fmt.Sprintf("%s:%d", service.Address, service.Port),
	// 	//TTL: "10s",
	// 	Timeout: "10s",
	// 	Interval: "15s",
	// 	DeregisterCriticalServiceAfter: "30s",
	// }

	// 定义健康监测，可以http,tcp,script,ttl的一种, 值为访问path, 能通过AgentService自动构建出完整的地址
	// attrs := make(map[string]string)
	// attrs["check_http"] = "/check"  通过http访问  是否返回200
	// attrs["check_tcp"] = "/check" //通过tcp访问
	// attrs["check_script"] = "cur http://127.0.0.1:8500/check"  //可以是Python， shell等脚本, 返回0正常
	// attrs["check_ttl"] = "ttl"
	// attrs["check_timeout"] = "10s"
	// attrs["check_interval"] = "15s"
	// attrs["check_deregister_after"] = "30s"

	if err := agent.Register(service, nil); err != nil {
		logs.Error("regist service fail: ", err)
	}

	if err := agent.Deregister(service.ID); err != nil {
		logs.Error("unregist service fail: ", err)
	}

	// quit := make(chan bool)
	// <-quit
}

func TestConsulAgent_PutKV(t *testing.T) {
	// uri := "http://127.0.0.1:8500"
	agent := GetConsulInst()
	if agent == nil {
		logs.Info("consul agent get fail.")
		return
	}

	if err := agent.PutKV("foo", "test value"); err != nil {
		logs.Info("put kv fail")
	} else {
		logs.Info("put kv success")
	}

	if val, err := agent.GetKV("foo"); err != nil {
		logs.Info("get kv fail")
	} else {
		logs.Info(val)
	}
}

func TestConsulAgent_FireEvent(t *testing.T) {
	// uri := "http://127.0.0.1:8500"
	agent := GetConsulInst()
	if agent == nil {
		logs.Info("consul agent get fail.")
		return
	}

	eName := "TestEvent"
	if err := agent.FireEvent(eName); err != nil {
		logs.Info("fire event fail")
	} else {
		logs.Info("fire event success:", eName)
	}

}

func TestConsulAgent_CheckService(t *testing.T) {
	// uri := "http://127.0.0.1:8500"
	agent := GetConsulInst()
	if agent == nil {
		logs.Info("consul agent get fail.")
		return
	}

	svrName := "TestApp1.0"

	res, err := agent.GetService(svrName, "")
	if err != nil {
		logs.Info("check service fail. ", svrName)
		return
	}

	// logs.Info(len(res))
	for i, v := range res {
		logs.Info(i, v.Service.ID, v.Service.Service, v.Service.Address, v.Service.Port)
	}

}
