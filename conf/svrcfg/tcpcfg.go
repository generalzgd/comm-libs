/**
 * @version: 1.0.0
 * @author: zhangguodong:general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: Gogland
 * @file: tcpcfg.go
 * @time: 2017/6/29 9:43
 */
package svrcfg

import (
	"encoding/json"
	"strconv"
)

// tcp侦听配置，包括post，内部服务之间的tcp
type TcpCfg struct {
	Host             string `json:"host" toml:"host" yaml:"host"`                                     //
	Port             int    `json:"port" toml:"port" yaml:"port"`                                     //
	ConnMax          int    `json:"connMax" toml:"connMax" yaml:"connMax"`                            //
	ReadBufferSize   int    `json:"readBufferSize" toml:"readBufferSize" yaml:"readBufferSize"`       //
	WriteBufferSize  int    `json:"writeBufferSize" toml:"writeBufferSize" yaml:"writeBufferSize"`    //
	SendChanLimit    int    `json:"sendChanLimit" toml:"sendChanLimit" yaml:"sendChanLimit"`          //
	ReceiveChanLimit int    `json:"receiveChanLimit" toml:"receiveChanLimit" yaml:"receiveChanLimit"` //
}

// 读取buffer字节大小，默认1024
func (p *TcpCfg) GetReadBufferSize() uint32 {
	if p.ReadBufferSize <= 0 {
		return 1024
	}
	return uint32(p.ReadBufferSize)
}

// 读取写入buffer字节大小，默认1024
func (p *TcpCfg) GetWriteBufferSize() uint32 {
	if p.WriteBufferSize <= 0 {
		return 1024
	}
	return uint32(p.WriteBufferSize)
}

// 读取发送队列长度，默认1024个包
func (p *TcpCfg) GetSendChanLimit() uint32 {
	if p.SendChanLimit <= 0 {
		return 1024
	}
	return uint32(p.SendChanLimit)
}

// 读物接收队列长度，默认1024个包
func (p *TcpCfg) GetReceiveChanLimit() uint32 {
	if p.ReceiveChanLimit <= 0 {
		return 1024
	}
	return uint32(p.ReceiveChanLimit)
}

// 读取侦听的地址
func (p *TcpCfg) GetListenAddr() string {
	return ":" + strconv.Itoa(p.Port)
}

// 读取连接的地址
func (p *TcpCfg) GetLinkAddr() string {
	return p.Host + ":" + strconv.Itoa(p.Port)
}

// func NewTcpCfg() *TcpCfg {
// 	p := &TcpCfg{readBufferSize: 1024, writeBufferSize: 1024, sendChanLimit: 1024, receiveChanLimit: 1024}
// 	return p
// }

func (p *TcpCfg) FromJson(b []byte) error {
	return json.Unmarshal(b, &p)
}

func (p *TcpCfg) ToJson() ([]byte, error) {
	return json.Marshal(p)
}
