package tgbotapi

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const StorageConvTimeout time.Duration = 600000000000 // 10 min

type StateStorage interface {
	GetState(ctx context.Context, update Update) (state string, err error)
	SetState(ctx context.Context, update Update, state string) (result bool, err error)
	ClearState(ctx context.Context, update Update) (err error)
}

type RedisStateStorage struct {
	Rdb *redis.Client
}

func NewRedisStateStorage(redisUrl string) (*RedisStateStorage, error) {
	opts, err := redis.ParseURL(redisUrl)
	if err != nil {
		return nil, fmt.Errorf("redis not connected: %v", redisUrl)
	}
	return &RedisStateStorage{Rdb: redis.NewClient(opts)}, nil
}
func (redisStateStorage *RedisStateStorage) Ping(ctx context.Context) error {
	return redisStateStorage.Rdb.Ping(ctx).Err()
}

func BuildKey(update Update) (key string) {
	if update.InlineQuery != nil {
		key += update.InlineQuery.ID
	} else if update.CallbackQuery != nil {
		key += strconv.Itoa(update.CallbackQuery.Message.MessageID)
	}
	if update.EffectiveUser() != nil {
		key += strconv.FormatInt(update.EffectiveUser().ID, 10)
	}
	if update.PreCheckoutQuery != nil {
		key += strconv.FormatInt(update.PreCheckoutQuery.From.ID, 10)
	} else if update.EffectiveChat() != nil {
		key += strconv.FormatInt(update.EffectiveChat().ID, 10)
	}
	return key
}
func (redisStateStorage *RedisStateStorage) GetState(ctx context.Context, update Update) (state string, err error) {
	key := BuildKey(update)
	state, err = redisStateStorage.Rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", fmt.Errorf("redis problem Get Key: %v", key)
	}
	return state, nil
}

func (redisStateStorage *RedisStateStorage) SetState(ctx context.Context, update Update, state string) (result bool, err error) {
	key := BuildKey(update)
	err = redisStateStorage.Rdb.Set(ctx, key, state, StorageConvTimeout).Err()
	if err == nil {
		return true, nil
	}
	return false, fmt.Errorf("redis problem Set Key: %v", key)
}
func (redisStateStorage *RedisStateStorage) ClearState(ctx context.Context, update Update) (err error) {
	_, err = redisStateStorage.SetState(ctx, update, "")
	return err
}

type NilStateStorage struct {
}

func NewNilStateStorage() (*NilStateStorage, error) {
	return &NilStateStorage{}, nil
}

func (nilStateStorage *NilStateStorage) GetState(ctx context.Context, update Update) (state string, err error) {
	return "", nil
}

func (nilStateStorage *NilStateStorage) SetState(ctx context.Context, update Update, state string) (result bool, err error) {
	return true, nil
}
func (nilStateStorage *NilStateStorage) ClearState(ctx context.Context, update Update) (err error) {
	return nil
}
