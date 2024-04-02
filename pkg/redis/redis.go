package redis

import (
	"demogo/config"
	"demogo/pkg/constant"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/go-redis/redis"
	"github.com/go-redsync/redsync/v4"
	redigoRedsync "github.com/go-redsync/redsync/v4/redis/redigo"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/nitishm/go-rejson"
)

const (
	defaultMaxIdle         int = 1
	defaultMaxActive       int = 5
	defaultMaxConnLifetime int = 240 // default from redigo
	ErrMutexAlreadyExist       = "mutex with same key already exist"
)

// RedisHandler is a struct for save redis Client and a receiver
type RedisHandler struct {
	Client      *redis.Client
	Pipe        redis.Pipeliner
	PipeCounter int
	RH          *rejson.Handler
	RS          map[string]*redisearch.Client
	Pool        *redigo.Pool
	Redisync    *redsync.Redsync
}

// ConnectRedis make redis connection
func (r *RedisHandler) ConnectRedis(ra *config.RedisAccount) {

	r.Client = redis.NewClient(&redis.Options{
		Addr:         ra.Host + ":" + strconv.Itoa(ra.Port),
		Password:     ra.Password,
		DB:           ra.DB,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		PoolSize:     ra.PoolSize,
		MinIdleConns: ra.MinIdleConns,
		PoolTimeout:  16 * time.Second,
		MaxRetries:   2,
		// MaxConnAge:         30 * time.Second,
		IdleCheckFrequency: 1 * time.Minute,
	})
	result, err := r.Client.Ping().Result()
	if err != nil {
		log.Fatalf("Could not connect redis, error %s", err.Error())
	}
	fmt.Printf("\n redis handler redisearch ping result %s \n", result)

	log.Println("Connected to redis")
	rh := rejson.NewReJSONHandler()

	flag.Parse()
	rh.SetGoRedisClient(r.Client)
	r.RH = rh

	fmt.Printf("\n redis handler redisearch index %+v \n", ra.RedisearchIndex)
	r.buildPool(ra)
	r.buildRedisHandler(ra, ra.RedisearchIndex)
}

// Get data
func (r *RedisHandler) Get(key string) (string, error) {
	val, err := r.Client.Get(key).Result()
	if err != nil {
		return "", err
	}
	return val, err
}

// Get data with pattern
// note : avoid pattern '*' for keys lookup. Consider adding prefix and/or suffix on the pattern.
func (r *RedisHandler) Keys(pattern string) ([]string, error) {
	val, err := r.Client.Keys(pattern).Result()
	if err != nil {
		return nil, err
	}
	return val, err
}

// HGetAll get hash redis data all fields
func (r *RedisHandler) HGetAll(key string) (map[string]string, error) {
	val, err := r.Client.HGetAll(key).Result()
	if err != nil {
		return make(map[string]string), err
	}
	return val, nil
}

// Set insert data to redis in string type
func (r *RedisHandler) Set(key string, data interface{}) error {
	statusCmd := r.Client.Set(key, data, 0)

	if statusCmd != nil {
		err := statusCmd.Err()
		if err != nil {
			return err
		}
	}

	return nil
}

// SetEx insert data to old redis in string type with Expiration (TTL)
func (r *RedisHandler) SetEx(key string, parameter SetExParameter) error {
	if parameter.IsTesting {
		return nil
	}
	statusCmd := r.Client.Set(key, parameter.Data, parameter.ExpireDuration)

	if statusCmd != nil {
		err := statusCmd.Err()
		if err != nil {
			return err
		}
	}
	return nil
}

// GetSet get and insert data to redis
func (r *RedisHandler) GetSet(key string, data interface{}) error {
	err := r.Client.GetSet(key, data).Err()
	if err != nil {
		return err
	}
	return nil
}

// HMSet insert hash data to redis type hash
func (r *RedisHandler) HMSet(key string, data interface{}) error {
	m := make(map[string]interface{})
	destVal := reflect.ValueOf(data)
	dataByte := destVal.Bytes()
	err := json.Unmarshal(dataByte, &m)
	if err != nil {
		return err
	}
	err = r.Client.HMSet(key, m).Err()
	if err != nil {
		return err
	}
	return nil
}

// HMSet insert hash data to redis type hash with Expiration (TTL)
func (r *RedisHandler) HMSetEx(key string, parameter SetExParameter) error {
	m := make(map[string]interface{})
	destVal := reflect.ValueOf(parameter.Data)
	dataByte := destVal.Bytes()
	err := json.Unmarshal(dataByte, &m)
	if err != nil {
		return err
	}
	err = r.Client.HMSet(key, m).Err()
	if err != nil {
		return err
	}

	statusCmd := r.Client.Expire(key, parameter.ExpireDuration)
	if statusCmd != nil {
		err := statusCmd.Err()
		if err != nil {
			return err
		}
	}
	return nil
}

// HMGet get hash data from redis
func (r *RedisHandler) HMGet(key string, fields ...string) ([]interface{}, error) {
	val, err := r.Client.HMGet(key, fields...).Result()
	return val, err
}

// Exists check exists data from redis
func (r *RedisHandler) Exists(key string) (bool, error) {
	checker, err := r.Client.Exists(key).Result()
	if err != nil {
		return false, err
	}

	return checker == 1, nil
}

// Close redis connection
func (r *RedisHandler) Close() {
	err := r.Client.Close()
	if err != nil {
		log.Println(err.Error())
	}
}

// Scan accept wildcard for searching ex : users*
func (r *RedisHandler) Scan(key string) ([]string, error) {
	c1 := uint64(10)
	n := int64(1)
	c := uint64(0)
	var res []string
	var res1 []string
	var err error
	for c1 != 0 {
		n = n * 10
		res1, c, err = r.Client.Scan(c, key, n).Result()
		if err != nil {
			return res, err
		}
		res = append(res, res1...)
		c1 = c
	}
	if len(res) <= 0 {
		return res, errors.New(constant.NotFound)
	}
	return res, nil
}

// Del data in redis from keys (one, or more than one key)
func (r *RedisHandler) Del(keys ...string) error {
	for _, key := range keys {
		_, err := r.Client.Del(key).Result()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RedisHandler) HDel(key string, fields ...string) error {
	_, err := r.Client.HDel(key, fields...).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisHandler) XDel(stream string, ids ...string) error {
	_, err := r.Client.XDel(stream, ids...).Result()
	if err != nil {
		return err
	}
	return nil
}

// LLen func to get total list
func (r *RedisHandler) LLen(key string) (int64, error) {
	total, err := r.Client.LLen(key).Result()
	if err != nil {
		return total, err
	}
	return total, nil
}

// LRange func to get value from list by range
func (r *RedisHandler) LRange(key string, start int64, stop int64) ([]string, error) {
	data, err := r.Client.LRange(key, start, stop).Result()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// LPush func to push into list
func (r *RedisHandler) LPush(key string, value ...interface{}) error {
	for _, val := range value {
		_, err := r.Client.LPush(key, val).Result()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RedisHandler) SPop(key string) (string, error) {
	result, err := r.Client.SPop(key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

// SAdd add members to a set stored in redis to a key. if key not exists, will create new and add member
func (r *RedisHandler) SAdd(key string, data ...interface{}) error {
	_, err := r.Client.SAdd(key, data...).Result()
	if err != nil {
		return err
	}
	return nil
}

// SIsMember will return true if data is member of set stored at a key or false if data is not member of set stored at a key
func (r *RedisHandler) SIsMember(key string, data interface{}) (bool, error) {
	exist, err := r.Client.SIsMember(key, data).Result()
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (r *RedisHandler) Pipeline() {
	r.Pipe = r.Client.Pipeline()
	r.PipeCounter = 0
}

func (r *RedisHandler) JSONSet(key string, data interface{}) (interface{}, error) {
	return r.RH.JSONSet(key, ".", data)
}

func (r *RedisHandler) JSONGet(key string) (interface{}, error) {
	return r.RH.JSONGet(key, ".")
}

func (r *RedisHandler) JSONMGet(path string, keys ...string) (interface{}, error) {
	return r.RH.JSONMGet(path, keys...)
}
func (r *RedisHandler) JSONDel(key, path string) (interface{}, error) {
	return r.RH.JSONDel(key, path)
}

func (r *RedisHandler) CreateRedisync() *redsync.Redsync {

	if r.Redisync == nil {

		// Create a pool with go-redis which is the pool redisync will use while communicating with redis
		pool := redigoRedsync.NewPool(r.Pool)

		// Create an instance of redisync to be used to obtain a mutual exclusion lock.
		rs := redsync.New(pool)
		r.Redisync = rs
	}

	return r.Redisync
}

func (r *RedisHandler) CreateRedisMutex(key string, options ...redsync.Option) *redsync.Mutex {
	redisync := r.Redisync

	if redisync == nil {
		redisync = r.CreateRedisync()
	}

	mutex := redisync.NewMutex(key, options...)

	return mutex
}

func (r *RedisHandler) LockRedisMutex(mutex *redsync.Mutex, tries int) error {
	var err error

	if mutex != nil {
		for i := 0; i < tries; i++ {
			err = mutex.Lock()
			if err == nil || err == redsync.ErrFailed {
				break
			}
		}
	}

	return err
}

func (r *RedisHandler) buildPool(ra *config.RedisAccount) {
	formatClient := fmt.Sprintf("%s:%d", ra.Host, ra.Port)

	maxIdle := ra.MaxIdle
	if maxIdle < 1 {
		maxIdle = defaultMaxIdle
	}

	maxActive := ra.MaxActive
	if maxActive < 1 {
		maxActive = defaultMaxActive
	}
	maxConnLifetime := ra.MaxConnLifetime
	if maxConnLifetime < 1 {
		maxConnLifetime = defaultMaxConnLifetime
	}

	pool := &redigo.Pool{
		MaxIdle:         maxIdle,
		MaxActive:       maxActive,
		MaxConnLifetime: time.Duration(maxConnLifetime) * time.Second,
		Wait:            true,
		TestOnBorrow: func(c redigo.Conn, t time.Time) (err error) {
			if time.Since(t) > time.Second {
				_, err = c.Do("PING")
			}

			return err
		},
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", formatClient)
			if err != nil {
				return nil, err
			}
			if ra.Password != "" {
				if _, err := c.Do("AUTH", ra.Password); err != nil {
					c.Close()
					return nil, err
				}
			}

			// if _, err := c.Do("SELECT", 0); err != nil {
			// 	c.Close()
			// 	return nil, err
			// }
			return c, nil
		},
	}

	r.Pool = pool
}

func (r *RedisHandler) buildRedisHandler(ra *config.RedisAccount, idx []string) {
	m := make(map[string]*redisearch.Client)
	for _, v := range idx {
		// m[v] = redisearch.NewClient(ra.URL+":"+strconv.Itoa(ra.Port), v)
		m[v] = redisearch.NewClientFromPool(r.Pool, v)
		res, err := m[v].Info()
		if err != nil {
			log.Println(err.Error())
		} else {
			// check can be marshal or not
			_, err = json.Marshal(res)
			if err != nil {
				log.Println(err.Error())
			} else {
				log.Println(res)
			}
		}

		r.RS = m
	}
}

func (r *RedisHandler) NewMutex(key string, timeout time.Duration) (*redsync.Mutex, error) {
	cache, err := r.Get(key)
	if err != nil && !strings.EqualFold(err.Error(), redis.Nil.Error()) {
		return nil, err
	}

	if !strings.EqualFold(cache, "") {
		return nil, fmt.Errorf(ErrMutexAlreadyExist)
	}

	return r.CreateRedisMutex(key, redsync.WithExpiry(timeout)), nil
}

func (r *RedisHandler) UnlockMutex(mutex *redsync.Mutex) (bool, error) {
	return mutex.Unlock()
}

func (r *RedisHandler) CreateAndLockMutex(key string, timeout time.Duration) (*redsync.Mutex, error) {
	mutex, err := r.NewMutex(key, timeout)
	if err != nil {
		return nil, err
	}

	if err = r.LockRedisMutex(mutex, 3); err != nil {
		return nil, err
	}

	return mutex, nil
}
