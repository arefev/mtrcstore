package service

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/arefev/mtrcstore/internal/agent/model"
	"github.com/stretchr/testify/require"
)

func TestDoRequestSuccess(t *testing.T) {
	t.Run("do request success", func(t *testing.T) {
		ctx := context.Background()
		s := httptest.NewServer(nil)
		defer s.Close()

		client := NewClient("", "", s.URL)
		err := client.Request(ctx, []model.Metric{})
		require.NoError(t, err)
	})
}

func TestDoRequestFail(t *testing.T) {
	t.Run("do request success", func(t *testing.T) {
		ctx := context.Background()
		client := NewClient("", "", "http://fail.lo")
		err := client.Request(ctx, []model.Metric{})
		require.ErrorIs(t, err, ErrRequestFail)
	})
}
