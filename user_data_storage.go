package tgbotapi

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

type UserDataStorage interface {
	GetUserData(ctx context.Context, userId int64) (data UserData, err error)
	SetUserData(ctx context.Context, userId int64, data UserData) (err error)
	NewUserData(ctx context.Context, userId int64) (data UserData, err error)
}

type UserData interface {
	SetValue(ctx context.Context, key, value string) (err error)
	GetValue(ctx context.Context, key string) (value string, err error)
}

type InMemoryUserDataStorage struct {
	Storage map[int64]UserData
	mu      sync.RWMutex
}
type InMemoryUserData struct {
	Data map[string]string
	mu   sync.RWMutex
}

func (ud *InMemoryUserData) SetValue(ctx context.Context, key, value string) (err error) {
	ud.mu.Lock()
	ud.Data[key] = value
	ud.mu.Unlock()
	return nil
}

func (ud *InMemoryUserData) GetValue(ctx context.Context, key string) (value string, err error) {
	ud.mu.RLock()
	value = ud.Data[key]
	ud.mu.RUnlock()
	return value, nil
}

func NewInMemoryUserDataStorage() *InMemoryUserDataStorage {
	return &InMemoryUserDataStorage{
		Storage: make(map[int64]UserData),
		mu:      sync.RWMutex{},
	}
}
func NewInMemoryUserData() *InMemoryUserData {
	return &InMemoryUserData{
		Data: make(map[string]string),
		mu:   sync.RWMutex{},
	}
}

func (uds *InMemoryUserDataStorage) GetUserData(ctx context.Context, userId int64) (data UserData, err error) {
	uds.mu.RLock()
	r := uds.Storage[userId]
	uds.mu.RUnlock()
	return r, nil
}

func (uds *InMemoryUserDataStorage) SetUserData(ctx context.Context, userId int64, data UserData) (err error) {
	uds.mu.Lock()
	uds.Storage[userId] = data
	uds.mu.Unlock()
	return nil
}

func (uds *InMemoryUserDataStorage) NewUserData(ctx context.Context, userId int64) (data UserData, err error) {
	return NewInMemoryUserData(), nil
}

type RedisUserDataStorage struct {
	Rdb         *redis.Client
	RedisPrefix string
}

func NewRedisUserDataStorage(redisUrl string) (*RedisUserDataStorage, error) {
	opts, err := redis.ParseURL(redisUrl)
	if err != nil {
		return nil, fmt.Errorf("redis not connected: %v", redisUrl)
	}
	return &RedisUserDataStorage{Rdb: redis.NewClient(opts), RedisPrefix: "RedisUserDataStorage"}, nil
}

func (ruds *RedisUserDataStorage) Ping(ctx context.Context) error {
	return ruds.Rdb.Ping(ctx).Err()
}

func (uds *RedisUserDataStorage) GetUserData(ctx context.Context, userId int64) (data UserData, err error) {
	// key := fmt.Sprintf("%v-%v", uds.RedisPrefix, userId)
	// userDataString, err := uds.Rdb.Get(ctx, key).Result()
	// if err != nil {
	// 	return nil, fmt.Errorf("redis problem Get Key: %v", key)
	// }
	// err = json.Unmarshal([]byte(userDataString), &data)
	// if err != nil {
	// 	return nil, fmt.Errorf("can't Unmarshal: %v", userDataString)
	// }
	// return data, nil
	return uds.NewUserData(ctx, userId)
}

func (uds *RedisUserDataStorage) SetUserData(ctx context.Context, userId int64, data UserData) (err error) {
	// key := fmt.Sprintf("%v-%v", uds.RedisPrefix, userId)
	// jsonBytes, err := json.Marshal(data)
	// if err != nil {
	// 	return fmt.Errorf("can't Marshal: %v", data)
	// }
	// err = uds.Rdb.Set(ctx, key, string(jsonBytes), StorageConvTimeout).Err()
	// if err == nil {
	// 	return nil
	// }
	// return fmt.Errorf("redis problem Set Key: %v", key)
	return nil
}

func (uds *RedisUserDataStorage) NewUserData(ctx context.Context, userId int64) (data UserData, err error) {
	return NewRedisUserData(uds.Rdb, fmt.Sprintf("%v-%v", uds.RedisPrefix, userId)), nil
}

type RedisUserData struct {
	Rdb *redis.Client
	Key string
}

func NewRedisUserData(rdb *redis.Client, key string) *RedisUserData {
	return &RedisUserData{
		Rdb: rdb,
		Key: key,
	}
}

func (rud *RedisUserData) GetValue(ctx context.Context, key string) (value string, err error) {
	userDataString, err := rud.Rdb.Get(ctx, rud.Key).Result()
	switch {
	case err == redis.Nil:
		fmt.Println("key does not exist")
	case err != nil:
		return "", fmt.Errorf("redis problem Get Key: %v", rud.Key)
	}
	data := make(map[string]string)
	err = json.Unmarshal([]byte(userDataString), &data)
	if err != nil {
		return "", fmt.Errorf("can't Unmarshal: %v", userDataString)
	}
	value = data[key]
	return value, nil
}

func (rud *RedisUserData) SetValue(ctx context.Context, key, value string) (err error) {
	userDataString, err := rud.Rdb.Get(ctx, rud.Key).Result()
	switch {
	case err == redis.Nil:
		fmt.Println("key does not exist")
	case err != nil:
		return fmt.Errorf("redis problem Set Key: %v", rud.Key)
	}
	data := make(map[string]string)
	err = json.Unmarshal([]byte(userDataString), &data)
	if userDataString != "" && err != nil {
		return fmt.Errorf("can't Unmarshal: %v", userDataString)
	}
	data[key] = value
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("can't Marshal: %v", data)
	}
	err = rud.Rdb.Set(ctx, rud.Key, string(jsonBytes), StorageConvTimeout).Err()
	if err == nil {
		return nil
	}
	return fmt.Errorf("redis problem Set Key: %v", key)
}
