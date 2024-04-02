package redis

//go:generate mockgen --destination=../mocks/mock_iredis.go --package=mocks --source=iredis.go
import (
	"github.com/arraisi/demogo/config"
	"github.com/go-redsync/redsync/v4"
	"time"
)

// IRedis is the general interface for redis commands usage
type IRedis interface {
	ConnectRedis(ra *config.RedisAccount)
	Get(key string) (string, error)
	Keys(pattern string) ([]string, error)
	HGetAll(key string) (map[string]string, error)
	Set(key string, data interface{}) error
	SetEx(key string, parameter SetExParameter) error
	HMSet(key string, data interface{}) error
	HMSetEx(key string, parameter SetExParameter) error
	GetSet(key string, data interface{}) error
	HMGet(key string, fields ...string) ([]interface{}, error)
	Exists(key string) (bool, error)
	Close()
	Scan(key string) ([]string, error)
	Del(keys ...string) error
	HDel(key string, fields ...string) error
	XDel(stream string, ids ...string) error
	LLen(key string) (int64, error)
	LRange(key string, start int64, stop int64) ([]string, error)
	LPush(key string, value ...interface{}) error
	SPop(key string) (string, error)
	SAdd(key string, data ...interface{}) error
	SIsMember(key string, data interface{}) (bool, error)
	Pipeline()
	JSONSet(key string, data interface{}) (interface{}, error)
	JSONGet(key string) (interface{}, error)
	JSONMGet(path string, keys ...string) (interface{}, error)
	JSONDel(key, path string) (interface{}, error)
	CreateRedisync() *redsync.Redsync
	CreateRedisMutex(key string, options ...redsync.Option) *redsync.Mutex
	LockRedisMutex(mutex *redsync.Mutex, tries int) error
	NewMutex(key string, timeout time.Duration) (*redsync.Mutex, error)
	UnlockMutex(mutex *redsync.Mutex) (bool, error)
	CreateAndLockMutex(key string, timeout time.Duration) (*redsync.Mutex, error)
}
