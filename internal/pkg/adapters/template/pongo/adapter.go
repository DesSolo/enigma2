package pongo

import (
	"fmt"
	"io"

	"github.com/flosch/pongo2/v6"

	"enigma/internal/pkg/adapters/template"
)

// Adapter ...
type Adapter struct {
	templateSet *pongo2.TemplateSet
}

// NewAdapter ...
func NewAdapter(templateSet *pongo2.TemplateSet) *Adapter {
	return &Adapter{
		templateSet: templateSet,
	}
}

// RenderFile ...
func (a *Adapter) RenderFile(name string, writer io.Writer, data template.Data) error {
	tmpl, err := a.templateSet.FromFile(name)
	if err != nil {
		return fmt.Errorf("could not load template: %w", err)
	}

	if err := tmpl.ExecuteWriter(pongo2.Context(data), writer); err != nil {
		return fmt.Errorf("could not render template: %w", err)
	}

	return nil
}
