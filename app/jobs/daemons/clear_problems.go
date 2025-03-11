package daemons

import (
	"context"
	"educ-gpt/services"
	"time"

	"go.uber.org/zap"
)

type ClearProblemsDaemon struct {
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
	go func() {
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
			roadmapSrv: roadmapSrv,
			logger:     logger,
			sleep:      sleep,
			ctx:        ctx,
			cancel:     cancel,
		},
		ctx,
		cancel
}
