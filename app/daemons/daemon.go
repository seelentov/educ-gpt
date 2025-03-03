package daemons

type DaemonService interface {
	Start()
	Work()
	Stop()
}
