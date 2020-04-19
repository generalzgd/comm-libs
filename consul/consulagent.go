/**
 * @version: 1.0.0
 * @author: zgd: general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: Gogland
 * @file: consulagent.go
 * @time: 2017/8/10 下午7:58
 */

package consul

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	`github.com/astaxie/beego/logs`
	consulapi "github.com/hashicorp/consul/api"

	libs "github.com/generalzgd/comm-libs"
)

type ConsulAgent struct {
	sync.Mutex
	client          *consulapi.Client
	curRegistration *consulapi.AgentServiceRegistration
	exit            bool
	onceClose       sync.Once
	nodeIp          string
}

var (
	once  sync.Once
	agent *ConsulAgent
)

const (
	defaultToken = "anonymous"
)

func makeInst(uriStr string, token string) (*ConsulAgent, error) {
	if len(uriStr) == 0 {
		uriStr = "http://127.0.0.1:8500"
	}

	uri, err := url.Parse(uriStr)
	if err != nil {
		logs.Error("url parse error: ", err)
		return nil, err
	}

	config := consulapi.DefaultConfig()
	config.Address = uriStr

	if len(token) > 0 {
		config.Token = token
	} else {
		config.Token = defaultToken
	}

	client, err := consulapi.NewClient(config)
	if err != nil {
		logs.Error("consul: ", uri.Scheme)
		return nil, err
	}

	agent := &ConsulAgent{client: client}
	return agent, nil
}

// sample: http://localhost:8500
// 获取默认（本地）consul连接, 其他连接使用GetConsulRemoteInst方法
func GetConsulInst() *ConsulAgent {
	if agent == nil {
		once.Do(func() {
			uriStr := "http://127.0.0.1:8500"
			p, err := makeInst(uriStr, "")
			if err != nil {
				logs.Error("make consul instance fail: ", err)
			}
			agent = p
		})
	}

	return agent
}

// 获取非本地consul连接
func GetConsulRemoteInst(uriStr string, token string) *ConsulAgent {
	if agent == nil {
		once.Do(func() {
			p, err := makeInst(uriStr, token)
			if err != nil {
				logs.Error("make consul instance fail: ", err)
			}
			agent = p
		})
	}

	return agent
}

func (p *ConsulAgent) Destroy() {
	p.onceClose.Do(func() {
		p.exit = true
	})
}

func (p *ConsulAgent) Run() {

}

// localIp本地调试时可以用，正式环境后请传空字符串
func (p *ConsulAgent) Init(svrName string, useType string, svrPort int, svrType int, healthPort int, healthType string, localIp string) *consulapi.AgentServiceRegistration {
	ip := localIp // p.GetNodeIp()
	if len(ip) == 0 {
		ip = libs.GetInnerIp()
	}
	reg := &consulapi.AgentServiceRegistration{
		ID:      strings.ToLower(fmt.Sprintf("%s_%s_%d", svrName, useType, libs.Ip2Long(ip))),
		Name:    strings.ToLower(fmt.Sprintf("%s", svrName)),
		Tags:    []string{"type:" + strconv.Itoa(svrType), strings.ToLower(useType)},
		Port:    svrPort,
		Address: ip,
	}
	reg.Check = &consulapi.AgentServiceCheck{
		TCP:                            fmt.Sprintf("%s:%d", ip, svrPort),
		Timeout:                        "1s",
		Interval:                       "15s",
		DeregisterCriticalServiceAfter: "30s",
		Status:                         "passing",
	}

	if healthType == "http" {
		reg.Check.HTTP = fmt.Sprintf("http://%s:%d%s", reg.Address, healthPort, "/health")
		reg.Check.TCP = ""
	}

	if len(reg.Check.HTTP) > 0 {
		p.RunHealthCheck(reg.Check.HTTP)
	}
	p.curRegistration = reg
	return reg
}

func (p *ConsulAgent) GetCurrentResgitration() *consulapi.AgentServiceRegistration {
	return p.curRegistration
}

func (p *ConsulAgent) RunHealthCheck(addr string) error {
	// addr := "http://127.0.0.1:"+strconv.Itoa(port)+"/health"
	uri, err := url.Parse(addr)
	if err != nil {
		return err
	}

	http.HandleFunc(uri.Path, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("success"))
	})
	go http.ListenAndServe(uri.Host, nil)
	return nil
}

// Ping will try to connect to consul by attempting to retrieve the current leader.
func (p *ConsulAgent) Ping() error {
	p.Lock()
	defer p.Unlock()

	status := p.client.Status()
	leader, err := status.Leader()
	if err != nil {
		return err
	}
	log.Println("consul: current leader ", leader)

	return nil
}

// regist new service
func (p *ConsulAgent) Register(registration *consulapi.AgentServiceRegistration, svrAttrs map[string]string) error {
	p.Lock()
	defer p.Unlock()

	if registration == nil {
		registration = p.curRegistration
	}

	if svrAttrs != nil {
		registration.Check = p.buildCheck(registration.Address, registration.Port, svrAttrs)
	}

	err := p.client.Agent().ServiceRegister(registration)
	if err != nil {
		return err
	}

	return nil
}

func (p *ConsulAgent) buildCheck(svrIp string, svrPort int, svrAttrs map[string]string) *consulapi.AgentServiceCheck {
	check := &consulapi.AgentServiceCheck{}

	if svrAttrs != nil {
		if status := svrAttrs["check_initial_status"]; status != "" {
			check.Status = status
		}
		if path := svrAttrs["check_http"]; path != "" {
			check.HTTP = fmt.Sprintf("http://%s:%d%s", svrIp, svrPort, path)
			if timeout := svrAttrs["check_timeout"]; timeout != "" {
				check.Timeout = timeout
			}
		} else if path := svrAttrs["check_https"]; path != "" {
			check.HTTP = fmt.Sprintf("https://%s:%d%s", svrIp, svrPort, path)
			if timeout := svrAttrs["check_timeout"]; timeout != "" {
				check.Timeout = timeout
			}
			// } else if cmd := svrAttrs["check_cmd"]; cmd != "" {
			// 	check.Script = fmt.Sprintf("check-cmd %s %s %s", service.Origin.ContainerID[:12], service.Origin.ExposedPort, cmd)
			// } else if script := svrAttrs["check_script"]; script != "" {
			// 	check.Script = p.interpolateService(script, svrIp, svrPort)

		} else if ttl := svrAttrs["check_ttl"]; ttl != "" {
			check.TTL = ttl
		} else if tcp := svrAttrs["check_tcp"]; tcp != "" {
			check.TCP = fmt.Sprintf("%s:%d", svrIp, svrPort)
			if timeout := svrAttrs["check_timeout"]; timeout != "" {
				check.Timeout = timeout
			}
		} else {
			return nil
		}
		/*if check.Script != "" || check.HTTP != "" || check.TCP != "" {
			if interval := svrAttrs["check_interval"]; interval != "" {
				check.Interval = interval
			} else {
				check.Interval = "10s"
			}
		}*/
		if deregisterAfter := svrAttrs["check_deregister_after"]; deregisterAfter != "" {
			check.DeregisterCriticalServiceAfter = deregisterAfter
		}
	}

	return check
}

func (p *ConsulAgent) interpolateService(script string, svrIp string, svrPort int) string {
	withIp := strings.Replace(script, "$SERVICE_IP", svrIp, -1)
	withPort := strings.Replace(withIp, "$SERVICE_PORT", strconv.Itoa(svrPort), -1)
	return withPort
}

// unregist service
func (p *ConsulAgent) Deregister(svrId string) error {
	p.Lock()
	defer p.Unlock()

	if len(svrId) == 0 {
		svrId = p.curRegistration.ID
	}

	return p.client.Agent().ServiceDeregister(svrId)
}

func (p *ConsulAgent) Refresh() error {
	return nil
}

// get services in map
func (p *ConsulAgent) Services() (map[string]*consulapi.AgentService, error) {
	p.Lock()
	defer p.Unlock()

	return p.client.Agent().Services()
}

func (p *ConsulAgent) GetService(svrName string, tag string) ([]*consulapi.ServiceEntry, error) {
	p.Lock()
	defer p.Unlock()

	if len(svrName) == 0 && p.curRegistration != nil {
		svrName = p.curRegistration.Name
	}

	hel := p.client.Health()
	res, _, err := hel.Service(svrName, tag, false, nil)
	if err != nil {
		return nil, err
	}

	// list := make([]*consulapi.AgentService, len(res))
	// for i, s := range res {
	// 	list[i] = s.Service
	// }
	return res, nil
}

func (p *ConsulAgent) PutKV(key string, value string) error {
	p.Lock()
	defer p.Unlock()

	kv := p.client.KV()

	meta, err := kv.Put(&consulapi.KVPair{Key: key, Value: []byte(value)}, nil)
	if err != nil {
		logs.Error("consulkv: failed to register service: ", err)
		return err
	}
	logs.Info(meta)
	return err
}

func (p *ConsulAgent) GetKV(key string) (string, error) {
	p.Lock()
	defer p.Unlock()

	v := ""
	kv := p.client.KV()

	pair, meta, err := kv.Get(key, nil)
	if err != nil {
		return v, err
	}
	if pair == nil {
		return v, errors.New("consul kv: unexpected nil value")
	}
	// if pair.Flags != 42 {
	// 	return v, errors.New("consul kv: unexpected value flags!=42")
	// }
	if meta.LastIndex == 0 {
		return v, errors.New("consul kv: unexpected value meta.LastIndex==0")
	}

	v = string(pair.Value)
	return v, nil
}

func (p *ConsulAgent) GetKVList(prefix string) (consulapi.KVPairs, error) {
	p.Lock()
	defer p.Unlock()

	kv := p.client.KV()

	kvpairs, meta, err := kv.List(prefix, nil)
	if err != nil {
		return nil, err
	}
	if meta.LastIndex == 0 {
		return nil, errors.New("consul kv: unexpected value meta.LastIndex==0")
	}
	return kvpairs, nil
}

func (p *ConsulAgent) DelKV(key string) error {
	p.Lock()
	defer p.Unlock()

	kv := p.client.KV()
	if _, err := kv.Delete(key, nil); err != nil {
		return err
	}
	return nil
}

//
func (p *ConsulAgent) FireEvent(eventType string) error {
	p.Lock()
	defer p.Unlock()

	event := p.client.Event()

	userEvent := &consulapi.UserEvent{Name: eventType}
	id, meta, err := event.Fire(userEvent, nil)
	if err != nil {
		return err
	}

	if meta.RequestTime == 0 {
		return errors.New("consul event: fire event RequestTime==0")
	}

	if id == "" {
		return errors.New("consul event: fire event id =='' ")
	}
	return nil
}

// func (p *ConsulAgent) ListEvents() ([]*consulapi.UserEvent,error) {
// 	p.Lock()
// 	defer p.Unlock()
//
// 	event := p.client.Event()
//
// 	events, _, err := event.List("", nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return events, nil
// }

// 获取当前节点ip
func (p *ConsulAgent) GetNodeIp() string {
	if len(p.nodeIp) > 0 {
		return p.nodeIp
	}
	m, err := p.client.Agent().Self()
	if err != nil {
		return ""
	}
	member, ok := m["Member"]
	if !ok {
		return ""
	}
	addr, ok := member["Addr"]
	if ok {
		p.nodeIp = addr.(string)
	}
	return p.nodeIp
}
