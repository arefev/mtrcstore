package service

import (
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
)

var ErrRequestFail = errors.New("doRequest failed")

type Client struct{}

func (c *Client) DoRequest(url string, headers map[string]string, body any) error {
	request := resty.New().R()
	for k, v := range headers {
		request.SetHeader(k, v)
	}

	if _, err := request.SetBody(body).Post(url); err != nil {
		return fmt.Errorf("%w: %w", ErrRequestFail, err)
	}

	return nil
}
