package daemons

import (
	"context"
	"educ-gpt/services"
	"sync"
	"time"

	"go.uber.org/zap"
)

type ClearResetTokens struct {
	wg             *sync.WaitGroup
	resetTokensSrv services.ResetTokenService
	logger         *zap.Logger
	sleep          time.Duration
	ctx            context.Context
	cancel         context.CancelFunc
}

func (c *ClearResetTokens) Work() {
	err := c.resetTokensSrv.Clear()
	if err != nil {
		c.logger.Error("Error clearing reset tokens", zap.Error(err))
	}
}

func (c *ClearResetTokens) Start() {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.logger.Info("ClearResetTokens started")
		ticker := time.NewTicker(c.sleep)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.logger.Info("ClearResetTokens works")
				c.Work()
				c.logger.Info("ClearResetTokens works end")
			case <-c.ctx.Done():
				c.logger.Info("ClearResetTokens stopped")
				return
			}
		}
	}()
}

func (c *ClearResetTokens) Stop() {
	c.logger.Info("Stopping ClearProblemsDaemon...")
	c.cancel()
}

func NewClearResetTokens(
	resetTokensSrv services.ResetTokenService,
	logger *zap.Logger,
	sleep time.Duration,
) (*ClearResetTokens, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	return &ClearResetTokens{
			wg:             &sync.WaitGroup{},
			resetTokensSrv: resetTokensSrv,
			logger:         logger,
			sleep:          sleep,
			ctx:            ctx,
			cancel:         cancel,
		},
		ctx,
		cancel
}
