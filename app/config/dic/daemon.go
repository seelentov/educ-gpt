package dic

import (
	"context"
	"educ-gpt/config/data"
	"educ-gpt/config/logger"
	"educ-gpt/jobs/daemons"
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

var ClearThemesDaemon *DaemonController

func clearThemes() {
	daemon, ctx, cancel := daemons.NewClearThemesDaemon(
		RoadmapService(),
		logger.Logger(),
		time.Hour*12,
	)
	daemonsSlice = append(daemonsSlice, daemon)
	ClearThemesDaemon = &DaemonController{ctx, cancel}
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

var SendMailDaemon *DaemonController

func sendMail() {
	daemon, ctx, cancel := daemons.NewSendMailDaemon(
		SenderService(),
		data.Redis(),
		logger.Logger(),
		0,
		"email_sender",
	)
	daemonsSlice = append(daemonsSlice, daemon)
	ClearNonActivatedUsersDaemon = &DaemonController{ctx, cancel}
}

var ClearUnusedFilesDaemon *DaemonController

func clearUnusedFiles() {
	daemon, ctx, cancel := daemons.NewClearUnusedFilesDaemon(
		data.DB(),
		logger.Logger(),
		time.Hour*24,
	)
	daemonsSlice = append(daemonsSlice, daemon)
	ClearUnusedFilesDaemon = &DaemonController{ctx, cancel}
}

func initServices() {
	clearProblems()
	clearTokens()
	clearNonActivatedUsers()
	clearUnusedFiles()
	clearThemes()

	sendMail()
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
