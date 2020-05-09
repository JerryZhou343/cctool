package app

type TaskState int

const (
	TaskStateUnknown = iota
	TaskStateInit    = 1
	TaskStateDoing   = 2
	TaskStateDone    = 3
	TaskStateFailed  = 4
)

func (t TaskState) String() string {
	switch t {
	case TaskStateDoing:
		return "正在进行"
	case TaskStateDone:
		return "已完成"
	case TaskStateFailed:
		return "失败"
	default:
		return "未知状态"
	}
}
