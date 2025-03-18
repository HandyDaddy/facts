package application

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/HandyDaddy/facts/internal/config"
	"github.com/HandyDaddy/facts/internal/domain/entities"
	"github.com/HandyDaddy/facts/internal/domain/repository"
	factclient "github.com/HandyDaddy/facts/internal/infrastructure/fact-client"
)

type FactService struct {
	queue      chan entities.Fact
	wg         *sync.WaitGroup
	stopChan   chan struct{}
	cfg        *config.Config
	factClient repository.FactRepository
	logger     *logrus.Logger
}

func NewFactService(cfg *config.Config) *FactService {
	return &FactService{
		queue:      make(chan entities.Fact, cfg.App.MaxBufferSize),
		wg:         &sync.WaitGroup{},
		factClient: factclient.NewHTTPClient(&cfg.HttpClient),
		stopChan:   make(chan struct{}),
		cfg:        cfg,
		logger:     logrus.New(),
	}
}

// Start - starting FactService
func (fs *FactService) Start(ctx context.Context) {
	fs.wg.Add(1)
	go func() {
		defer fs.wg.Done()
		fs.processQueue(ctx)
	}()

	fs.wg.Add(1)
	go func() {
		defer fs.wg.Done()
		fs.simulateIncomingData(ctx)
	}()
}

// Shutdown - graceful shutdown
func (fs *FactService) Shutdown() {
	fs.logger.Info("Starting graceful shutdown...")
	close(fs.stopChan)

	timeout := time.After(30 * time.Second)
	for {
		select {
		case <-timeout:
			fs.logger.Warn("Shutdown timeout reached, forcing exit")
			return
		default:
			if len(fs.queue) == 0 {
				fs.wg.Wait()
				close(fs.queue)
				fs.logger.Info("Graceful shutdown completed")
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// AddFact - add fact to queue
func (fs *FactService) AddFact(fact entities.Fact) error {
	select {
	case fs.queue <- fact:
		fs.logger.Infof("Fact added to queue. Queue length: %d", len(fs.queue))
		return nil
	case <-fs.stopChan:
		return fmt.Errorf("service is shutting down")
	}
}

func (fs *FactService) processQueue(ctx context.Context) {
	for {
		select {
		case fact := <-fs.queue:
			if err := fs.factClient.SaveFact(ctx, &fact); err != nil {
				fs.logger.Errorf("Failed to save fact: %v", err)
				// При ошибке возвращаем факт в очередь
				select {
				case fs.queue <- fact:
				case <-fs.stopChan:
					return
				case <-ctx.Done():
					return
				}
				time.Sleep(time.Second)
				continue
			}
			fs.logger.Info("Fact successfully saved")
		case <-fs.stopChan:
			return
		case <-ctx.Done():
			return
		}
	}
}

func (fs *FactService) simulateIncomingData(ctx context.Context) {
	fact := entities.Fact{
		PeriodStart:         "2024-05-01",
		PeriodEnd:           "2024-05-31",
		PeriodKey:           "year",
		IndicatorToMoId:     227373,
		IndicatorToMoFactId: 0,
		Value:               1,
		FactTime:            "2024-05-31",
		IsPlan:              0,
		AuthUserId:          40,
		Comment:             "buffer Last_name",
	}

	ticker := time.NewTicker(60 * time.Millisecond) // ~1000 записей в минуту
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := fs.AddFact(fact); err != nil {
				fs.logger.Errorf("Failed to add fact: %v", err)
			}
		case <-fs.stopChan:
			return
		case <-ctx.Done():
			return
		}
	}
}
