package agent

import (
	"runtime"
	"testing"

	"github.com/arefev/mtrcstore/internal/agent/repository"
	"github.com/arefev/mtrcstore/internal/agent/service"
	"github.com/stretchr/testify/assert"
)

func TestWorker_read(t *testing.T) {
	type fields struct {
		ReportInterval int
		PollInterval   int
	}
	type args struct {
		memStats *runtime.MemStats
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "positive test №1",
			fields: fields{
				ReportInterval: 2,
				PollInterval:   0,
			},
			args: args{
				memStats: &runtime.MemStats{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serverHost := "http://localhost:8080"
			storage := repository.NewMemory()
			report := service.NewReport(&storage, serverHost)
			w := Worker{
				Report:         &report,
				ReportInterval: tt.fields.ReportInterval,
				PollInterval:   tt.fields.PollInterval,
			}

			assert.NoError(t, w.read(tt.args.memStats))
			w.Report.IncrementCounter()

			assert.Contains(t, storage.GetGauges(), "Alloc")
			assert.Contains(t, storage.GetCounters(), "PollCount")
		})
	}
}
