package service

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/arefev/mtrcstore/internal/agent/model"
	"github.com/arefev/mtrcstore/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcClient struct {
	url string
}

func NewGRPCClient(url string) *grpcClient {
	return &grpcClient{
		url: url,
	}
}

func (gc *grpcClient) Request(ctx context.Context, data []model.Metric) error {
	conn, err := grpc.NewClient(gc.url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("grpc request NewClient failed: %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("grpc request conn close failed: %s", err.Error())
		}
	}()
	client := proto.NewMetricsClient(conn)

	pMetrics := make([]*proto.Metric, 0, len(data))
	for _, m := range data {
		pm := &proto.Metric{
			ID:   m.ID,
			Type: m.MType,
		}

		if m.Value != nil {
			pm.Value = *m.Value
		}

		if m.Delta != nil {
			pm.Delta = *m.Delta
		}

		pMetrics = append(pMetrics, pm)
	}

	_, err = client.UpdateMetric(ctx, &proto.UpdateMetricRequest{
		Metrics: pMetrics,
	})

	if err != nil {
		return fmt.Errorf("grpc request UpdateMetric failed: %w", err)
	}

	return nil
}

func (gc *grpcClient) IsConnRefused(err error) bool {
	return strings.Contains(err.Error(), "connection refused")
}
