package dic

import (
	"context"
	"educ-gpt/daemons"
	"educ-gpt/logger"
	"time"
)

type DaemonController struct {
	Ctx    context.Context
	Cancel context.CancelFunc
}

var (
	daemonsSlice []daemons.DaemonService
)

var ClearProblemsDaemon *DaemonController

func clearProblems() {
	daemon, ctx, cancel := daemons.NewClearProblemsDaemon(
		RoadmapService(),
		logger.Logger(),
		time.Hour*12,
	)
	daemonsSlice = append(daemonsSlice, daemon)
	ClearProblemsDaemon = &DaemonController{ctx, cancel}
}

var ClearResetTokensDaemon *DaemonController

func clearResetTokens() {
	daemon, ctx, cancel := daemons.NewClearResetTokens(
		ResetTokenService(),
		logger.Logger(),
		time.Hour*2,
	)
	daemonsSlice = append(daemonsSlice, daemon)
	ClearResetTokensDaemon = &DaemonController{ctx, cancel}
}

func initServices() {
	clearProblems()
	clearResetTokens()
}

func InitDaemons() {
	initServices()
	startAllDaemons()
}

func startAllDaemons() {
	for i := range daemonsSlice {
		daemonsSlice[i].Start()
	}
}

func StopAllDaemons() {
	for i := range daemonsSlice {
		daemonsSlice[i].Stop()
	}
}
