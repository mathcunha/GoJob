package gojob

import (
	"time"
)

type Worker struct {
	Ops, Executed, SuccessOp, ErrorOp int
	begin, end                        int
	task                              Task
}

type Task interface {
	Run() error
}

func (w *Worker) Work() {
	w.Executed = 0
	w.begin = time.Now().Nanosecond()
	for ; w.Executed < w.Ops; w.Executed++ {
		if err := w.task.Run(); err != nil {
			w.ErrorOp++
		} else {
			w.SuccessOp++
		}
	}
	w.end = time.Now().Nanosecond()
}

func (w *Worker) IsWorking() bool {
	return w.Executed == w.Ops
}
