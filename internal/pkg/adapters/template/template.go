//go:generate mockery --case snake --with-expecter --name Template

package template

import "io"

// Data ...
type Data map[string]any

// Template ...
type Template interface {
	RenderFile(name string, writer io.Writer, data Data) error
}
