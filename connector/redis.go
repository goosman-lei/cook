package connector

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	cook_log "gitlab.niceprivate.com/golang/cook/log"
	cook_util "gitlab.niceprivate.com/golang/cook/util"
	"math/rand"
	"time"
)

type RedisConnWrapper struct {
	redis.Conn
}

type RedisConf struct {
	Addrs []string

	TestInterval time.Duration

	MaxActive   int
	MaxIdle     int
	IdleTimeout time.Duration

	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

var redisConnMapping *cook_util.CMap = cook_util.NewCMap()

func setup_one_redis(sn string, config RedisConf) {
	redisConnMapping.Set(sn, &redis.Pool{
		MaxActive:   config.MaxActive,
		MaxIdle:     config.MaxIdle,
		IdleTimeout: config.IdleTimeout,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			var (
				addr string
				conn redis.Conn
				err  error
			)
			addr = config.Addrs[rand.Intn(len(config.Addrs))]
			conn, err = redis.DialTimeout("tcp", addr, config.ConnectTimeout, config.ReadTimeout, config.WriteTimeout)
			if err != nil {
				cook_log.Warnf("connect to redis[%s] failed: %s", addr, err)
				return nil, err
			}

			_, err = conn.Do("PING")
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < config.TestInterval {
				return nil
			}
			_, err := conn.Do("PING")
			return err
		},
	})
}

func SetupRedis(configs map[string]RedisConf) error {
	for sn, config := range configs {
		setup_one_redis(sn, config)
	}
	return nil
}

func GetRedis(sn string) (*RedisConnWrapper, error) {
	if conn, exists := redisConnMapping.Get(sn); exists {
		return &RedisConnWrapper{conn.(*redis.Pool).Get()}, nil
	}
	cook_log.Warnf("get redis conn[%s], but not ready", sn)
	return nil, fmt.Errorf("have no mysql cluster: %s", sn)
}

func MustGetRedis(sn string) *RedisConnWrapper {
	conn, err := GetRedis(sn)
	if err != nil {
		panic(err)
	}
	return conn
}

func (c *RedisConnWrapper) DoBool(command string, argv ...interface{}) (bool, error) {
	return redis.Bool(c.Do(command, argv...))
}
func (c *RedisConnWrapper) DoByteSlices(command string, argv ...interface{}) ([][]byte, error) {
	return redis.ByteSlices(c.Do(command, argv...))
}
func (c *RedisConnWrapper) DoBytes(command string, argv ...interface{}) ([]byte, error) {
	return redis.Bytes(c.Do(command, argv...))
}
func (c *RedisConnWrapper) DoFloat64(command string, argv ...interface{}) (float64, error) {
	return redis.Float64(c.Do(command, argv...))
}
func (c *RedisConnWrapper) DoInt(command string, argv ...interface{}) (int, error) {
	return redis.Int(c.Do(command, argv...))
}
func (c *RedisConnWrapper) DoInt64(command string, argv ...interface{}) (int64, error) {
	return redis.Int64(c.Do(command, argv...))
}
func (c *RedisConnWrapper) DoInt64Map(command string, argv ...interface{}) (map[string]int64, error) {
	return redis.Int64Map(c.Do(command, argv...))
}
func (c *RedisConnWrapper) DoIntMap(command string, argv ...interface{}) (map[string]int, error) {
	return redis.IntMap(c.Do(command, argv...))
}
func (c *RedisConnWrapper) DoInts(command string, argv ...interface{}) ([]int, error) {
	return redis.Ints(c.Do(command, argv...))
}
func (c *RedisConnWrapper) DoMultiBulk(command string, argv ...interface{}) ([]interface{}, error) {
	return redis.MultiBulk(c.Do(command, argv...))
}
func (c *RedisConnWrapper) DoPositions(command string, argv ...interface{}) ([]*[2]float64, error) {
	return redis.Positions(c.Do(command, argv...))
}
func (c *RedisConnWrapper) DoString(command string, argv ...interface{}) (string, error) {
	return redis.String(c.Do(command, argv...))
}
func (c *RedisConnWrapper) DoStringMap(command string, argv ...interface{}) (map[string]string, error) {
	return redis.StringMap(c.Do(command, argv...))
}
func (c *RedisConnWrapper) DoStrings(command string, argv ...interface{}) ([]string, error) {
	return redis.Strings(c.Do(command, argv...))
}
func (c *RedisConnWrapper) DoUint64(command string, argv ...interface{}) (uint64, error) {
	return redis.Uint64(c.Do(command, argv...))
}
func (c *RedisConnWrapper) DoValues(command string, argv ...interface{}) ([]interface{}, error) {
	return redis.Values(c.Do(command, argv...))
}
