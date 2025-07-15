package secrets

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"

	mocks_hasher "enigma/internal/pkg/hasher/mocks"
	mocks_storage "enigma/internal/pkg/storage/mocks"
)

type mk struct {
	storage  *mocks_storage.SecretStorage
	hasher   *mocks_hasher.Hasher
	provider *Provider
}

func newMK(t *testing.T) *mk {
	storage := mocks_storage.NewSecretStorage(t)
	hasher := mocks_hasher.NewHasher(t)
	return &mk{
		storage:  storage,
		hasher:   hasher,
		provider: New(storage, hasher),
	}
}

func TestSaveSecret_ExpectOk(t *testing.T) {
	t.Parallel()

	m := newMK(t)

	m.storage.EXPECT().
		IsUniq(
			mock.AnythingOfType("context.backgroundCtx"),
			mock.AnythingOfType("string"),
		).
		Return(true, nil).
		Once()

	m.hasher.EXPECT().
		Encrypt("test_message").
		Return("encrypter_test_message", nil).
		Times(1)

	m.storage.EXPECT().
		Save(
			mock.AnythingOfType("context.backgroundCtx"),
			mock.AnythingOfType("string"),
			"encrypter_test_message", 24*time.Hour,
		).
		Return(nil).
		Once()

	got, err := m.provider.SaveSecret(context.Background(), "test_message", 1)
	require.NoError(t, err)
	require.NotEmpty(t, got)
}
