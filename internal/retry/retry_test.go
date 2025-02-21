package retry

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigSuccess(t *testing.T) {
	t.Run("retry success with error", func(t *testing.T) {
		r := New(func() error { return errors.New("test") }, func(err error) bool { return true }, 2)
		require.Error(t, r.Run())
	})
}
