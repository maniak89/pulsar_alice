package wrap_logger

import (
	"context"
	"strconv"
	"sync"

	"pulsar_alice/internal/meter"
	"pulsar_alice/internal/models/storage"
)

type wrapper struct {
	child   meter.ValueProvider
	logger  Logger
	meterID string
	isInit  bool
	isInitM sync.Mutex
}

type Logger interface {
	Log(ctx context.Context, meterID string, level storage.LogLevel, msg string)
}

func New(child meter.ValueProvider, meterID string, logger Logger) meter.ValueProvider {
	return &wrapper{
		child:   child,
		logger:  logger,
		meterID: meterID,
	}
}

func (w *wrapper) insure(ctx context.Context) error {
	w.isInitM.Lock()
	defer w.isInitM.Unlock()
	if w.isInit {
		return nil
	}
	if err := w.child.Init(ctx); err != nil {
		w.logger.Log(ctx, w.meterID, storage.Error, err.Error())
		return err
	}
	w.logger.Log(ctx, w.meterID, storage.Info, "Success connected")
	w.isInit = true
	return nil
}

func (w *wrapper) Init(ctx context.Context) error {
	return w.insure(ctx)
}

func (w *wrapper) Value(ctx context.Context) meter.Value {
	if err := w.insure(ctx); err != nil {
		return meter.Value{Error: err}
	}

	result := w.child.Value(ctx)
	if result.Error != nil {
		w.logger.Log(ctx, w.meterID, storage.Error, "Failed get values: "+result.Error.Error())

		return result
	}

	w.logger.Log(ctx, w.meterID, storage.Info, "Success get value. Total "+strconv.FormatFloat(result.Value, 'f', -1, 64))

	return result
}
