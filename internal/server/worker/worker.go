package worker

import (
	"encoding/json"
	"os"
	"time"

	"github.com/arefev/mtrcstore/internal/server/repository"
	"go.uber.org/zap"
)

const filePermission = 0o644

var Worker *workerStore

type workerStore struct {
	Storage         repository.Storage
	FileStoragePath string
	log             *zap.Logger
	StoreInterval   int
	storeByEvent    bool
	restore         bool
}

func Init(interval int, fileStoragePath string, restore bool, storage repository.Storage, log *zap.Logger) *workerStore {
	if Worker != nil {
		return Worker
	}

	Worker = &workerStore{
		StoreInterval:   interval,
		Storage:         storage,
		FileStoragePath: fileStoragePath,
		storeByEvent:    interval == 0,
		restore:         restore,
		log:             log,
	}

	if restore {
		Worker.load()
	}

	return Worker
}

func (w *workerStore) Run() {
	w.log.Info(
		"worker running with params",
		zap.Int("interval in seconds", w.StoreInterval),
		zap.String("file", w.FileStoragePath),
		zap.Bool("restore", w.restore),
		zap.Bool("storeByEvent", w.storeByEvent),
	)

	if w.storeByEvent {
		return
	}

	start := time.Now()
	for {
		period := int(time.Since(start).Seconds())
		if period > w.StoreInterval {
			w.save()
			start = time.Now()
		}
	}
}

func (w *workerStore) load() {
	file, err := os.OpenFile(w.FileStoragePath, os.O_RDONLY|os.O_CREATE, filePermission)
	if err != nil {
		w.log.Error("worker open file failed", zap.Error(err))
		return
	}

	r := json.NewDecoder(file)

	if err := r.Decode(&w.Storage); err != nil {
		w.log.Error("worker decode data failed", zap.Error(err))
		return
	}

	w.log.Info("worker data loaded")
}

func (w *workerStore) SaveEvent() {
	if !w.storeByEvent {
		return
	}

	w.save()
}

func (w *workerStore) save() {
	file, err := os.OpenFile(w.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, filePermission)
	if err != nil {
		w.log.Error("worker open file failed", zap.Error(err))
		return
	}

	defer func() {
		err := file.Close()
		if err != nil {
			w.log.Error("worker close file failed", zap.Error(err))
			return
		}
	}()

	wr := json.NewEncoder(file)
	if err := wr.Encode(w.Storage); err != nil {
		w.log.Error("worker encode data failed", zap.Error(err))
		return
	}

	w.log.Info("worker data saved")
}
