package service

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

var (
	ErrInvalidUrl   = errors.New("path url must be view as /update/{type}/{name}/{value}")
	ErrInvalidValue = errors.New("value is invalid")
	ErrInvalidType  = errors.New("type is invalid")
)

type UrlParams struct {
	Type  string
	Name  string
	Value float64
}

type UrlParser struct {
	Url *url.URL
}

func NewUrlParser(url *url.URL) UrlParser {
	return UrlParser{
		Url: url,
	}
}

func (p *UrlParser) Exec() (UrlParams, error) {
	parsed, error := p.parse()

	if error != nil {
		return UrlParams{}, error
	}

	params, error := p.getParams(parsed)
	if error != nil {
		return UrlParams{}, error
	}

	return params, error
}

func (p *UrlParser) parse() ([]string, error) {
	const pathSize = 4
	var parsed = make([]string, pathSize)
	parsed = strings.Split(p.Url.Path, "/")
	if len(parsed) > 1 {
		parsed = parsed[1:]
	}

	if len(parsed) != pathSize {
		return parsed, ErrInvalidUrl
	}

	return parsed, nil
}

func (p *UrlParser) getParams(parsed []string) (UrlParams, error) {
	value, err := strconv.ParseFloat(parsed[3], 64)
	if err != nil {
		return UrlParams{}, ErrInvalidValue
	}

	if err := p.checkType(parsed[1]); err != nil {
		return UrlParams{}, ErrInvalidType
	}

	param := UrlParams{
		Type:  parsed[1],
		Name:  parsed[2],
		Value: value,
	}

	return param, nil
}

func (p *UrlParser) checkType(t string) error {
	if t != "counter" && t != "gauge" {
		return ErrInvalidType
	}

	return nil
}
