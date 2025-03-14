package service

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDoRequestSuccess(t *testing.T) {
	t.Run("do request success", func(t *testing.T) {
		ctx := context.Background()
		s := httptest.NewServer(nil)
		defer s.Close()

		client := Client{}
		err := client.DoRequest(ctx, s.URL, map[string]string{"Content-type": "application/json"}, nil)
		require.NoError(t, err)
	})
}

func TestDoRequestFail(t *testing.T) {
	t.Run("do request success", func(t *testing.T) {
		ctx := context.Background()
		client := Client{}
		err := client.DoRequest(ctx, "http://fail.lo", map[string]string{}, nil)
		require.ErrorIs(t, err, ErrRequestFail)
	})
}
