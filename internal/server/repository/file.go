package repository

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/arefev/mtrcstore/internal/server/model"
	"go.uber.org/zap"
)

type file struct {
	memory
	log             *zap.Logger
	fileStoragePath string
	filePermission  fs.FileMode
	storeInterval   int
	restore         bool
	storeByEvent    bool
}

func NewFile(intrvl int, filePath string, restore bool, log *zap.Logger) *file {
	const filePermission fs.FileMode = 0o644
	
	file := file{
		memory:          *NewMemory(),
		fileStoragePath: filePath,
		filePermission:  filePermission,
		storeInterval:   intrvl,
		restore:         restore,
		storeByEvent:    intrvl == 0,
		log:             log,
	}

	if restore {
		file.load()
	}

	return &file
}

func (f *file) load() {
	file, err := os.OpenFile(f.fileStoragePath, os.O_RDONLY|os.O_CREATE, f.filePermission)
	if err != nil {
		f.log.Error("worker open file failed", zap.Error(err))
		return
	}

	r := json.NewDecoder(file)

	if err := r.Decode(&f); err != nil {
		f.log.Error("worker decode data failed", zap.Error(err))
		return
	}

	f.log.Info("worker data loaded")
}

func (f *file) WorkerRun() *file {
	go f.worker()
	return f
}

func (f *file) worker() {
	f.log.Info(
		"worker running with params",
		zap.Int("interval in seconds", f.storeInterval),
		zap.String("file", f.fileStoragePath),
		zap.Bool("restore", f.restore),
		zap.Bool("storeByEvent", f.storeByEvent),
	)

	if f.storeByEvent {
		return
	}

	start := time.Now()
	for {
		period := int(time.Since(start).Seconds())
		if period > f.storeInterval {
			f.write()
			start = time.Now()
		}
	}
}

func (f *file) write() {
	file, err := os.OpenFile(f.fileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.filePermission)
	if err != nil {
		f.log.Error("worker open file failed", zap.Error(err))
		return
	}

	defer func() {
		err := file.Close()
		if err != nil {
			f.log.Error("worker close file failed", zap.Error(err))
			return
		}
	}()

	wr := json.NewEncoder(file)
	if err := wr.Encode(f); err != nil {
		f.log.Error("worker encode data failed", zap.Error(err))
		return
	}

	f.log.Info("worker data saved by worker")
}

func (f *file) Save(m model.Metric) error {
	if err := f.memory.Save(m); err != nil {
		return err
	}

	return f.writeEvent()
}

func (f *file) writeEvent() error {
	if !f.storeByEvent {
		return nil
	}

	file, err := os.OpenFile(f.fileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.filePermission)
	if err != nil {
		return fmt.Errorf("worker open file failed: %w", err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			f.log.Error("worker close file failed", zap.Error(err))
			return
		}
	}()

	wr := json.NewEncoder(file)
	if err := wr.Encode(f); err != nil {
		return fmt.Errorf("worker encode data failed: %w", err)
	}

	f.log.Info("worker data written by event")

	return nil
}
