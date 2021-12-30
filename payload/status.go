package payload

type TaskStatus string

const (
	TaskRunning = TaskStatus("RUNNING")
	TaskFailed  = TaskStatus("FAILED")
)

// ResultConnectorStatus for goroutine
type ResultConnectorStatus struct {
	Err               error
	ConnectorStatus   ConnectorStatus
	ConnectorNotFound ConnectorNotFound
}

type ConnectorStatus struct {
	Name      string    `json:"name"`
	Connector Connector `json:"connector"`
	Tasks     []Task    `json:"tasks"`
}

type Connector struct {
	Status   string `json:"status"`
	WorkerID string `json:"worker_id"`
}

type Task struct {
	ID       int32  `json:"id"`
	State    string `json:"state"`
	WorkerID string `json:"worker_id"`
	Trace    string `json:"trace"`
}

// ConnectorNotFound connector name not found
type ConnectorNotFound struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}

func (t Task) IsFailed() bool {
	if TaskStatus(t.State) == TaskFailed {
		return true
	}
	return false
}

func (t Task) ErrorMessage() string {
	if t.IsFailed() {
		return t.Trace
	}
	return ""
}

func (c Connector) IsFailed() bool {
	if TaskStatus(c.Status) == TaskFailed {
		return true
	}
	return false
}
