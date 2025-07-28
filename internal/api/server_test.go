package api

import (
	"testing"

	"enigma/internal/api/mocks"
	mocks_template "enigma/internal/pkg/adapters/template/mocks"
)

const (
	externalURL = "http://localhost"
)

type mk struct {
	secretsProvider *mocks.SecretsProvider
	template        *mocks_template.Template
	server          *Server
}

func newMK(t *testing.T) *mk {
	secretsProvider := mocks.NewSecretsProvider(t)
	template := mocks_template.NewTemplate(t)
	return &mk{
		secretsProvider: secretsProvider,
		template:        template,
		server:          NewServer(secretsProvider, template, externalURL),
	}
}
