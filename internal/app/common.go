package app

type TaskState int

const (
	TaskStateUnknown = iota
	TaskStateInit    = 1
	TaskStateDoing   = 2
	TaskStateDone    = 3
	TaskStateFailed  = 4
)
