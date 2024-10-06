package mocks

import (
	"context"
	"github.com/cuctemeh/rstream-consumer/internal/storage"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
)

type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) SCard(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRedisClient) XRead(ctx context.Context, args *redis.XReadArgs) ([]redis.XStream, error) {
	ret := m.Called(ctx, args)
	return ret.Get(0).([]redis.XStream), ret.Error(1)
}

func (m *MockRedisClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRedisClient) LLen(ctx context.Context, key string) (int64, error) {
	testifyArgs := m.Called(ctx, key)
	return testifyArgs.Get(0).(int64), testifyArgs.Error(1)
}

func (m *MockRedisClient) Subscribe(ctx context.Context, channel string) storage.PubSub {
	testifyArgs := m.Called(ctx, channel)
	return testifyArgs.Get(0).(*MockRedisPubSub)
}

func (m *MockRedisClient) SIsMember(ctx context.Context, set string, member interface{}) (bool, error) {
	testifyArgs := m.Called(ctx, set, member)
	return testifyArgs.Bool(0), testifyArgs.Error(1)
}

func (m *MockRedisClient) XAdd(ctx context.Context, args *redis.XAddArgs) (string, error) {
	testifyArgs := m.Called(ctx, args)
	return testifyArgs.String(0), testifyArgs.Error(1)
}

func (m *MockRedisClient) XLen(ctx context.Context, key string) (int64, error) {
	testifyArgs := m.Called(ctx, key)
	return testifyArgs.Get(0).(int64), testifyArgs.Error(1)
}

func (m *MockRedisClient) SAdd(ctx context.Context, key string, members ...interface{}) (int64, error) {
	testifyArgs := m.Called(ctx, key, members)
	return testifyArgs.Get(0).(int64), testifyArgs.Error(1)
}

func (m *MockRedisClient) SRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	testifyArgs := m.Called(ctx, key, members)
	return testifyArgs.Get(0).(int64), testifyArgs.Error(1)
}
