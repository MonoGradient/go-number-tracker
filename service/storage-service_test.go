package service

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

func init() {
	os.Setenv("REDIS_HOST", "localhost:6379")
	os.Setenv("REDIS_PASSWORD", "")
}

func BenchmarkIdempotentIncrement(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IdempotentIncrement()
	}
}


func TestIdempotentIncrement_CreateSuccess(t *testing.T) {
	client := newRedisClient()
	defer client.Close()
	value, err := IdempotentIncrement()
	if err != nil {
		assert.Error(t, err, "Failed to execute IdempotentIncrement!")
		t.Error(err)
	}


	if result, err := client.Get(value.Key).Result(); err != nil {
		assert.Error(t, err, "Failed to retrieve data from Redis with key")
	} else {
		convertedResult, _ := strconv.ParseInt(result, 10, 64)
		assert.Equal(t, convertedResult, value.Value, "Values from Redis should be equal")
	}
	t.Cleanup(cleanUpRedis)
}

func TestIdempotentIncrement_IncrementSuccess(t *testing.T) {
	client := newRedisClient()
	defer client.Close()

	key :=uuid.New().String()
	initialValue := rand.Int63()
	client.Set(key, initialValue, 15 * time.Minute)
	value, err := IdempotentIncrement(key)
	if err != nil {
		assert.Fail(t, "Failed to execute IdempotentIncrement!")
		t.Error(err)
	}


	if result, err := client.Get(value.Key).Result(); err != nil {
		assert.Fail(t, "Failed to retrieve data from Redis with key")
	} else {
		convertedResult, _ := strconv.ParseInt(result, 10, 64)
		assert.Equal(t, value.Value, initialValue+1, "Value should match initial Value + 1")
		assert.Equal(t, convertedResult, value.Value, "Values from Redis should be equal")
	}
	t.Cleanup(cleanUpRedis)
}

func TestIdempotentIncrement_CreateSuccess_InvalidKey(t *testing.T) {
	client := newRedisClient()
	defer client.Close()
	value, err := IdempotentIncrement("NOTAVALIDKEY")
	if err != nil {
		assert.Fail(t, "Failed to execute IdempotentIncrement!")
		t.Error(err)
	}


	if result, err := client.Get(value.Key).Result(); err != nil {
		assert.Fail(t, "Failed to retrieve data from Redis with key")
	} else {
		convertedResult, _ := strconv.ParseInt(result, 10, 64)
		assert.Equal(t, convertedResult, value.Value, "Values from Redis should be equal")
	}
	t.Cleanup(cleanUpRedis)
}

func TestIdempotentIncrement_Fail_NonIntValue(t *testing.T) {
	client := newRedisClient()
	defer client.Close()

	key :=uuid.New().String()
	client.Set(key, "HUHULULU", 15 * time.Minute)
	_, err := IdempotentIncrement(key)
	if err != nil {
		assert.Error(t, err, "Expected error to be thrown.")
	} else {
		assert.Fail(t, "Expected an error to be thrown!")
	}
	t.Cleanup(cleanUpRedis)
}

func TestDelete_Success(t *testing.T) {
	client := newRedisClient()
	defer client.Close()
	key := uuid.New().String()
	client.Set(key, "1", 15 * time.Minute)
	result := Delete(key)
	assert.True(t, result, "Expected deletion to be successful")
	t.Cleanup(cleanUpRedis)
}

func TestCheck_Success(t *testing.T) {
	result, err1 := IdempotentIncrement()
	checkResult, err2 := Check(result.Key)
	if err1 != nil || err2 != nil {
		assert.Fail(t, "Error setting up test")
	}
	assert.Equal(t, checkResult.Value, result.Value, "Check Results must be equal to Original")
	t.Cleanup(cleanUpRedis)
}


func cleanUpRedis() {
	client := newRedisClient()
	defer client.Close()
	client.FlushAll()
}

