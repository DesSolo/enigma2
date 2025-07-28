package closer

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func newTestHandler(t *testing.T, err error) handler {
	t.Helper()

	return func() error {
		return err
	}
}

func Test_Global_Add(t *testing.T) {
	Add(newTestHandler(t, nil))
	require.NotEmpty(t, globalCloser.handlers)
}

func Test_Global_Close_ExpectOk(t *testing.T) {
	t.Parallel()

	globalCloser = New()

	Add(newTestHandler(t, nil))

	err := Close()
	require.NoError(t, err)
}

func Test_Global_Close_ExpectErr(t *testing.T) {
	t.Parallel()

	globalCloser = New()

	Add(newTestHandler(t, errors.New("test error")))

	err := Close()
	require.EqualError(t, err, "test error")
}
