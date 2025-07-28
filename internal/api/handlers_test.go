package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"enigma/internal/pkg/adapters/template"
	"enigma/internal/pkg/storage"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newTestRequestWithForm(t *testing.T, data map[string]string) *http.Request {
	t.Helper()

	form := url.Values{}

	for k, v := range data {
		form.Add(k, v)
	}

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req
}

func newTestRequestWithChiContext(t *testing.T, method, path string, pathParams map[string]string) *http.Request {
	t.Helper()

	req := httptest.NewRequest(method, path, nil)
	chiCtx := chi.NewRouteContext()
	for k, v := range pathParams {
		chiCtx.URLParams.Add(k, v)
	}
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx)
	return req.WithContext(ctx)
}

func Test_healthHandler_ExpectOl(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()

	healthHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, `{"healthy": true}`, rec.Body.String())
}

func Test_indexHandler_ExpectOK(t *testing.T) {
	t.Parallel()

	m := newMK(t)

	m.template.EXPECT().
		RenderFile("index.html", mock.AnythingOfType("*httptest.ResponseRecorder"), template.Data(nil)).
		Return(nil).
		Once()

	rec := httptest.NewRecorder()

	m.server.indexHandler(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	require.Equal(t, http.StatusOK, rec.Code)
}

func Test_createSecretHandler_ExpectOK(t *testing.T) {
	t.Parallel()

	m := newMK(t)

	m.secretsProvider.EXPECT().
		SaveSecret(mock.AnythingOfType("context.backgroundCtx"), "secret", 1).
		Return("test_token", nil).
		Once()

	rec := httptest.NewRecorder()

	m.server.createSecretHandler(rec, newTestRequestWithForm(t, map[string]string{
		"msg": "secret",
		"due": "1",
	}))
	require.Equal(t, http.StatusOK, rec.Code)
}

func Test_createSecretHandler_notValid_ExpectErr(t *testing.T) {
	t.Parallel()

	m := newMK(t)

	rec := httptest.NewRecorder()

	m.server.createSecretHandler(rec, newTestRequestWithForm(t, map[string]string{
		"msg": "secret",
		"due": "",
	}))
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func Test_createSecretHandler_notValidDues_ExpectErr(t *testing.T) {
	t.Parallel()

	m := newMK(t)

	rec := httptest.NewRecorder()

	m.server.createSecretHandler(rec, newTestRequestWithForm(t, map[string]string{
		"msg": "secret",
		"due": "5",
	}))
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func Test_createSecretHandler_providerError_ExpectErr(t *testing.T) {
	t.Parallel()

	m := newMK(t)

	m.secretsProvider.EXPECT().
		SaveSecret(mock.AnythingOfType("context.backgroundCtx"), "secret", 1).
		Return("", errors.New("test error")).
		Once()

	rec := httptest.NewRecorder()

	m.server.createSecretHandler(rec, newTestRequestWithForm(t, map[string]string{
		"msg": "secret",
		"due": "1",
	}))
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func Test_viewSecretHandler_noToken_ExpectNotFound(t *testing.T) {
	t.Parallel()

	m := newMK(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/get/", nil)

	m.server.viewSecretHandler(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func Test_viewSecretHandler_viewConformPage_Exists_ExpectOK(t *testing.T) {
	t.Parallel()

	m := newMK(t)
	token := "test_token" // nolint:goconst

	m.secretsProvider.EXPECT().
		CheckExistsSecret(mock.Anything, token).
		Return(nil).
		Once()

	m.template.EXPECT().
		RenderFile("conform.html", mock.Anything, template.Data(nil)).
		Return(nil).
		Once()

	rec := httptest.NewRecorder()
	req := newTestRequestWithChiContext(t, http.MethodGet, "/get/"+token, map[string]string{"token": token})

	m.server.viewSecretHandler(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func Test_viewSecretHandler_viewConformPage_NotFound_ExpectNotFound(t *testing.T) {
	t.Parallel()

	m := newMK(t)
	token := "test_token"

	m.secretsProvider.EXPECT().
		CheckExistsSecret(mock.Anything, token).
		Return(storage.ErrNotFound).
		Once()

	m.template.EXPECT().
		RenderFile("expired.html", mock.Anything, template.Data(nil)).
		Return(nil).
		Once()

	rec := httptest.NewRecorder()
	req := newTestRequestWithChiContext(t, http.MethodGet, "/get/"+token, map[string]string{"token": token})

	m.server.viewSecretHandler(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func Test_viewSecretHandler_viewConformPage_ProviderError_ExpectErr(t *testing.T) {
	t.Parallel()

	m := newMK(t)
	token := "test_token"

	m.secretsProvider.EXPECT().
		CheckExistsSecret(mock.Anything, token).
		Return(errors.New("test error")).
		Once()

	rec := httptest.NewRecorder()
	req := newTestRequestWithChiContext(t, http.MethodGet, "/get/"+token, map[string]string{"token": token})

	m.server.viewSecretHandler(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func Test_viewSecretHandler_viewSecretPage_Exists_ExpectOK(t *testing.T) {
	t.Parallel()

	m := newMK(t)
	token := "test_token"
	secret := "my_secret"

	m.secretsProvider.EXPECT().
		GetSecret(mock.Anything, token).
		Return(secret, nil).
		Once()

	m.template.EXPECT().
		RenderFile("secret.html", mock.Anything, template.Data{"secret": secret}).
		Return(nil).
		Once()

	rec := httptest.NewRecorder()
	req := newTestRequestWithChiContext(t, http.MethodGet, "/get/"+token+"?conform=true", map[string]string{"token": token})

	m.server.viewSecretHandler(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func Test_viewSecretHandler_viewSecretPage_NotFound_ExpectNotFound(t *testing.T) {
	t.Parallel()

	m := newMK(t)
	token := "test_token"

	m.secretsProvider.EXPECT().
		GetSecret(mock.Anything, token).
		Return("", storage.ErrNotFound).
		Once()

	m.template.EXPECT().
		RenderFile("expired.html", mock.Anything, template.Data(nil)).
		Return(nil).
		Once()

	rec := httptest.NewRecorder()
	req := newTestRequestWithChiContext(t, http.MethodGet, "/get/"+token+"?conform=true", map[string]string{"token": token})

	m.server.viewSecretHandler(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func Test_viewSecretHandler_viewSecretPage_ProviderError_ExpectErr(t *testing.T) {
	t.Parallel()

	m := newMK(t)
	token := "test_token"

	m.secretsProvider.EXPECT().
		GetSecret(mock.Anything, token).
		Return("", errors.New("test error")).
		Once()

	rec := httptest.NewRecorder()
	req := newTestRequestWithChiContext(t, http.MethodGet, "/get/"+token+"?conform=true", map[string]string{"token": token})

	m.server.viewSecretHandler(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
}
