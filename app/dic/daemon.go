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

var ClearTokensDaemon *DaemonController

func clearTokens() {
	daemon, ctx, cancel := daemons.NewClearTokensDaemon(
		TokenService(),
		logger.Logger(),
		time.Hour*2,
	)
	daemonsSlice = append(daemonsSlice, daemon)
	ClearTokensDaemon = &DaemonController{ctx, cancel}
}

var ClearNonActivatedUsersDaemon *DaemonController

func clearNonActivatedUsers() {
	daemon, ctx, cancel := daemons.NewClearNonActivatedUsersDaemon(
		UserService(),
		logger.Logger(),
		time.Hour*2,
	)
	daemonsSlice = append(daemonsSlice, daemon)
	ClearNonActivatedUsersDaemon = &DaemonController{ctx, cancel}
}

func initServices() {
	clearProblems()
	clearTokens()
	clearNonActivatedUsers()
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
