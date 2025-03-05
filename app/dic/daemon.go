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
		time.Hour*24,
	)
	daemonsSlice = append(daemonsSlice, daemon)
	ClearProblemsDaemon = &DaemonController{ctx, cancel}
}

func initServices() {
	clearProblems()
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
