package daemons

import (
	"context"
	"educ-gpt/services"
	"time"

	"go.uber.org/zap"
)

type ClearThemesDaemon struct {
	roadmapSrv services.RoadmapService
	logger     *zap.Logger
	sleep      time.Duration
	ctx        context.Context
	cancel     context.CancelFunc
}

func (c *ClearThemesDaemon) Work() {
	err := c.roadmapSrv.ClearThemes()
	if err != nil {
		c.logger.Error("Error clearing themes", zap.Error(err))
	}
}

func (c *ClearThemesDaemon) Start() {
	go func() {
		c.logger.Info("ClearThemesDaemon started")
		ticker := time.NewTicker(c.sleep)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.logger.Info("ClearThemesDaemon works")
				c.Work()
				c.logger.Info("ClearThemesDaemon works end")
			case <-c.ctx.Done():
				c.logger.Info("ClearThemesDaemon stopped")
				return
			}
		}
	}()
}

func (c *ClearThemesDaemon) Stop() {
	c.logger.Info("Stopping ClearThemesDaemon...")
	c.cancel()
}

func NewClearThemesDaemon(
	roadmapSrv services.RoadmapService,
	logger *zap.Logger,
	sleep time.Duration,
) (*ClearThemesDaemon, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	return &ClearThemesDaemon{
			roadmapSrv: roadmapSrv,
			logger:     logger,
			sleep:      sleep,
			ctx:        ctx,
			cancel:     cancel,
		},
		ctx,
		cancel
}
