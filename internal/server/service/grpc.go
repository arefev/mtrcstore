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

func (gs *GRPCServer) UpdateMetric(
	ctx context.Context,
	in *proto.UpdateMetricRequest,
) (*proto.UpdateMetricResponse, error) {
	metrics := make([]model.Metric, 0, len(in.GetMetrics()))
	for _, m := range in.GetMetrics() {
		value := m.GetValue()
		delta := m.GetDelta()
		metrics = append(metrics, model.Metric{
			MType: m.GetType(),
			ID:    m.GetID(),
			Value: &value,
			Delta: &delta,
		})
	}

	err := gs.Storage.MassSave(ctx, metrics)
	if err != nil {
		return &proto.UpdateMetricResponse{}, fmt.Errorf("grpc update metric mass save failed: %w", err)
	}

	return &proto.UpdateMetricResponse{}, nil
}
