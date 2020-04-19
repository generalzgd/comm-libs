/**
 * @version: 1.0.0
 * @author: zhangguodong:general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: Gogland
 * @file: databasecfg.go
 * @time: 2017/6/29 9:45
 */
package svrcfg

import (
	"encoding/json"
	"fmt"
	"time"
)

// 数据库连接配置
type DatabaseCfg struct {
	Id           int    `json:"-" toml:"-" yaml:"-"`
	Name         string `json:"name" toml:"name" yaml:"name"`                         //
	Host         string `json:"host" toml:"host" yaml:"host"`                         //
	Port         uint   `json:"port" toml:"port" yaml:"port"`                         //
	Username     string `json:"user" toml:"username" yaml:"username"`                 //
	Password     string `json:"pwd" toml:"password" yaml:"password"`                  //
	Charset      string `json:"charset" toml:"charset" yaml:"charset"`                //
	Timeout      int    `json:"timeout" toml:"timeout" yaml:"timeout"`                //
	MaxIdleConns int    `json:"maxIdleConns" toml:"maxIdleConns" yaml:"maxIdleConns"` // 默认1
	MaxOpenConns int    `json:"maxOpenConns" toml:"maxOpenConns" yaml:"maxOpenConns"` // 默认1
}

// func NewDatabaseCfg() *DatabaseCfg {
// 	p := &DatabaseCfg{Charset: "utf8", Timeout: 2}
// 	return p
// }

func (p *DatabaseCfg) GetCharset() string {
	if len(p.Charset) == 0 {
		return "utf8"
	}
	return p.Charset
}

func (p *DatabaseCfg) GetTimeout() time.Duration {
	if p.Timeout <= 0 {
		return 2
	}
	return time.Duration(p.Timeout)
}

// 读取数据源
func (p *DatabaseCfg) GetDSN() string {
	// user,pwd, err := getDBAuthor()
	// if err != nil {
	//
	// }

	s := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&loc=Local",
		p.Username,
		p.Password,
		p.Host,
		p.Port,
		p.Name,
		p.GetCharset())
	return s
}

func (p *DatabaseCfg) FromJson(b []byte) error {
	return json.Unmarshal(b, &p)
}

func (p *DatabaseCfg) ToJson() ([]byte, error) {
	return json.Marshal(p)
}

//
// func getDBAuthor() (string, string, error) {
// 	switch libs.GetEnvName() {
// 	case "dev":
// 		return "live", "admin", nil
// 	case "beta":
// 		return "app_zhanqi", "xLgU^8hd", nil
// 	case "online":
// 		return "","",nil
// 	}
// 	return "","",errors.New("get database author fail")
// }
