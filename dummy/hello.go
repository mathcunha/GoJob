package dummy

import (
	"fmt"
	"github.com/mathcunha/wreckit/gojob"
	"time"
)

type Dummy struct {
	num int
}

func (d *Dummy) Run() error {
	fmt.Printf("hello %v\n", d.num)
	return nil
}

func (d *Dummy) NewTask(prop map[string]string) gojob.Task {
	return &Dummy{num: time.Now().Nanosecond()}
}
