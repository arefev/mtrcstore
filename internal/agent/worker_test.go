package agent

import (
	"runtime"
	"testing"

	"github.com/arefev/mtrcstore/internal/agent/repository"
	"github.com/stretchr/testify/assert"
)

func TestWorker_read(t *testing.T) {
	// var memStats runtime.MemStats
	storage := repository.NewMemory()
	serverHost := "http://localhost:8080"
	type fields struct {
		ReportInterval float64
		PollInterval   int
		Storage        repository.Storage
		ServerHost     string
	}
	type args struct {
		memStats *runtime.MemStats
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name: "positive test №1",
			fields: fields{
				ReportInterval: 2,
				PollInterval:  0,
				Storage:        &storage,
				ServerHost:     serverHost,
			},
			args: args{
				memStats: &runtime.MemStats{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w := &Worker{
				ReportInterval: tt.fields.ReportInterval,
				PollInterval:   tt.fields.PollInterval,
				Storage:        tt.fields.Storage,
				ServerHost:     tt.fields.ServerHost,
			}
			w.read(tt.args.memStats)

			assert.Contains(t, w.Storage.GetGauges(), "Alloc")
		})
	}
}
