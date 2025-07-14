package memory_test

import (
	"testing"

	"enigma/internal/storage/memory"

	"github.com/stretchr/testify/assert"
)

func TestGetInfo(t *testing.T) {
	s := memory.NewStorage()
	assert.Equal(t, s.GetInfo(), "Memory")
}

func TestIsReady(t *testing.T) {
	s := memory.NewStorage()
	ready, err := s.IsReady()
	assert.NoError(t, err)
	assert.True(t, ready)
}

func TestSave(t *testing.T) {
	s := memory.NewStorage()
	assert.NoError(t, s.Save("example", "msg", 1))
}

var cases = []struct {
	Key     string
	Message string
}{
	{Key: "1", Message: "1"},
	{Key: "2", Message: "2"},
	{Key: "3", Message: "3"},
}

func TestGet(t *testing.T) {
	s := memory.NewStorage()
	for _, tc := range cases {
		s.Save(tc.Key, tc.Message, 1) // nolint:errcheck
		msg, err := s.Get(tc.Key)
		assert.NoError(t, err)
		assert.Equal(t, msg, tc.Message)
	}
}

func TestGetExpired(t *testing.T) {
	s := memory.NewStorage()
	for _, tc := range cases {
		s.Save(tc.Key, tc.Message, -1) // nolint:errcheck
		msg, err := s.Get(tc.Key)
		assert.NoError(t, err)
		assert.NotEqual(t, msg, tc.Message)
	}
}

func TestDelete(t *testing.T) {
	s := memory.NewStorage()
	for _, tc := range cases {
		s.Save(tc.Key, tc.Message, 1) // nolint:errcheck
		assert.NoError(t, s.Delete(tc.Key))

		msg, _ := s.Get(tc.Key)
		assert.NotEqual(t, msg, tc.Message)
	}
}

// TODO: fix tests

func TestIsUniq(t *testing.T) {
	s := memory.NewStorage()
	for _, tc := range cases {
		s.Save(tc.Key, tc.Message, 1) // nolint:errcheck
		uniq, err := s.IsUniq(tc.Key)
		assert.NoError(t, err)
		assert.False(t, uniq)

		s.Delete(tc.Key) // nolint:errcheck
		uniq1, _ := s.IsUniq(tc.Key)
		assert.True(t, uniq1)
	}
}
