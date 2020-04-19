/**
 * @version: 1.0.0
 * @author: zhangguodong:general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: Gogland
 * @file: rediscfg.go
 * @time: 2017/6/29 9:45
 */
package svrcfg

import (
	"encoding/json"
	"strconv"
)

// redis连接配置
type RedisCfg struct {
	Id           string `json:"-" toml:"-" yaml:"-"`
	Name         string `json:"name" toml:"name" yaml:"name"`                         //
	Address      string `json:"address" toml:"address" yaml:"address"`                // 可以是ip地址，也可以是unix domain socket地址(port必须为0)
	Port         int    `json:"port" toml:"port" yaml:"port"`                         //
	MaxIdleConns int    `json:"maxIdleConns" toml:"maxIdleConns" yaml:"maxIdleConns"` // 默认1
	MaxOpenConns int    `json:"maxOpenConns" toml:"maxOpenConns" yaml:"maxOpenConns"` // 默认10
	IdleTimeout  int    `json:"idleTimeout" toml:"idleTimeout" yaml:"idleTimeout"`    // 默认180‘
}

func (p *RedisCfg) GetDialType() string {
	// port 为0时，使用unix socket
	if p.Port <= 0 {
		return "unix"
	}
	// tcp连接
	return "tcp"
}

func (p *RedisCfg) GetAddress() string {
	// port 为0时，使用unix socket
	if p.Port <= 0 {
		return p.Address
	}
	// tcp连接
	return p.Address + ":" + strconv.Itoa(p.Port)
}

func (p *RedisCfg) FromJson(b []byte) error {
	return json.Unmarshal(b, &p)
}

func (p *RedisCfg) ToJson() ([]byte, error) {
	return json.Marshal(p)
}
