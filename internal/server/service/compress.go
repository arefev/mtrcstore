package service

import (
	"compress/gzip"
	"fmt"
	"net/http"
)

type compressWriter struct {
	http.ResponseWriter
	zw          *gzip.Writer
	wroteHeader bool
}

func NewCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		ResponseWriter: w,
		zw:             gzip.NewWriter(w),
		wroteHeader:    false,
	}
}

func (c *compressWriter) Write(p []byte) (int, error) {
	if !c.wroteHeader {
		c.WriteHeader(http.StatusOK)
	}
	res, err := c.zw.Write(p)
	if err != nil {
		return res, fmt.Errorf("compressWriter Write failed: %w", err)
	}
	return res, nil
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < http.StatusMultipleChoices {
		c.Header().Set("Content-Encoding", "gzip")
	}
	c.ResponseWriter.WriteHeader(statusCode)
	c.wroteHeader = true
}

func (c *compressWriter) Close() error {
	if err := c.zw.Close(); err != nil {
		return fmt.Errorf("compressWriter Close failed: %w", err)
	}
	return nil
}
