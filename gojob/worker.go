package gojob

import (
	"time"
)

type Worker struct {
	Ops, Executed, SuccessOp, ErrorOp int
	Begin, End                        int
	task                              Task
}

type Task interface {
	Run() error
	LoadProperties(map[string]string)
	Clone() Task
}

func (w *Worker) Work() {
	w.Executed = 0
	w.Begin = time.Now().Nanosecond()
	for ; w.Executed < w.Ops; w.Executed++ {
		if err := w.task.Run(); err != nil {
			w.ErrorOp++
		} else {
			w.SuccessOp++
		}
	}
	w.End = time.Now().Nanosecond()
}

func (w *Worker) IsWorking() bool {
	return w.Executed == w.Ops
}
