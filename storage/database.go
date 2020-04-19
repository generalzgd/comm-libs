/**
 * @version: 1.0.0
 * @author: zgd: general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: Gogland
 * @file: database.go
 * @time: 2017/9/15 上午9:48
 */

package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
	"xorm.io/core"
	`github.com/astaxie/beego/logs`
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/generalzgd/comm-libs/conf/svrcfg"
)

// 数据库名
const (
	ActiveDBName    = "actives"
	ContentDbName   = "contents"
	VideoDbName     = "new_videos"
	LogDbName       = "log"
	ConfigDbName    = "config"
	UserDbname      = "user"
	UsersDataDbName = "users_data"
	AuthUserDbName  = "oauth_user"
	RichDbName      = "rich"
	ConfDbName      = "conf"
	ImDbName        = "im"
	HfImDBName      = "hf_im"
	HfEsports		= "hf_esports"
)

var (
	NoDbErr      = errors.New("no db setted")
	TimeoutDbErr = errors.New("db connect timeout")
	// LinkFailDbErr    = errors.New("db link fail")
	LogFileFailDbErr = errors.New("db create log file fail")
	DsnErr           = errors.New("mysql dsn empty")
)

// 数据库链接对象
type DbLinkHolder struct {
	sync.Mutex

	Name         string
	DSN          string
	Conn         *sql.DB
	lastUseTime  time.Time
	MaxIdleConns int
	MaxOpenConns int
	//
	engin   *xorm.Engine
	dbDebug int
}


func (p *DbLinkHolder) GetEngin() (*xorm.Engine, error) {

	if len(p.DSN) == 0 {
		logs.Debug("GetEngin 3.", DsnErr)
		return nil, DsnErr
	}
	if p.engin == nil && len(p.DSN) > 0 {
		engine, err := p.newEngine()
		if err != nil {
			logs.Error("get xorm.Engine err: %v , create: one more time.", err)
			engine, err = p.newEngine()
		}
		p.engin = engine
		return engine, err
	}
	if time.Since(p.lastUseTime) > time.Minute {
		err := p.engin.Ping()
		if err != nil {
			engine, err := p.newEngine()
			if err != nil {
				logs.Error("get xorm.Engine err: %v ,timeout: one more time.", err)
				return engine, err
			}
			p.engin = engine
			// 需要外部重新获取
			return engine, err
		}
	}
	return p.engin, nil
}

func (p *DbLinkHolder) newEngine() (*xorm.Engine, error) {
	e, err := xorm.NewEngine("mysql", p.DSN)
	if err != nil {
		logs.Debug("GetEngin 1.", err)
		return nil, err
	} else {
		logs.Debug("NewEngin mysql.", p.DSN)
	}

	timeout := make(chan bool)
	go func() {
		err = e.Ping()
		timeout <- true
	}()
	// 数据库连接超时控制
	select {
	case <-time.After(time.Second * 2):
		return nil, TimeoutDbErr
	case <-timeout:
		break
	}
	if err != nil {
		logs.Debug("GetEngin 2.", err)
		return nil, err
	}

	if p.dbDebug == 1 {
		// 打印执行的每句SQL
		e.ShowSQL(true)
		// 设置日志级别
		e.Logger().SetLevel(core.LOG_INFO)
		// 设置日志写入文件
		logFile, err := os.Create("./logs/sql.log")
		if err != nil {
			return nil, LogFileFailDbErr
		}
		e.SetLogger(xorm.NewSimpleLogger(logFile))
	}

	// 设置连接池最大最小连接数
	if p.MaxIdleConns > 0 {
		e.SetMaxIdleConns(p.MaxIdleConns)
	} else {
		e.SetMaxIdleConns(1)
	}
	if p.MaxOpenConns > 0 {
		e.SetMaxOpenConns(p.MaxOpenConns)
	} else {
		e.SetMaxOpenConns(10)
	}
	return e, nil
}

func (p *DbLinkHolder) GetLink() error {
	if p.Conn == nil && len(p.DSN) > 0 {
		c, err := sql.Open("mysql", p.DSN)
		if err != nil {
			logs.Info("init mysql link failure. Error:", err.Error())
			return err
		}
		if p.MaxIdleConns > 0 {
			c.SetMaxIdleConns(p.MaxIdleConns)
		} else {
			c.SetMaxIdleConns(1)
		}

		if p.MaxOpenConns > 0 {
			c.SetMaxOpenConns(p.MaxOpenConns)
		} else {
			c.SetMaxOpenConns(1)
		}

		err = c.Ping()
		if err != nil {
			logs.Info("init mysql ping failure. Error:", err.Error())
			return err
		}
		p.Conn = c
		logs.Info("init mysql db success. dsn:", p.DSN)
	}

	if len(p.DSN) == 0 {
		return errors.New("mysql dsn empty")
	}

	if time.Since(p.lastUseTime) > time.Minute {
		err := p.Conn.Ping()
		return err
	}

	return nil
}

func (p *DbLinkHolder) CloseLink() {
	if p.Conn != nil {
		p.Conn.Close()
		p.Conn = nil
		logs.Info("mysql db closed. ", p.Name)
	}
}

func (p *DbLinkHolder) RecordQueryTime() {
	p.lastUseTime = time.Now()
}

// //////////////////////////////////////////////////////////
type IDbCtrl interface {
	AddCfg(configs map[string]*svrcfg.DatabaseCfg)
	GetDbLink(dbName string) (*DbLinkHolder, error)
	CloseDbLink(dbName string)
	CloseAllDbLink()
}

// //////////////////////////////////////////////////////////
// 数据库控制器
type DbCtrl struct {
	// sync.Mutex

	DbMap  map[string]*DbLinkHolder
	ticker *time.Ticker

	bClosed  bool
	exitChan chan bool
	onClose  sync.Once
}

func NewDbCtrl() *DbCtrl {
	t := &DbCtrl{
		DbMap:    make(map[string]*DbLinkHolder, 10),
		ticker:   time.NewTicker(time.Minute * 5),
		exitChan: make(chan bool, 1),
	}
	go t.wait()
	return t
}

func (p *DbCtrl) destroy() {
	p.onClose.Do(func() {
		if p.ticker != nil {
			p.ticker.Stop()
		}
		if p.exitChan != nil {
			close(p.exitChan)
		}

		p.CloseAllDbLink()
		p.bClosed = true
	})
}

func (p *DbCtrl) wait() {
	defer p.destroy()

	for {
		select {
		case <-p.ticker.C:
			p.onTicker()
		case <-p.exitChan:
			break
		}
	}
}

func (p *DbCtrl) onTicker() {
	// p.Lock()
	// defer p.Unlock()

	for _, holder := range p.DbMap {
		span := time.Since(holder.lastUseTime)
		if span > time.Minute*5 {
			holder.CloseLink()
		}
	}
}

func (p *DbCtrl) AddCfg(configs map[string]*svrcfg.DatabaseCfg) {
	for _, c := range configs {
		p.DbMap[c.Name] = &DbLinkHolder{
			DSN: c.GetDSN(),
		}
		logs.Debug("init db AddCfg.", c.Name, c.GetDSN())
	}
}

func (p *DbCtrl) GetDbLink(dbName string) (*DbLinkHolder, error) {
	if holder, ok := p.DbMap[dbName]; ok {
		err := holder.GetLink()
		if err != nil {
			return nil, err
		}

		return holder, nil
	}
	return nil, NoDbErr
}

func (p *DbCtrl) CloseDbLink(dbName string) {
	if holder, ok := p.DbMap[dbName]; ok {
		holder.CloseLink()
	}
}

func (p *DbCtrl) CloseAllDbLink() {
	for _, h := range p.DbMap {
		h.CloseLink()
	}
	logs.Info("all mysql db closed.")
}

//
func (p *DbCtrl) GetDbEngin(dbname string) (*xorm.Engine, error) {
	if h, ok := p.DbMap[dbname]; ok {
		e, err := h.GetEngin()
		if err != nil {
			logs.Error("GetDbEngin err and try ty.", err, dbname)
			// h.engin = nil
			e, err = h.GetEngin() // 再获取一次
			if err != nil {
				return nil, err
			}
		}
		h.RecordQueryTime()
		return e, nil
	}
	return nil, NoDbErr
}

// ****************************************************************
// func (p *DbCtrl) GetUserInfo(uid uint) (*db_ds.UserInfo, error) {
// 	h, err := p.GetDbLink(UserDbname)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	h.Lock()
// 	defer h.Unlock()
//
// 	sqlstr := fmt.Sprintf("SELECT nickname,gender,avatar FROM users WHERE uid=%d", uid)
// 	rows, err := h.Conn.Query(sqlstr)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	user := &db_ds.UserInfo{}
// 	if rows.Next() {
// 		user.Uid = uid
// 		if err := rows.Scan(&user.Nickname, &user.Gender, &user.Avatar); err != nil {
// 			return nil, err
// 		}
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	if err := rows.Close(); err != nil {
// 		return nil, err
// 	}
//
// 	return user, nil
// }
