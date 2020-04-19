/**
 * @version: 1.0.0
 * @author: zgd: general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: Gogland
 * @file: redis.go
 * @time: 2017/9/15 下午1:56
 */

package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	`github.com/astaxie/beego/logs`
	"github.com/bitly/go-simplejson"
	"github.com/garyburd/redigo/redis"

	libs "github.com/generalzgd/comm-libs"
	"github.com/generalzgd/comm-libs/conf/svrcfg"
)

// redis 库名定义
const (
	DefaultRdsName = "default" // 默认为common实例
	ActiveRdsName  = "active"
	UserRdsName    = "user"
	SessionRdsName = "session"
	CacheRdsName   = "cache"
	CommonRdsName  = "common"

	StorageRdsName = "storage"
	QueueRdsName   = "queue" // storage
	LevelRdsName   = "level"
	SocialRdsName  = "social"
	CredisRdsName  = "credis"
	DataRdsName    = "data"   // storage
	PubSubRdsName  = "pubsub" // 发布订阅
)

type RedisLinkPool struct {
	sync.Mutex

	Cfg *svrcfg.RedisCfg
	// Conn        redis.Conn
	lastUseTime time.Time
	Pool        *redis.Pool // 连接池
}

func (p *RedisLinkPool) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

func (p *RedisLinkPool) GetLink() (*RedisLink, error) {
	if p.Pool != nil {
		c := p.Pool.Get()
		if c != nil {
			return &RedisLink{c}, nil
		}
	}
	return nil, errors.New("Get Redis Link Fail")
}

type RedisLink struct {
	redis.Conn
}

func (p *RedisLink) Close() {
	p.RecordQueryTime()
	if p.Conn != nil {
		p.Conn.Close()
		p.Conn = nil
	}
}

func (p *RedisLink) RecordQueryTime() {
	// p.lastUseTime = time.Now()
}

func (p *RedisLink) Flush() error {
	return p.Conn.Flush()
}

// Key（键）
func (p *RedisLink) Del(key string) error {
	return p.Conn.Send("DEL", key)
}

func (p *RedisLink) Dump(key string) (interface{}, error) {
	return p.Conn.Do("DUMP", key)
}

func (p *RedisLink) Exists(key string) (interface{}, error) {
	return p.Conn.Do("EXISTS", key)
}

func (p *RedisLink) Expire(key string, seconds uint32) error {
	return p.Conn.Send("EXPIRE", key, seconds)
}

func (p *RedisLink) ExpireAt(key string, timestamp int64) error {
	return p.Conn.Send("EXPIREAT", key, timestamp)
}

func (p *RedisLink) Keys(pattern string) (interface{}, error) {
	return p.Conn.Do("KEYS", pattern)
}

// func (p *RedisLink) Migrate(host, port, key, destination, timeout interface{}) (interface{}, error) {
// 	return p.Conn.Do("MIGRATE", host, port, key, destination, timeout)
// }

// func (p *RedisLink) Move(key string, db interface{}) (interface{}, error) {
// 	return p.Conn.Do("MOVE", key, db)
// }

// func (p *RedisLink) Object(subCommand string, arguments interface{}) (interface{}, error) {
//
// }

func (p *RedisLink) Presist(key string) error {
	return p.Conn.Send("PERSIST", key)
}

func (p *RedisLink) PExpire(key string, milliseconds uint64) error {
	return p.Conn.Send("PEXPIRE", key, milliseconds)
}

func (p *RedisLink) PExpireAt(key string, millisecondsTimestamp uint64) error {
	return p.Conn.Send("PEXPIREAT", key, millisecondsTimestamp)
}

func (p *RedisLink) PTtl(key string) error {
	return p.Conn.Send("PTTL", key)
}

func (p *RedisLink) RandomKey() (interface{}, error) {
	return p.Conn.Do("RANDOMKEY")
}

func (p *RedisLink) Rename(key, newKey string) (interface{}, error) {
	return p.Conn.Do("RENAME", key, newKey)
}

func (p *RedisLink) Renamenx(key, newKey string) (interface{}, error) {
	return p.Conn.Do("RENAMENX", key, newKey)
}

func (p *RedisLink) Restore(key string, ttl, serializedValue interface{}) error {
	return p.Conn.Send("RESTORE", key, ttl, serializedValue)
}

// func (p *RedisLink) Sort(key string) (interface{}, error) {
//
// }

func (p *RedisLink) Ttl(key string) (interface{}, error) {
	return p.Conn.Do("TTL", key)
}

func (p *RedisLink) Type(key string) (interface{}, error) {
	return p.Conn.Do("TYPE", key)
}

// func (p *RedisLink) Scan(cursor, pattern, count interface{}) (interface{}, error) {
// 	return p.Conn.Do("SCAN", cursor, pattern, count)
// }

// String（字符串）/////////////////////////////////////////////////////////////////////////////////////////
func (p *RedisLink) Append(key string, value interface{}) error {
	return p.Conn.Send("APPEND", key, value)
}

func (p *RedisLink) BitCount(key string, start, end int) (interface{}, error) {
	return p.Conn.Do("BITCOUNT", key, start, end)
}

// func (p *RedisLink) BitOp(operation string, key string) (interface{}, error) {
// 	return p.Conn.Do("BITOP", operation, key)
// }

func (p *RedisLink) Decr(key string) (interface{}, error) {
	return p.Conn.Do("DECR", key)
}

func (p *RedisLink) DecrBy(key string, decrement int) error {
	return p.Conn.Send("DECRBY", key, decrement)
}

func (p *RedisLink) GetBit(key string, offset int) (interface{}, error) {
	return p.Conn.Do("GETBIT", key, offset)
}

func (p *RedisLink) Get(key string) (interface{}, error) {
	return p.Conn.Do("GET", key)
}

func (p *RedisLink) GetRange(key string, start, end int) (interface{}, error) {
	return p.Conn.Do("GETRANGE", key, start, end)
}

func (p *RedisLink) GetSet(key string, value interface{}) (interface{}, error) {
	return p.Conn.Do("GETSET", key, value)
}

func (p *RedisLink) Incr(key string) (interface{}, error) {
	return p.Conn.Do("INCR", key)
}

func (p *RedisLink) IncrBy(key string, increment int) (interface{}, error) {
	return p.Conn.Do("INCRBY", key, increment)
}

func (p *RedisLink) IncrByFloat(key string, increment float64) error {
	return p.Conn.Send("INCRBYFLOAT", key, increment)
}

func (p *RedisLink) MGet(key string) (interface{}, error) {
	return p.Conn.Do("MGET", key)
}

func (p *RedisLink) MSet(key string, value interface{}) error {
	return p.Conn.Send("MSET", key, value)
}

func (p *RedisLink) MSetnx(key string, value interface{}) error {
	return p.Conn.Send("MSETNX", key, value)
}

func (p *RedisLink) PSetex(key string, milliseconds int64, value interface{}) error {
	return p.Conn.Send("PSETEX", key, milliseconds, value)
}

func (p *RedisLink) Set(key string, value interface{}) error {
	return p.Conn.Send("SET", key, value)
}

func (p *RedisLink) Lock(key string, value interface{}, EX int, PX int, NX bool, XX bool) bool {
	args := append([]interface{}{}, key, value)
	if EX > 0 {
		args = append(args, "EX", EX) // 过期时间设置为 seconds 秒
	}
	if PX > 0 {
		args = append(args, "PX", PX) // 过期时间设置为 milliseconds 毫秒
	}
	if NX {
		args = append(args, "NX") // 只在键不存在时， 才对键进行设置操作
	}
	if XX {
		args = append(args, "XX") // 只在键已经存在时， 才对键进行设置操作
	}
	res, _ := redis.String(p.Conn.Do("SET", args...))
	return res == "OK"
}

func (p *RedisLink) Unlock(key ...interface{}) bool {
	res, _ := redis.Int(p.Conn.Do("DEL", key...))
	return res == len(key)
}

func (p *RedisLink) SetBit(key string, offset int, value interface{}) error {
	return p.Conn.Send("SETBIT", key, offset, value)
}

func (p *RedisLink) Setex(key string, seconds int64, value interface{}) error {
	return p.Conn.Send("SETEX", key, seconds, value)
}

func (p *RedisLink) Setnx(key string, value interface{}) error {
	return p.Conn.Send("SETNX", key, value)
}

func (p *RedisLink) SetRange(key string, offset int, value interface{}) error {
	return p.Conn.Send("SETRANGE", key, offset, value)
}

func (p *RedisLink) StrLen(key string) (interface{}, error) {
	return p.Conn.Do("STRLEN", key)
}

// SortedSet（有序集合）/////////////////////////////////////////////////////////////////////////////////////
func (p *RedisLink) ZAdd(key string, score interface{}, member interface{}) (interface{}, error) {
	return p.Conn.Do("ZADD", key, score, member)
}

func (p *RedisLink) ZCard(key string) (interface{}, error) {
	return p.Conn.Do("ZCARD", key)
}

func (p *RedisLink) ZCount(key string, min, max int) (interface{}, error) {
	return p.Conn.Do("ZCOUNT", key, min, max)
}

func (p *RedisLink) ZIncrBy(key string, increment int64, member interface{}) (interface{}, error) {
	return p.Conn.Do("ZINCRBY", key, increment, member)
}

func (p *RedisLink) ZRange(key string, start, stop int, bWithScores bool) (interface{}, error) {
	if bWithScores {
		return p.Conn.Do("ZRANGE", key, start, stop, "WITHSCORES")
	}
	return p.Conn.Do("ZRANGE", key, start, stop)
}

func (p *RedisLink) ZRangeByScore(key string, min, max int, bWithScores bool) (interface{}, error) {
	if bWithScores {
		return p.Conn.Do("ZRANGEBYSCORE", key, min, max, "WITHSCORES")
	}
	return p.Conn.Do("ZRANGEBYSCORE", key, min, max)
}

func (p *RedisLink) ZRank(key string, member interface{}) (interface{}, error) {
	return p.Conn.Do("ZRANK", key, member)
}

func (p *RedisLink) ZRem(key string, member interface{}) (interface{}, error) {
	return p.Conn.Do("ZREM", key, member)
}

func (p *RedisLink) ZRemRangeByRank(key string, start, stop int) (interface{}, error) {
	return p.Conn.Do("ZREMRANGEBYRANK", key, start, stop)
}

func (p *RedisLink) ZRemRangeByScore(key string, min, max int) (interface{}, error) {
	return p.Conn.Do("ZREMRANGEBYSCORE", key, min, max)
}

func (p *RedisLink) ZRevRange(key string, start, stop int, bWithScores bool) (interface{}, error) {
	if bWithScores {
		return p.Conn.Do("ZREVRANGE", key, start, stop, "WITHSCORES")
	}
	return p.Conn.Do("ZREVRANGE", key, start, stop)
}

func (p *RedisLink) ZRevRangeByScore(key string, max, min int, bWithScores bool) (interface{}, error) {
	if bWithScores {
		return p.Conn.Do("ZREVRANGEBYSCORE", key, max, min, "WITHSCORES")
	}
	return p.Conn.Do("ZREVRANGEBYSCORE", key, max, min)
}

func (p *RedisLink) ZRevRank(key string, member interface{}) (interface{}, error) {
	return p.Conn.Do("ZREVRANK", key, member)
}

func (p *RedisLink) ZScore(key string, member interface{}) (interface{}, error) {
	return p.Conn.Do("ZSCORE", key, member)
}

// func (p *RedisLink) ZUnionStore(keys []string, weights []interface{}, aggregate interface{}) (interface{}, error) {
//  return p.Conn.Do("ZUNIONSTORE", keys...)
// }

// func (p *RedisLink) ZInterStore(keys []string, weights []interface{}, aggregate interface{}) (interface{}, error) {
// 	return p.Conn.Do("ZINTERSTORE", keys...)
// }

// func (p *RedisLink) ZScan(key string, cursor, pattern, count interface{}) (interface{}, error) {
// 	return p.Conn.Do("ZSCAN", key, cursor, pattern, count)
// }

// Hash（哈希表）
func (p *RedisLink) HDel(key string, field ...interface{}) (interface{}, error) {
	args := make([]interface{}, 0, len(field)+1)
	args = append(args, key)
	args = append(args, field...)
	return p.Conn.Do("HDEL", args...)
}

func (p *RedisLink) HExists(key string, field interface{}) (interface{}, error) {
	return p.Conn.Do("HEXISTS", key, field)
}

func (p *RedisLink) HGet(key string, field interface{}) (interface{}, error) {
	return p.Conn.Do("HGET", key, field)
}

func (p *RedisLink) HGetAll(key string) (interface{}, error) {
	return p.Conn.Do("HGETALL", key)
}

func (p *RedisLink) HIncrBy(key string, field interface{}, increment int) (interface{}, error) {
	return p.Conn.Do("HINCRBY", key, field, increment)
}

func (p *RedisLink) HIncrByFloat(key string, field interface{}, increment float64) (interface{}, error) {
	return p.Conn.Do("HINCRBYFLOAT", key, field, increment)
}

func (p *RedisLink) HKeys(key string) (interface{}, error) {
	return p.Conn.Do("HKEYS", key)
}

func (p *RedisLink) HLen(key string) (interface{}, error) {
	return p.Conn.Do("HLEN", key)
}

func (p *RedisLink) HMGet(key string, field ...interface{}) (interface{}, error) {
	field = append([]interface{}{key}, field...)
	return p.Conn.Do("HMGET", field...)
}

func (p *RedisLink) HMSet(key string, kvs ...interface{}) (interface{}, error) {
	tot := make([]interface{}, 0, 1+len(kvs))
	tot = append(tot, key)
	tot = append(tot, kvs...)
	return p.Conn.Do("HMSET", tot...)
}

func (p *RedisLink) HSet(key string, field, value interface{}) (interface{}, error) {
	return p.Conn.Do("HSET", key, field, value)
}

func (p *RedisLink) HSetnx(key string, field, value interface{}) (interface{}, error) {
	return p.Conn.Do("HSETNX", key, field, value)
}

func (p *RedisLink) HVals(key string) (interface{}, error) {
	return p.Conn.Do("HVALS", key)
}

// func (p *RedisLink) HScan(key string, cursor, pattern, count interface{}) (interface{}, error) {
// 	return p.Conn.Do("HSCAN", key, cursor, pattern, count)
// }

// 列表操作 ///////////////////////////////////
func (p *RedisLink) RPush(key string, values ...interface{}) (interface{}, error) {
	args := []interface{}{key}
	args = append(args, values...)
	return p.Conn.Do("RPUSH", args...)
}

func (p *RedisLink) LPush(key string, values ...interface{}) (interface{}, error) {
	args := []interface{}{key}
	args = append(args, values...)
	return p.Conn.Do("LPUSH", args...)
}

func (p *RedisLink) LPop(key string) (interface{}, error) {
	return p.Conn.Do(key)
}

func (p *RedisLink) LRange(key string, start, end int) (interface{}, error) {
	return p.Conn.Do("LRANGE", key, start, end)
}

func (p *RedisLink) LRem(key string, count int, value interface{}) (interface{}, error) {
	return p.Conn.Do("LREM", key, count, value)
}

// 集合
// 将一个或多个 member 元素加入到集合 key 当中，已经存在于集合的 member 元素将被忽略
func (p *RedisLink) SAdd(key string, members ...interface{}) (interface{}, error) {
	args := []interface{}{key}
	args = append(args, members...)
	return p.Conn.Do("SADD", args...)
}

func (p *RedisLink) SRem(key string, members ...interface{}) (interface{}, error) {
	args := make([]interface{}, len(members)+1)
	args = append(args, key)
	args = append(args, members...)
	return p.Conn.Do("SREM", args...)
}

func (p *RedisLink) SIsMember(key string, member interface{}) (bool, error) {
	v, err := p.Conn.Do("SISMEMBER", key, member)
	if err != nil {
		return false, err
	}
	val := libs.Interface2Int(v)
	return val == 1, nil
}

// 发布订阅消息
func (p *RedisLink) Publish(channel interface{}, message interface{}) (interface{}, error) {
	return p.Conn.Do("Publish", channel, message)
}

// 订阅发布

// ///////////////////////////////////////////////////////
type IRedisCtrl interface {
	AddCfg(configs map[string]*svrcfg.RedisCfg)
	GetRdsLink(rdsName string) (*RedisLink, error)
	CloseRdsLink(rdsName string)
	CloseAllRdsLink()
}

// ///////////////////////////////////////////////////////

type RedisCtrl struct {
	// sync.Mutex

	exitChan  chan bool
	bClosed   bool
	onceClose sync.Once
	RdsPool   map[string]*RedisLinkPool
	ticker    *time.Ticker
}

func NewRedisCtrl() *RedisCtrl {
	t := &RedisCtrl{
		exitChan: make(chan bool, 1),
		RdsPool:  make(map[string]*RedisLinkPool, 5),
		ticker:   time.NewTicker(time.Minute * 5),
	}
	// go t.wait()
	return t
}

func (p *RedisCtrl) destroy() {
	p.onceClose.Do(func() {
		if p.exitChan != nil {
			close(p.exitChan)
			p.exitChan = nil
		}
		p.CloseAllRdsLink()
		p.bClosed = true
	})
}

// func (p *RedisCtrl) wait() {
// 	defer p.destroy()
//
// 	for {
// 		select {
// 		case <-p.ticker.C:
// 			p.onTicker()
// 		case <-p.exitChan:
// 			break
// 		}
// 	}
// }
//
// func (p *RedisCtrl) onTicker() {
// 	for _, h := range p.RdsMap {
// 		span := time.Since(h.lastUseTime)
// 		if span > time.Minute*5 {
// 			h.CloseLink()
// 		}
// 	}
// }

func (p *RedisCtrl) AddCfg(configs map[string]*svrcfg.RedisCfg) {
	for _, c := range configs {
		openConns := c.MaxOpenConns
		if openConns <= 0 {
			openConns = 10
		}
		idleConns := c.MaxIdleConns
		if idleConns <= 0 {
			idleConns = 1
		}
		idleTimeout := c.IdleTimeout
		if idleTimeout <= 0 {
			idleTimeout = 180
		}
		h := &RedisLinkPool{
			Cfg:         c,
			lastUseTime: time.Now(),
			Pool:        NewRedisPool(c.Name, c.GetDialType(), c.GetAddress(), openConns, idleConns, idleTimeout),
		}
		p.RdsPool[c.Name] = h
	}
}

func NewRedisPool(name, network, address string, maxOpen, maxIdle, idleTimeout int) *redis.Pool {
	pool := &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxOpen,
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
		// MaxConnLifetime: 0,
		Dial: func() (redis.Conn, error) {
			logs.Debug("Redis dial link. network:", network, " address:", address)
			c, err := redis.Dial(network, address)
			if err != nil {
				logs.Error("Dial Redis failure. err:", err.Error())
				return nil, err
			}

			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				logs.Error("TestOnBorrow ping err.", err)
			}
			return err
		},
	}

	logs.Info("Init Redis Pool success. name:", name, " network:", network, " address:", address)
	return pool
}

func (p *RedisCtrl) GetRdsLink(rdsName string) (*RedisLink, error) {
	if h, ok := p.RdsPool[rdsName]; ok {
		link, err := h.GetLink()
		if err != nil {
			return nil, err
		}
		return link, nil
	}
	return nil, errors.New("no redis cfg setted")
}

func (p *RedisCtrl) CloseRdsLink(rdsName string) {
	if h, ok := p.RdsPool[rdsName]; ok {
		h.Close()
	}
}

func (p *RedisCtrl) CloseAllRdsLink() {
	for _, h := range p.RdsPool {
		h.Close()
	}
	logs.Info("all redis closed.")
}

// ////////////////////  获取通用信息方法  ///////////////////////////////////

// 获取房间信息
func (p *RedisCtrl) GetRoomInfo(roomId int) (*db_ds.RoomInfo, error) {
	h, err := p.GetRdsLink(UserRdsName)
	if err != nil {
		return nil, err
	}
	defer h.Close()

	key := fmt.Sprintf("room_cache_roomid_%d", roomId)
	// rep, err := redis.String(h.Conn.Do("GET", key))
	rep, err := redis.String(h.Get(key))
	if err != nil {
		logs.Info("redis get room info failure. Error:", err.Error())
		return nil, err
	}

	h.RecordQueryTime()

	roomInfo := new(db_ds.RoomInfo)
	var js map[string]interface{}

	err = json.Unmarshal([]byte(rep), &js)
	if err != nil {
		logs.Info("redis room info parse json failure. Error:", err.Error())
		return nil, err
	}

	if v, ok := js["id"]; ok {
		roomInfo.Id = libs.ToInt(v)
	}
	if v, ok := js["nickname"]; ok {
		roomInfo.Nickname = v.(string)
	}
	if v, ok := js["code"]; ok {
		roomInfo.Code = libs.ToInt(v)
	}
	if v, ok := js["gameId"]; ok {
		roomInfo.GameId = libs.ToInt(v)
	}
	if v, ok := js["spic"]; ok {
		roomInfo.Spic = v.(string)
	}
	if v, ok := js["bpic"]; ok {
		roomInfo.Bpic = v.(string)
	}
	if v, ok := js["status"]; ok {
		roomInfo.Status = libs.ToInt(v)
	}
	if v, ok := js["liveTime"]; ok {
		roomInfo.LiveTime = int64(libs.ToInt(v))
	}
	if v, ok := js["publishUrl"]; ok {
		roomInfo.PublishUrl = v.(string)
	}
	if v, ok := js["videoIdKey"]; ok {
		roomInfo.VideoId = v.(string)
	}
	return roomInfo, nil
}

// 获取用户手机信息
func (p *RedisCtrl) GetUserCell(uid int) (string, error) {
	h, err := p.GetRdsLink(UserRdsName)
	if err != nil {
		return "", err
	}
	defer h.Close()

	key := fmt.Sprintf("userdata:mobile:%d", int(uid%5000))
	// cell, err := redis.String(h.Conn.Do("HGET", key, uid))
	cell, err := redis.String(h.HGet(key, uid))
	if err != nil {
		logs.Info("redis get user cell failure: ", "HGET", key, uid)
		return "", err
	}

	return cell, nil
}

func (p *RedisCtrl) GetUserEmail(uid int) (string, error) {
	h, err := p.GetRdsLink(UserRdsName)
	if err != nil {
		return "", err
	}
	defer h.Close()

	key := fmt.Sprintf("userdata:email:%d", int(uid%5000))
	// email, err := redis.String(h.Conn.Do("HGET", key, uid))
	email, err := redis.String(h.HGet(key, uid))
	if err != nil {
		return "", err
	}

	return email, nil
}

func (p *RedisCtrl) GetUserInfo(uid uint) (*db_ds.UserInfo, error) {
	h, err := p.GetRdsLink(UserRdsName)
	if err != nil {
		return nil, err
	}
	defer h.Close()

	// fileds := []string{"avatar","nickname","gender","account"}
	// logs.Info("RedisCtrl::GetUserInfo 1>>", uid)

	key := fmt.Sprintf("usersvr_user_%d", uid)
	// logs.Info("RedisCtrl::GetUserInfo 2>>", key)

	res, err := redis.Strings(h.Conn.Do("HMGET", key, "avatar", "nickname", "gender"))
	if err != nil {
		// logs.Info("RedisCtrl:GetUserInfo 3>>", err)
		return nil, err
	}
	// logs.Info("RedisCtrl:GetUserInfo 4>>", res)

	ll := len(res)
	if ll > 0 {
		user := new(db_ds.UserInfo)
		user.Uid = uint(uid)
		if ll >= 1 && len(res[0]) > 0 {
			user.Avatar = res[0]
		}
		if ll >= 2 && len(res[1]) > 0 {
			user.Nickname = res[1]
		}
		if ll >= 3 && len(res[2]) > 0 {
			v, _ := strconv.Atoi(res[2])
			user.Gender = uint8(v)
		}
		if ll >= 4 && len(res[3]) > 0 {
			user.Account = res[3]
		}
		return user, nil
	}
	return nil, errors.New("user info get fail")
}

func (p *RedisCtrl) CheckToken(r *http.Request) (*db_ds.UserInfo, error) {
	guidCookie, guidErr := r.Cookie("ZQ_GUID")
	if guidErr != nil {
		logs.Info("check token ZQ_GUID_C cookie sid. Error:", guidErr.Error())
		return nil, guidErr
	}

	_, sidErr := r.Cookie("PHPSESSID") // sidCookie
	if sidErr != nil {
		logs.Info("check token no PHPSESSID cookie sid. Error:", sidErr.Error())
		return nil, sidErr
	}

	var (
		body string
		ok   bool
	)

	// 优先从redis获取用户数据
	if body, ok = p.GetUserInfoBodyByCookie(guidCookie.Value); ok {
		// logs.Info("GetUserInfoBodyByCookie::body=", body)

		user := &db_ds.UserInfo{}
		dataObj, err := simplejson.NewJson([]byte(body))
		if err != nil {
			logs.Info("check token parse js result failure. Error:", err.Error(), "body:", string(body))
			return nil, err
		}
		user.FromSimpleJson(dataObj)
		return user, nil
	} else {
		logs.Info("GetUserInfoBodyByCookie:: fail")
	}

	// 直接返回失败，不再http查询
	return nil, errors.New("check token fail")
}

func (p *RedisCtrl) GetUserInfoBodyByCookie(GUID string) (string, bool) {
	c, err := p.GetRdsLink(SessionRdsName)
	if err != nil {
		return "", false
	}

	if c == nil {
		logs.Error("getUserInfoByRedis, get redis connect failure. GUID:", GUID)
		return "", false
	}

	defer c.Close()

	keyGUID := fmt.Sprintf("apiauth_session_guid_%s", GUID)
	// uid, err := redis.Int(c.Do("get", keyGUID))
	uid, err := redis.Int(c.Get(keyGUID))
	if err != nil {
		logs.Error("getUserInfoByRedis get session guid failure. GUID:", GUID, " err:", err.Error())
		return "", false
	}

	if uid <= 0 {
		return "", false
	}

	keyUID := fmt.Sprintf("apiauth_session_auth_%d", uid)
	// userInfoJson, err := redis.String(c.Do("get", keyUID))
	userInfoJson, err := redis.String(c.Get(keyUID))
	if err != nil {
		logs.Error("getUserInfoByRedis get session auth failure. GUID:", GUID, " uid:", uid, " err:", err.Error())
		return "", false
	}

	return userInfoJson, true
}
