// mocks/mock_redis_pubsub.go
package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockRedisPubSub struct {
	mock.Mock
}

func (m *MockRedisPubSub) ReceiveMessage(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockRedisPubSub) Close() error {
	args := m.Called()
	return args.Error(0)
}
