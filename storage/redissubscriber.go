/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: subscriber.go
 * @time: 2020/2/17 3:00 下午
 * @project: duoduoteamsvr
 */

package storage

import (
	`sync`

	`github.com/astaxie/beego/logs`
	`github.com/garyburd/redigo/redis`

	`github.com/generalzgd/comm-libs/conf/svrcfg`
)

type SubscribeCallback func([]byte)

type RedisSubscriber struct {
	closeOnce sync.Once
	//
	cfg  *svrcfg.RedisCfg
	conn redis.PubSubConn
	//
	handlers map[string]SubscribeCallback
}

func NewSubscriber(cfg *svrcfg.RedisCfg) *RedisSubscriber {
	return &RedisSubscriber{
		cfg:      cfg,
		handlers: map[string]SubscribeCallback{},
	}
}

func (p *RedisSubscriber) AddListener(key string, callback SubscribeCallback) {
	p.handlers[key] = callback
}

func (p *RedisSubscriber) Init() error {
	conn, err := redis.Dial("tcp", p.cfg.GetAddress())
	if err != nil {
		return err
	}
	//
	p.conn = redis.PubSubConn{Conn: conn} // conn
	//
	ch := make([]interface{}, 0, len(p.handlers))
	for k := range p.handlers {
		ch = append(ch, k)
	}
	p.conn.Subscribe(ch...)

	//p.conn.Unsubscribe()
	//
	go p.run()
	return nil
}

func (p *RedisSubscriber) Close() {
	p.closeOnce.Do(func() {
		p.conn.Unsubscribe()
		p.conn.Close()
	})
}

func (p *RedisSubscriber) run() {
	defer func() {
		p.Close()
		if r := recover(); r != nil {

		}
	}()

	for {
		switch res := p.conn.Receive().(type) {
		case redis.Message:
			channel := res.Channel // (*string)(unsafe.Pointer(&res.Channel))
			message := res.Data    //  (*string)(unsafe.Pointer(&res.Data))
			if handler, ok := p.handlers[channel]; ok {
				handler(message)
			}
		case redis.Subscription:
			logs.Info("redis.Subscription got=[%s: %s %d]", res.Channel, res.Kind, res.Count)
		case redis.Error:
			logs.Info("redis.Subscription got err=[%v]", res)
		}
	}

}
