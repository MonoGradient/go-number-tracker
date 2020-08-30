package service

import (
	"Go-Tracker/model"
	"context"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

var ctx = context.Background()

func newRedisClient() *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB: 0,
	})

	return rdb
}

func IdempotentIncrement(opts ...string) (*model.StorageResponse, error) {
	executionTime := time.Now()
	var key string
	if len(opts) > 0 {
		key = opts[0]
		if _, err := uuid.Parse(key); err != nil {
			log.Infof("Unable to parse provided key of %v. Generating new key.\n", key)
			key = uuid.New().String()
			log.Infof("Generating new ID for key. New ID: %v", key)
		}
	} else {
		key = uuid.New().String()
		log.Infof("Generating new ID for key. New ID: %v", key)
	}
	client := newRedisClient()
	defer client.Close()
	if res,err := client.Incr(key).Result(); err != nil {
		log.Errorf("Error incrementing data: %v", err)
		return nil, err
	} else {
		return &model.StorageResponse{
			Key:             key,
			Value:           res,
			ActionTimestamp: executionTime.Format(time.RFC3339),
		}, nil
	}
}

func Decrement(key string) (*model.StorageResponse, error) {
	client := newRedisClient()
	defer client.Close()
	if _, err := client.Get(key).Result(); err != nil {
		return nil, err
	}
	result, err := client.Decr(key).Result()
	if err != nil {
		return nil, err
	}
	return &model.StorageResponse{
		Key:             key,
		Value:           result,
		ActionTimestamp: time.Now().Format(time.RFC3339),
	}, nil
}

func Delete(key string) bool {
	client := newRedisClient()
	defer client.Close()
	if _, err := client.Get(key).Result(); err != nil {
		log.Errorf("Failed to find key in DB: %v", err.Error())
		return false
	} else {
		client.Del(key)
		return true
	}
}

func Check(key string) (*model.StorageResponse, error) {
	client := newRedisClient()
	defer client.Close()
	if result, err := client.Get(key).Result(); err != nil {
		return nil, err
	} else {
		if value, err := strconv.ParseInt(result, 10, 64); err != nil {
			return nil, err
		} else {
			return &model.StorageResponse{
				Key:             key,
				Value:           value,
				ActionTimestamp: time.Now().Format(time.RFC3339),
			}, nil
		}
	}
}
