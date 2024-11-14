package service

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
)

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func NewCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	res, err := c.zw.Write(p)
	if err != nil {
		return res, fmt.Errorf("compressWriter Write failed: %w", err)
	}
	return res, nil
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < http.StatusMultipleChoices {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	if err := c.zw.Close(); err != nil {
		return fmt.Errorf("compressWriter Close failed: %w", err)
	}
	return nil
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func NewCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("gzip reader failed: %w", err)
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	n, err = c.zr.Read(p)
	if err != nil {
		return n, fmt.Errorf("compressReader Read failed: %w", err)
	}
	return n, nil
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return fmt.Errorf("compressReader Close failed: %w", err)
	}
	return nil
}
