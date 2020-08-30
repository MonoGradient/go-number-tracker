package api

import (
	"Go-Tracker/model"
	"Go-Tracker/service"
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestIdempotentIncrementStorageApi_Success(t *testing.T) {
	r, _ := http.NewRequest("POST", "/api/v1/storage", nil)
	w := httptest.NewRecorder()
	IdempotentIncrementStorageApi(w, r)
	assert.Equal(t, w.Code, http.StatusOK, "Expected successful status code")
	var responseData model.StorageResponse
	json.Unmarshal(w.Body.Bytes(), &responseData)
	assert.NotEmpty(t, responseData.Value)
	assert.Equal(t, getDataFromRedis(responseData.Key), responseData.Value)
	t.Cleanup(cleanUpRedis)
}

func TestIncrementStorageWithKeyApi_CreateSuccess(t *testing.T) {
	key := uuid.New().String()
	r, _ := http.NewRequest("POST", "/api/v1/storage", nil)
	r = mux.SetURLVars(r, map[string]string {
		"key": key,
	})
	w := httptest.NewRecorder()

	IncrementStorageWithKeyApi(w, r)
	var responseData model.StorageResponse

	assert.Equal(t, w.Code, http.StatusOK)
	dbValue := getDataFromRedis(key)

	json.Unmarshal(w.Body.Bytes(), &responseData)

	assert.Equal(t, dbValue, responseData.Value, "Values should match")
	t.Cleanup(cleanUpRedis)
}

func TestIncrementStorageWithKeyApi_Fail_InvalidKey(t *testing.T) {
	key := "fakekey"
	r, _ := http.NewRequest("POST", "/api/v1/storage", nil)
	r = mux.SetURLVars(r, map[string]string {
		"key": key,
	})
	w := httptest.NewRecorder()

	IncrementStorageWithKeyApi(w, r)
	var responseData model.StorageResponse

	assert.Equal(t, w.Code, http.StatusOK)

	json.Unmarshal(w.Body.Bytes(), &responseData)

	assert.Equal(t, getDataFromRedis(responseData.Key), responseData.Value, "Values should match")
	t.Cleanup(cleanUpRedis)
}

func TestIncrementStorageWithKeyApi_Success_InvalidKey(t *testing.T) {
	key := "fakekey"
	r, _ := http.NewRequest("POST", "/api/v1/storage", nil)
	r = mux.SetURLVars(r, map[string]string {
		"key": key,
	})
	w := httptest.NewRecorder()

	IncrementStorageWithKeyApi(w, r)
	var responseData model.StorageResponse

	assert.Equal(t, w.Code, http.StatusOK)

	json.Unmarshal(w.Body.Bytes(), &responseData)

	assert.Equal(t, getDataFromRedis(responseData.Key), responseData.Value, "Values should match")
	t.Cleanup(cleanUpRedis)
}

func TestDecrementStorageApi_Success(t *testing.T) {

	data, _ := service.IdempotentIncrement()
	r, _ := http.NewRequest("PUT", "/api/v1/decrement", nil)
	r = mux.SetURLVars(r, map[string]string{
		"key": data.Key,
	})
	w := httptest.NewRecorder()

	DecrementStorageApi(w, r)
	var responseData model.StorageResponse
	json.Unmarshal(w.Body.Bytes(), &responseData)

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, data.Value-1, responseData.Value, "Expected values to be the equal")

}

func TestDecrementStorageApi_Fail_InvalidKey(t *testing.T) {
	key := "invalidkey"
	r, _ := http.NewRequest("PUT", "/api/v1/decrement", nil)
	r = mux.SetURLVars(r, map[string]string{
		"key": key,
	})
	w := httptest.NewRecorder()
	DecrementStorageApi(w, r)

	assert.Equal(t, w.Code, http.StatusInternalServerError)
}

func TestCheckStorageApi_Success(t *testing.T) {
	data, _ := service.IdempotentIncrement()
	r, _ := http.NewRequest("GET", "/api/v1/storage", nil)
	r = mux.SetURLVars(r, map[string]string{
		"key": data.Key,
	})
	w := httptest.NewRecorder()
	CheckStorageApi(w, r)
	var responseData model.StorageResponse
	json.Unmarshal(w.Body.Bytes(), &responseData)

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, responseData.Value, data.Value)
	assert.Equal(t, responseData.Key, data.Key)

}

func TestCheckStorageApi_Fail_InvalidKey(t *testing.T) {
	key := "invalidkey"
	r, _ := http.NewRequest("GET", "/api/v1/storage", nil)
	r = mux.SetURLVars(r, map[string]string{
		"key": key,
	})
	w := httptest.NewRecorder()
	CheckStorageApi(w, r)
	var responseData model.StorageResponse
	json.Unmarshal(w.Body.Bytes(), &responseData)

	assert.Equal(t, w.Code, http.StatusInternalServerError)
}

func TestDeleteStorageApi_Success(t *testing.T) {
	data, _ := service.IdempotentIncrement()
	r, _ := http.NewRequest("DELETE", "/api/v1/storage", nil)
	r = mux.SetURLVars(r, map[string]string{
		"key": data.Key,
	})
	w := httptest.NewRecorder()
	DeleteStorageApi(w, r)

	assert.Equal(t, w.Code, http.StatusOK)
}

func TestDeleteStorageApi_Fail_InvalidKey(t *testing.T) {
	key := "invalidkey"
	r, _ := http.NewRequest("DELETE", "/api/v1/storage", nil)
	r = mux.SetURLVars(r, map[string]string{
		"key": key,
	})
	w := httptest.NewRecorder()
	DeleteStorageApi(w, r)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestIncrementStorageWithKeyApi_Success_Increment(t *testing.T) {
	key := uuid.New().String()
	value := rand.Int63()
	client := newRedisClient()
	defer client.Close()
	client.Set(key, value, 15 * time.Minute)
	r, _ := http.NewRequest("POST", "/api/v1/storage", nil)
	r = mux.SetURLVars(r, map[string]string {
		"key": key,
	})
	w := httptest.NewRecorder()

	IncrementStorageWithKeyApi(w, r)
	var responseData model.StorageResponse

	assert.Equal(t, w.Code, http.StatusOK)

	json.Unmarshal(w.Body.Bytes(), &responseData)

	assert.Equal(t, getDataFromRedis(responseData.Key), responseData.Value, "Values should match")
	assert.Equal(t, value+1, responseData.Value, "Values should match")
	t.Cleanup(cleanUpRedis)
}

func getDataFromRedis(key string) int64 {
	client := newRedisClient()
	defer client.Close()
	result, _ := client.Get(key).Result()
	conv, _ := strconv.ParseInt(result, 10, 64)
	return conv
}

func init() {
	os.Setenv("REDIS_HOST", "localhost:6379")
	os.Setenv("REDIS_PASSWORD", "")
}

func cleanUpRedis() {
	client := newRedisClient()
	defer client.Close()
	client.FlushAll()
}

func newRedisClient() *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB: 0,
	})

	return rdb
}
