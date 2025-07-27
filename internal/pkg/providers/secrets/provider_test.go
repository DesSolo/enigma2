package secrets

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"

	mocks_hasher "enigma/internal/pkg/hasher/mocks"
	"enigma/internal/pkg/storage"
	mocks_storage "enigma/internal/pkg/storage/mocks"
)

type mk struct {
	storage  *mocks_storage.SecretStorage
	hasher   *mocks_hasher.Hasher
	provider *Provider
}

func newMK(t *testing.T) *mk {
	secretsStorage := mocks_storage.NewSecretStorage(t)
	hasher := mocks_hasher.NewHasher(t)
	return &mk{
		storage:  secretsStorage,
		hasher:   hasher,
		provider: New(secretsStorage, hasher),
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
		Once()

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

func TestSaveSecret_generateUniqToken_ExpectErr(t *testing.T) {
	t.Parallel()

	m := newMK(t)

	m.storage.EXPECT().
		IsUniq(
			mock.AnythingOfType("context.backgroundCtx"),
			mock.AnythingOfType("string"),
		).
		Return(false, nil).
		Times(m.provider.tokenSaveRetries)

	got, err := m.provider.SaveSecret(context.Background(), "test_message", 1)
	require.EqualError(t, err, "s.generateUniqToken: maximum retries save")
	require.Empty(t, got)
}

func TestSaveSecret_generateUniqToken_IsUniq_ExpectErr(t *testing.T) {
	t.Parallel()

	m := newMK(t)

	m.storage.EXPECT().
		IsUniq(
			mock.AnythingOfType("context.backgroundCtx"),
			mock.AnythingOfType("string"),
		).
		Return(false, errors.New("test_err")).
		Once()

	got, err := m.provider.SaveSecret(context.Background(), "test_message", 1)
	require.EqualError(t, err, "s.generateUniqToken: storage.IsUniq: test_err")
	require.Empty(t, got)
}

func TestSaveSecret_Encrypt_ExpectErr(t *testing.T) {
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
		Return("", errors.New("test_err")).
		Once()

	got, err := m.provider.SaveSecret(context.Background(), "test_message", 1)
	require.EqualError(t, err, "s.Encrypt: test_err")
	require.Empty(t, got)
}

func TestSaveSecret_Save_ExpectErr(t *testing.T) {
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
		Once()

	m.storage.EXPECT().
		Save(
			mock.AnythingOfType("context.backgroundCtx"),
			mock.AnythingOfType("string"),
			"encrypter_test_message", 24*time.Hour,
		).
		Return(errors.New("test_err")).
		Once()

	got, err := m.provider.SaveSecret(context.Background(), "test_message", 1)
	require.EqualError(t, err, "storage.Save: test_err")
	require.Empty(t, got)
}

func TestGetSecret_ExpectOk(t *testing.T) {
	t.Parallel()

	m := newMK(t)

	m.storage.EXPECT().
		Get(
			mock.AnythingOfType("context.backgroundCtx"),
			"test_token",
		).
		Return("encrypted_message", nil).
		Once()

	m.hasher.EXPECT().
		Decrypt("encrypted_message").
		Return("decrypted_message", nil).
		Once()

	m.storage.EXPECT().
		Delete(
			mock.AnythingOfType("context.backgroundCtx"),
			"test_token",
		).
		Return(nil).
		Once()

	got, err := m.provider.GetSecret(context.Background(), "test_token")
	require.NoError(t, err)
	require.Equal(t, "decrypted_message", got)
}

func TestGetSecret_Get_ExpectErr(t *testing.T) {
	t.Parallel()

	m := newMK(t)

	m.storage.EXPECT().
		Get(
			mock.AnythingOfType("context.backgroundCtx"),
			"test_token",
		).
		Return("", errors.New("test_err")).
		Once()

	got, err := m.provider.GetSecret(context.Background(), "test_token")
	require.EqualError(t, err, "storage.Get: test_err")
	require.Empty(t, got)
}

func TestGetSecret_Decrypt_ExpectErr(t *testing.T) {
	t.Parallel()

	m := newMK(t)

	m.storage.EXPECT().
		Get(
			mock.AnythingOfType("context.backgroundCtx"),
			"test_token",
		).
		Return("encrypted_message", nil).
		Once()

	m.hasher.EXPECT().
		Decrypt("encrypted_message").
		Return("", errors.New("test_err")).
		Once()

	got, err := m.provider.GetSecret(context.Background(), "test_token")
	require.EqualError(t, err, "p.hasher.Decrypt: test_err")
	require.Empty(t, got)
}

func TestGetSecret_Delete_ExpectErr(t *testing.T) {
	t.Parallel()

	m := newMK(t)

	m.storage.EXPECT().
		Get(
			mock.AnythingOfType("context.backgroundCtx"),
			"test_token",
		).
		Return("encrypted_message", nil).
		Once()

	m.hasher.EXPECT().
		Decrypt("encrypted_message").
		Return("decrypted_message", nil).
		Once()

	m.storage.EXPECT().
		Delete(
			mock.AnythingOfType("context.backgroundCtx"),
			"test_token",
		).
		Return(errors.New("test_err")).
		Once()

	got, err := m.provider.GetSecret(context.Background(), "test_token")
	require.EqualError(t, err, "storage.Delete: test_err")
	require.Empty(t, got)
}

func TestCheckExistsSecret_ExpectOk(t *testing.T) {
	t.Parallel()

	m := newMK(t)

	m.storage.EXPECT().
		IsUniq(
			mock.AnythingOfType("context.backgroundCtx"),
			"test_token",
		).
		Return(false, nil).
		Once()

	err := m.provider.CheckExistsSecret(context.Background(), "test_token")
	require.NoError(t, err)
}

func TestCheckExistsSecret_ExpectErr(t *testing.T) {
	t.Parallel()

	m := newMK(t)

	m.storage.EXPECT().
		IsUniq(
			mock.AnythingOfType("context.backgroundCtx"),
			"test_token",
		).
		Return(true, nil).
		Once()

	err := m.provider.CheckExistsSecret(context.Background(), "test_token")
	require.ErrorIs(t, err, storage.ErrNotFound)
}

func TestCheckExistsSecret_IsUniq_ExpectErr(t *testing.T) {
	t.Parallel()

	m := newMK(t)

	m.storage.EXPECT().
		IsUniq(
			mock.AnythingOfType("context.backgroundCtx"),
			"test_token",
		).
		Return(false, errors.New("test_err")).
		Once()

	err := m.provider.CheckExistsSecret(context.Background(), "test_token")
	require.EqualError(t, err, "storage.IsUniq: test_err")
}
