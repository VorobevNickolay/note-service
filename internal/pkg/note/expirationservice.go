package note

import (
	"go.uber.org/zap"
	"time"
)

type expStore interface {
	ExpireNotes() error
}

type ExpService struct {
	store  expStore
	ticker *time.Ticker
	logger *zap.Logger
}

func NewExpService(store expStore, expirationRunInterval time.Duration, logger *zap.Logger) *ExpService {
	ticker := time.NewTicker(expirationRunInterval)
	return &ExpService{store: store, ticker: ticker, logger: logger}
}

func (service *ExpService) Run() error {
	for range service.ticker.C {
		if err := service.store.ExpireNotes(); err != nil {
			service.logger.Error("failed to expire notes", zap.Error(err))
		}
	}
	return nil
}
