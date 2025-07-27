package memory

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"

	"github.com/stretchr/testify/assert"
)

func TestIsReady(t *testing.T) {
	t.Parallel()

	s := NewStorage()
	ready, err := s.IsReady(context.Background())
	assert.NoError(t, err)
	assert.True(t, ready)
}

func TestGet_ExpectOK(t *testing.T) {
	t.Parallel()

	s := NewStorage()
	s.secrets["test"] = &secret{
		text:   "test",
		expire: time.Now().UTC().Add(time.Hour),
	}

	got, err := s.Get(context.Background(), "test")
	require.NoError(t, err)
	assert.Equal(t, "test", got)
}

func TestGet_notFound_ExpectErr(t *testing.T) {
	t.Parallel()

	s := NewStorage()

	got, err := s.Get(context.Background(), "test")
	require.EqualError(t, err, "not found")
	assert.Equal(t, "", got)
}

func TestGet_expired_ExpectErr(t *testing.T) {
	t.Parallel()

	s := NewStorage()
	s.secrets["test"] = &secret{
		text:   "test",
		expire: time.Now().UTC().Add(-time.Hour),
	}

	got, err := s.Get(context.Background(), "test")
	require.EqualError(t, err, "not found")
	assert.Equal(t, "", got)
}

func TestSave_ExpectOK(t *testing.T) {
	t.Parallel()

	s := NewStorage()

	err := s.Save(context.Background(), "test", "test message", time.Hour)
	require.NoError(t, err)
}

func TestDelete_ExpectOK(t *testing.T) {
	t.Parallel()

	s := NewStorage()
	err := s.Delete(context.Background(), "test")
	require.NoError(t, err)
}

func TestIsUniq_ExpectOK(t *testing.T) {
	t.Parallel()

	s := NewStorage()
	s.secrets["test"] = &secret{
		text:   "test",
		expire: time.Now().UTC().Add(time.Hour),
	}

	got, err := s.IsUniq(context.Background(), "test")
	require.NoError(t, err)
	assert.False(t, got)
}

func TestClose_ExpectOK(t *testing.T) {
	t.Parallel()

	s := NewStorage()

	err := s.Close()
	require.NoError(t, err)
}
