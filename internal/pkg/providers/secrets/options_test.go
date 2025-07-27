package secrets

import (
	"testing"

	"github.com/stretchr/testify/require"

	mocks_hasher "enigma/internal/pkg/hasher/mocks"
	mocks_storage "enigma/internal/pkg/storage/mocks"
)

func TestOptions(t *testing.T) {
	t.Parallel()

	storage := mocks_storage.NewSecretStorage(t)
	hasher := mocks_hasher.NewHasher(t)

	const (
		customTokenLength = 64
		customRetries     = 5
	)

	provider := New(
		storage,
		hasher,
		WithTokenLength(customTokenLength),
		WithTokenSaveRetries(customRetries),
	)

	require.Equal(t, customTokenLength, provider.tokenLength)
	require.Equal(t, customRetries, provider.tokenSaveRetries)
}
