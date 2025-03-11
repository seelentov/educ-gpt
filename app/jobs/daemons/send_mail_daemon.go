package daemons

import (
	"context"
	"educ-gpt/jobs/tasks"
	"educ-gpt/services"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"

	"go.uber.org/zap"
)

type SendMailDaemon struct {
	senderSrv   services.SenderService
	redisClient *redis.Client
	queueName   string
	logger      *zap.Logger
	sleep       time.Duration
	ctx         context.Context
	cancel      context.CancelFunc
}

func (c *SendMailDaemon) Work() {
	result, err := c.redisClient.BLPop(c.ctx, 0, c.queueName).Result()
	if err != nil {
		c.logger.Error("Error while receiving email sender", zap.Error(err))
		return
	}

	var task tasks.MailTask
	err = json.Unmarshal([]byte(result[1]), &task)
	if err != nil {
		c.logger.Error("Error while parsing email sender", zap.Error(err))
		return
	}

	err = c.senderSrv.SendMessage(task.To, task.Subject, task.Body)
	if err != nil {
		c.logger.Error("Error while sending email sender. Retry after all", zap.Error(err))
		c.redisClient.LPush(c.ctx, c.queueName, result[1])
		return
	}

	c.logger.Info("Sent email", zap.String(c.queueName, task.To))
}

func (c *SendMailDaemon) Start() {
	go func() {
		c.logger.Info("SendMailDaemon started")
		for {
			select {
			case <-c.ctx.Done():
				c.logger.Info("SendMailDaemon stopped")
				return
			default:
				c.logger.Info("SendMailDaemon works")
				c.Work()
				c.logger.Info("SendMailDaemon works end")
			}
		}
	}()
}

func (c *SendMailDaemon) Stop() {
	c.logger.Info("Stopping SendMailDaemon...")
	c.cancel()
}

func NewSendMailDaemon(
	senderSrv services.SenderService,
	redisClient *redis.Client,
	logger *zap.Logger,
	sleep time.Duration,
	queueName string,
) (*SendMailDaemon, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	return &SendMailDaemon{
			senderSrv:   senderSrv,
			redisClient: redisClient,
			logger:      logger,
			sleep:       sleep,
			ctx:         ctx,
			cancel:      cancel,
			queueName:   queueName,
		},
		ctx,
		cancel
}
