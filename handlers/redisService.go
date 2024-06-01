package handlers

import (
	"FeishuCozeRobot/config"
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"time"
)

type RedisUtil struct {
	client *redis.Client
}

// RedisClient  全局变量, 外部使用utils.RedisClient来访问
var RedisClient RedisUtil

// InitRedisUtil  初始化redis
func InitRedisUtil(redisConig config.RedisConfig) (*RedisUtil, error) {
	//连接redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisConig.Addr + ":" + strconv.Itoa(redisConig.Port),
		Password: redisConig.Password,
		DB:       redisConig.Db,
		PoolSize: 10,
	})

	//初始化全局redis结构体
	RedisClient = RedisUtil{client: client}
	return &RedisClient, nil
}

// SetStr 设置数据到redis中（string）
func (rs *RedisUtil) SetStr(ctx context.Context, key string, value string, expiration time.Duration) error {
	_, err := rs.client.Set(ctx, key, value, expiration).Result()
	return err
}

// SetStrNotExist 设置数据到redis中（string）
func (rs *RedisUtil) SetStrNotExist(ctx context.Context, key string, value string, expireSecond int) bool {
	val, err := rs.client.Do(ctx, "SET", key, value, "EX", expireSecond, "NX").Result()
	if err != nil || val == nil {
		return false
	}
	return true
}

// SetStrWithExpire 设置数据到redis中（string）
func (rs *RedisUtil) SetStrWithExpire(ctx context.Context, key string, value string, expireSecond int) error {
	err := rs.client.Do(ctx, "Set", key, value, "ex", expireSecond).Err()
	return err
}

// GetStr 获取redis中数据（string）
func (rs *RedisUtil) GetStr(ctx context.Context, key string) (string, error) {
	val, err := rs.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

// HSet 设置数据到redis中（hash）
func (rs *RedisUtil) HSet(ctx context.Context, key string, field string, value string) error {
	return rs.client.Do(ctx, "HSet", key, field, value).Err()
}

// HGet 获取redis中数据（hash）
func (rs *RedisUtil) HGet(ctx context.Context, key string, field string) (string, error) {
	val, err := rs.client.Do(ctx, "HGet", key, field).Result()
	if err != nil {
		return "", err
	}
	return string(val.([]byte)), nil
}

// DelByKey 删除
func (rs *RedisUtil) DelByKey(ctx context.Context, key string) error {
	return rs.client.Do(ctx, "DEL", key).Err()

}

// SetExpire 设置key过期时间
func (rs *RedisUtil) SetExpire(ctx context.Context, key string, expireSecond int) error {
	return rs.client.Do(ctx, "EXPIRE", key, expireSecond).Err()
}

// KEYEXISTS 判断KEY在redis中是否存在
func (rs *RedisUtil) KEYEXISTS(ctx context.Context, KEY string) bool {
	exists, err := rs.client.Do(ctx, "EXISTS", KEY).Bool()
	if err != nil {
		return false
	}
	return exists
}

// KEYEXISTSGetStr 判断KEY在redis中是否存在,存在则获取内容
func (rs *RedisUtil) KEYEXISTSGetStr(ctx context.Context, KEY string) (bool, string) {
	if rs.KEYEXISTS(ctx, KEY) {
		str, err := rs.GetStr(ctx, KEY)
		if err == nil {
			return true, str
		}
	}
	return false, ""
}

// GetBytes  获取redis中数据（string）
func (rs *RedisUtil) GetBytes(ctx context.Context, key string) ([]byte, error) {
	return rs.client.Get(ctx, key).Bytes()
}

// KEYEXISTSGetBytes  判断KEY在redis中是否存在,存在则获取内容
func (rs *RedisUtil) KEYEXISTSGetBytes(ctx context.Context, KEY string) (bool, []byte) {
	if rs.KEYEXISTS(ctx, KEY) {
		str, err := rs.GetBytes(ctx, KEY)
		if err == nil {
			return true, str
		}
	}
	return false, nil
}

// KEYEXISTSGetScan  判断KEY在redis中是否存在,存在则获取指定类型的内容
func (rs *RedisUtil) KEYEXISTSGetScan(ctx context.Context, KEY string, Val interface{}) bool {
	if rs.KEYEXISTS(ctx, KEY) {
		body, err := rs.client.Get(ctx, KEY).Bytes()
		if err == nil {
			if json.Unmarshal(body, Val) == nil {
				return true
			}
		}
		log.Println("获取redis 数据报错了", err.Error())
	}
	return false
}

// Set 设置数据到redis中 泛型
func (rs *RedisUtil) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rs.client.Set(ctx, key, value, expiration).Err()
}
