package tgbotapi

import (
	"context"
	"sync"
)

type UserDataStorage interface {
	GetUserData(ctx context.Context, userId int64) (data UserData, err error)
	SetUserData(ctx context.Context, userId int64, ud UserData) (err error)
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

func (uds *InMemoryUserDataStorage) SetUserData(ctx context.Context, userId int64, ud UserData) (err error) {
	uds.mu.Lock()
	uds.Storage[userId] = ud
	uds.mu.Unlock()
	return nil
}

func (uds *InMemoryUserDataStorage) NewUserData(ctx context.Context, userId int64) (data UserData, err error) {
	return NewInMemoryUserData(), nil
}
