package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
	"login/modle"

	"github.com/teampui/pac"
	"os"
	"strings"
	"time"
)

type RedisRepoInterface interface {
	pac.Service

	// member
	GetSortedSetAllMembers(key string) []string
	GetLowestScoreMember(key string) string
	GetSortedSetAllMembersWithScores(key string) []redis.Z
	AddSortedSetMember(key string, member string, score float64) error
	RemoveSortedSetMember(key string, member string) error
	CountSortedSetMembers(key string) int64

	// expired time
	GetExpiredTime(key string) time.Duration
	SetExpiredTime(key string, duration time.Duration) error

	HSetItem(key string, item *model.Item) error

	SetString(key string, value any) error
	GetString(key string) (string, error)
}

type RedisRepo struct {
	rdb *redis.Client
}

func (repo *RedisRepo) Register(app *pac.App) {
	app.Repositories.Add("redis", repo)

	rdb := redis.NewClient(&redis.Options{
		Addr:     strings.TrimPrefix(os.Getenv("REDIS_DSN"), "redis://"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	repo.rdb = rdb
}

func (repo *RedisRepo) CountSortedSetMembers(key string) int64 {
	ctx := context.Background()

	return repo.rdb.ZCard(ctx, key).Val()
}

// GetExpiredTime get expired time
func (repo *RedisRepo) GetExpiredTime(key string) time.Duration {
	ctx := context.Background()
	return repo.rdb.TTL(ctx, key).Val()
}

// SetExpiredTime set expired time
func (repo *RedisRepo) SetExpiredTime(key string, duration time.Duration) error {
	ctx := context.Background()
	return repo.rdb.Expire(ctx, key, duration).Err()
}

// AddSortedSetMember add member to sorted set
func (repo *RedisRepo) AddSortedSetMember(key string, member string, score float64) error {
	ctx := context.Background()
	return repo.rdb.ZAdd(ctx, key, &redis.Z{
		Score:  score,
		Member: member,
	}).Err()
}

// GetSortedSetAllMembers get sorted set members
func (repo *RedisRepo) GetSortedSetAllMembers(key string) []string {
	ctx := context.Background()
	return repo.rdb.ZRange(ctx, key, 0, -1).Val()
}

// GetSortedSetAllMembersWithScores get sorted set members with scores
func (repo *RedisRepo) GetSortedSetAllMembersWithScores(key string) []redis.Z {
	ctx := context.Background()
	return repo.rdb.ZRangeWithScores(ctx, key, 0, -1).Val()
}

// RemoveSortedSetMember Remove member from sorted set
func (repo *RedisRepo) RemoveSortedSetMember(key string, member string) error {
	ctx := context.Background()
	return repo.rdb.ZRem(ctx, key, member).Err()
}

// GetLowestScoreMember get the lowest score member
func (repo *RedisRepo) GetLowestScoreMember(key string) string {
	ctx := context.Background()
	return repo.rdb.ZRange(ctx, key, 0, 0).Val()[0]
}

func (repo *RedisRepo) HSetItem(key string, item *model.Item) error {
	ctx := context.Background()
	if _, err := repo.rdb.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, key, "id", item.Id)
		rdb.HSet(ctx, key, "str2", "world")
		rdb.HSet(ctx, key, "int", 123)
		rdb.HSet(ctx, key, "bool", 1)
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (repo *RedisRepo) SetString(key string, value any) error {
	ctx := context.Background()
	return repo.rdb.Set(ctx, key, value, 86400*time.Second).Err()
}

func (repo *RedisRepo) GetString(key string) (string, error) {
	ctx := context.Background()
	return repo.rdb.Get(ctx, key).Result()
}
