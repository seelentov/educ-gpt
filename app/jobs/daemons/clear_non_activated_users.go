package daemons

import (
	"context"
	"educ-gpt/services"
	"time"

	"go.uber.org/zap"
)

type ClearNonActivatedUsersDaemon struct {
	userSrv services.UserService
	logger  *zap.Logger
	sleep   time.Duration
	ctx     context.Context
	cancel  context.CancelFunc
}

func (c *ClearNonActivatedUsersDaemon) Work() {
	err := c.userSrv.ClearNonActivatedUsers()
	if err != nil {
		c.logger.Error("Error clearing non activated users", zap.Error(err))
	}
}

func (c *ClearNonActivatedUsersDaemon) Start() {
	go func() {
		c.logger.Info("ClearNonActivatedUsersDaemon started")
		ticker := time.NewTicker(c.sleep)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.logger.Info("ClearNonActivatedUsersDaemon works")
				c.Work()
				c.logger.Info("ClearNonActivatedUsersDaemon works end")
			case <-c.ctx.Done():
				c.logger.Info("ClearNonActivatedUsersDaemon stopped")
				return
			}
		}
	}()
}

func (c *ClearNonActivatedUsersDaemon) Stop() {
	c.logger.Info("Stopping ClearNonActivatedUsersDaemon...")
	c.cancel()
}

func NewClearNonActivatedUsersDaemon(
	userSrv services.UserService,
	logger *zap.Logger,
	sleep time.Duration,
) (*ClearNonActivatedUsersDaemon, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	return &ClearNonActivatedUsersDaemon{
			userSrv: userSrv,
			logger:  logger,
			sleep:   sleep,
			ctx:     ctx,
			cancel:  cancel,
		},
		ctx,
		cancel
}
