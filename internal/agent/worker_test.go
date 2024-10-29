package agent

import (
	"runtime"
	"testing"

	"github.com/arefev/mtrcstore/internal/agent/repository"
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
				PollInterval:  0,
			},
			args: args{
				memStats: &runtime.MemStats{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			storage := repository.NewMemory()
			serverHost := "http://localhost:8080"

			w := &Worker{
				ReportInterval: tt.fields.ReportInterval,
				PollInterval:   tt.fields.PollInterval,
				Storage:        &storage,
				ServerHost:     serverHost,
			}
			w.read(tt.args.memStats)
			w.Storage.IncrementCounter()

			assert.Contains(t, w.Storage.GetGauges(), "Alloc")
			assert.Contains(t, w.Storage.GetCounters(), "PollCount")
		})
	}
}
