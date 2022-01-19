package payload

import "testing"

func TestTask_IsFailed(t *testing.T) {
	task := Task{Trace: "", WorkerID: "testing1", ID: 1, State: "RUNNING"}
	if task.IsFailed() {
		t.Error("return of state is not as expected")
	}
	task = Task{Trace: "", WorkerID: "testing2", ID: 2, State: "FAILED"}
	if !task.IsFailed() {
		t.Error("return of state is not as expected")
	}
}

func TestTask_ErrorMessage(t *testing.T) {
	task := Task{Trace: "exit code 0", WorkerID: "testing1", ID: 1, State: "RUNNING"}
	if task.ErrorMessage() != "" {
		t.Error("return of error message is not as expected")
	}
	task = Task{Trace: "exit code 255", WorkerID: "testing2", ID: 2, State: "FAILED"}
	if task.ErrorMessage() != "exit code 255" {
		t.Error("return of error message is not as expected")
	}
}

func TestConnector_IsFailed(t *testing.T) {
	c := Connector{WorkerID: "testing1", Status: "RUNNING"}
	if c.IsFailed() {
		t.Error("return of state is not as expected")
	}
	c = Connector{WorkerID: "testing1", Status: "FAILED"}
	if !c.IsFailed() {
		t.Error("return of state is not as expected")
	}
}