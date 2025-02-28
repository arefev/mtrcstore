package service

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDoRequestSuccess(t *testing.T) {
	t.Run("do request success", func(t *testing.T) {
		s := httptest.NewServer(nil)
		defer s.Close()

		client := Client{}
		err := client.DoRequest(s.URL, map[string]string{"Content-type": "application/json"}, nil)
		require.NoError(t, err)
	})
}

func TestDoRequestFail(t *testing.T) {
	t.Run("do request success", func(t *testing.T) {
		client := Client{}
		err := client.DoRequest("http://fail.lo", map[string]string{}, nil)
		require.ErrorIs(t, err, ErrRequestFail)
	})
}
