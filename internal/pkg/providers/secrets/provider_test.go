package secrets

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"

	"enigma/internal/pkg/storage/mocks"
)

type mk struct {
	storage  *mocks.SecretStorage
	provider *Provider
}

func newMK(t *testing.T) *mk {
	storage := mocks.NewSecretStorage(t)
	return &mk{
		storage:  storage,
		provider: New(storage),
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

	m.storage.EXPECT().
		Save(
			mock.AnythingOfType("context.backgroundCtx"),
			mock.AnythingOfType("string"),
			"test_message", 1,
		).
		Return(nil).
		Once()

	got, err := m.provider.SaveSecret(context.Background(), "test_message", 1)
	require.NoError(t, err)
	require.NotEmpty(t, got)
}
