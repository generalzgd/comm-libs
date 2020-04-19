/**
 * @version: 1.0.0
 * @author: zhangguodong:general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: Gogland
 * @file: httpcfg.go
 * @time: 2017/6/29 9:44
 */
package svrcfg

import (
	"encoding/json"
	"strconv"
)

// http侦听配置
type HttpCfg struct {
	Host string `json:"host" toml:"host" yaml:"host"`
	Port int    `json:"port" toml:"port" yaml:"port"`
}

// 读取侦听地址
func (p *HttpCfg) GetListenAddr() string {
	return ":" + strconv.Itoa(p.Port)
}

// 读取连接地址
func (p *HttpCfg) GetLinkAddr() string {
	return p.Host + ":" + strconv.Itoa(p.Port)
}

func (p *HttpCfg) FromJson(b []byte) error {
	return json.Unmarshal(b, &p)
}

func (p *HttpCfg) ToJson() ([]byte, error) {
	return json.Marshal(p)
}
