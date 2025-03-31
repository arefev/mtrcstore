package service

import (
	"context"
	"fmt"

	"github.com/arefev/mtrcstore/internal/proto"
	"github.com/arefev/mtrcstore/internal/server/model"
	"github.com/arefev/mtrcstore/internal/server/repository"
)

type GRPCServer struct {
	proto.UnimplementedMetricsServer
	Storage repository.Storage
}

func (gs *GRPCServer) UpdateMetric(ctx context.Context, in *proto.UpdateMetricRequest) (*proto.UpdateMetricResponse, error) {
	var metrics []model.Metric
	for _, m := range in.Metrics {
		metrics = append(metrics, model.Metric{
			MType: m.Type,
			ID:    m.ID,
			Value: &m.Value,
			Delta: &m.Delta,
		})
	}

	err := gs.Storage.MassSave(ctx, metrics)
	if err != nil {
		return &proto.UpdateMetricResponse{}, fmt.Errorf("grpc update metric mass save failed: %w", err)
	}

	return &proto.UpdateMetricResponse{}, nil
}
