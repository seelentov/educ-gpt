package daemons

import (
	"context"
	"educ-gpt/services"
	"time"

	"go.uber.org/zap"
)

type ClearTokensDaemon struct {
	tokenSrv services.TokenService
	logger   *zap.Logger
	sleep    time.Duration
	ctx      context.Context
	cancel   context.CancelFunc
}

func (c *ClearTokensDaemon) Work() {
	err := c.tokenSrv.Clear()
	if err != nil {
		c.logger.Error("Error clearing tokens", zap.Error(err))
	}
}

func (c *ClearTokensDaemon) Start() {
	go func() {
		c.logger.Info("ClearTokensDaemon started")
		ticker := time.NewTicker(c.sleep)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.logger.Info("ClearTokensDaemon works")
				c.Work()
				c.logger.Info("ClearTokensDaemon works end")
			case <-c.ctx.Done():
				c.logger.Info("ClearTokensDaemon stopped")
				return
			}
		}
	}()
}

func (c *ClearTokensDaemon) Stop() {
	c.logger.Info("Stopping ClearTokensDaemon...")
	c.cancel()
}

func NewClearTokensDaemon(
	tokenSrv services.TokenService,
	logger *zap.Logger,
	sleep time.Duration,
) (*ClearTokensDaemon, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	return &ClearTokensDaemon{
			tokenSrv: tokenSrv,
			logger:   logger,
			sleep:    sleep,
			ctx:      ctx,
			cancel:   cancel,
		},
		ctx,
		cancel
}
