package daemons

import (
	"context"
	"educ-gpt/services"
	"sync"
	"time"

	"go.uber.org/zap"
)

type ClearProblemsDaemon struct {
	wg         *sync.WaitGroup
	roadmapSrv services.RoadmapService
	logger     *zap.Logger
	sleep      time.Duration
	ctx        context.Context
	cancel     context.CancelFunc
}

func (c *ClearProblemsDaemon) Work() {
	err := c.roadmapSrv.ClearProblems()
	if err != nil {
		c.logger.Error("Error clearing problems", zap.Error(err))
	}
}

func (c *ClearProblemsDaemon) Start() {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.logger.Info("ClearProblemsDaemon started")
		ticker := time.NewTicker(c.sleep)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.logger.Info("ClearProblemsDaemon works")
				c.Work()
				c.logger.Info("ClearProblemsDaemon works end")
			case <-c.ctx.Done():
				c.logger.Info("ClearProblemsDaemon stopped")
				return
			}
		}
	}()
}

func (c *ClearProblemsDaemon) Stop() {
	c.logger.Info("Stopping ClearProblemsDaemon...")
	c.cancel()
}

func NewClearProblemsDaemon(
	roadmapSrv services.RoadmapService,
	logger *zap.Logger,
	sleep time.Duration,
) (*ClearProblemsDaemon, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	return &ClearProblemsDaemon{
			wg:         &sync.WaitGroup{},
			roadmapSrv: roadmapSrv,
			logger:     logger,
			sleep:      sleep,
			ctx:        ctx,
			cancel:     cancel,
		},
		ctx,
		cancel
}
