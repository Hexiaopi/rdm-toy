package model

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/redis/go-redis/v9"
)

type Conn struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	UserName string `json:"username"`
	PassWord string `json:"password"`
}

func (c Conn) NewClient() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", c.Host, c.Port),
		Username: c.UserName,
		Password: c.PassWord,
	})
	rdb.AddHook(RedisLog{})
	result, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	if result != "PONG" {
		return nil, errors.New("wrong client")
	}
	return rdb, nil
}

type RedisLog struct {
}

func (rl RedisLog) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

func (rl RedisLog) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		next(ctx, cmd)
		log.Println("redis", cmd.Args())
		return nil
	}
}

func (rl RedisLog) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		return next(ctx, cmds)
	}
}
