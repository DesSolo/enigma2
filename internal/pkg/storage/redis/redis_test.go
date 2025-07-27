//nolint:dupl
package redis

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"enigma/internal/pkg/storage"
	"enigma/internal/pkg/storage/redis/mocks"
)

func TestIsReady_ExpectOK(t *testing.T) {
	t.Parallel()

	client := mocks.NewUniversalClient(t)
	client.On("Ping", mock.Anything).Return(redis.NewStatusCmd(context.Background()))

	s := NewStorage(client)
	ready, err := s.IsReady(context.Background())
	assert.NoError(t, err)
	assert.True(t, ready)
}

func TestIsReady_ExpectErr(t *testing.T) {
	t.Parallel()

	client := mocks.NewUniversalClient(t)
	cmd := redis.NewStatusCmd(context.Background())
	cmd.SetErr(errors.New("test"))
	client.On("Ping", mock.Anything).Return(cmd)

	s := NewStorage(client)
	ready, err := s.IsReady(context.Background())
	assert.Error(t, err)
	assert.False(t, ready)
}

func TestGet_ExpectOK(t *testing.T) {
	t.Parallel()

	client := mocks.NewUniversalClient(t)
	client.On("Get", mock.Anything, "test").Return(redis.NewStringResult("test", nil))

	s := NewStorage(client)
	got, err := s.Get(context.Background(), "test")
	require.NoError(t, err)
	assert.Equal(t, "test", got)
}

func TestGet_NotFound_ExpectErr(t *testing.T) {
	t.Parallel()

	client := mocks.NewUniversalClient(t)
	client.On("Get", mock.Anything, "test").Return(redis.NewStringResult("", redis.Nil))

	s := NewStorage(client)
	got, err := s.Get(context.Background(), "test")
	require.ErrorIs(t, err, storage.ErrNotFound)
	assert.Equal(t, "", got)
}

func TestGet_Failed_ExpectErr(t *testing.T) {
	t.Parallel()

	client := mocks.NewUniversalClient(t)
	client.On("Get", mock.Anything, "test").Return(redis.NewStringResult("", errors.New("test")))

	s := NewStorage(client)
	got, err := s.Get(context.Background(), "test")
	require.Error(t, err)
	assert.Equal(t, "", got)
}

func TestSave_ExpectOK(t *testing.T) {
	t.Parallel()

	client := mocks.NewUniversalClient(t)
	client.On("Set", mock.Anything, "test", "test message", time.Hour).
		Return(redis.NewStatusCmd(context.Background()))

	s := NewStorage(client)
	err := s.Save(context.Background(), "test", "test message", time.Hour)
	require.NoError(t, err)
}

func TestSave_ExpectErr(t *testing.T) {
	t.Parallel()

	client := mocks.NewUniversalClient(t)
	cmd := redis.NewStatusCmd(context.Background())
	cmd.SetErr(errors.New("test"))
	client.On("Set", mock.Anything, "test", "test message", time.Hour).
		Return(cmd)

	s := NewStorage(client)
	err := s.Save(context.Background(), "test", "test message", time.Hour)
	require.Error(t, err)
}

func TestDelete_ExpectOK(t *testing.T) {
	t.Parallel()

	client := mocks.NewUniversalClient(t)
	client.On("Del", mock.Anything, "test").Return(redis.NewIntCmd(context.Background()))

	s := NewStorage(client)
	err := s.Delete(context.Background(), "test")
	require.NoError(t, err)
}

func TestDelete_ExpectErr(t *testing.T) {
	t.Parallel()

	client := mocks.NewUniversalClient(t)
	client.On("Del", mock.Anything, "test").
		Return(redis.NewIntResult(0, errors.New("test")))

	s := NewStorage(client)
	err := s.Delete(context.Background(), "test")
	require.Error(t, err)
}

func TestIsUniq_ExpectTrue(t *testing.T) {
	t.Parallel()

	client := mocks.NewUniversalClient(t)
	client.On("Exists", mock.Anything, "test").Return(redis.NewIntResult(0, nil))

	s := NewStorage(client)
	got, err := s.IsUniq(context.Background(), "test")
	require.NoError(t, err)
	assert.True(t, got)
}

func TestIsUniq_ExpectFalse(t *testing.T) {
	t.Parallel()

	client := mocks.NewUniversalClient(t)
	client.On("Exists", mock.Anything, "test").Return(redis.NewIntResult(1, nil))

	s := NewStorage(client)
	got, err := s.IsUniq(context.Background(), "test")
	require.NoError(t, err)
	assert.False(t, got)
}

func TestIsUniq_ExpectErr(t *testing.T) {
	t.Parallel()

	client := mocks.NewUniversalClient(t)
	client.On("Exists", mock.Anything, "test").Return(redis.NewIntResult(0, errors.New("test")))

	s := NewStorage(client)
	got, err := s.IsUniq(context.Background(), "test")
	require.Error(t, err)
	assert.False(t, got)
}

func TestClose_ExpectOK(t *testing.T) {
	t.Parallel()

	client := mocks.NewUniversalClient(t)
	client.On("Close").Return(nil)

	s := NewStorage(client)
	err := s.Close()
	require.NoError(t, err)
}

func TestClose_ExpectErr(t *testing.T) {
	t.Parallel()

	client := mocks.NewUniversalClient(t)
	client.On("Close").Return(errors.New("test"))

	s := NewStorage(client)
	err := s.Close()
	require.Error(t, err)
}
