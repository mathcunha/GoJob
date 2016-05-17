package dummy

import (
	"fmt"
	"github.com/mathcunha/wreckit/gojob"
	"time"
)

type Dummy struct {
	num int
}

func (d *Dummy) LoadProperties(prop map[string]string) {
	fmt.Printf("loading properties %p\n", d)
}

func (d *Dummy) Run() error {
	fmt.Printf("hello %v\n", d.num)
	return nil
}

func (d *Dummy) Clone() gojob.Task {
	return &Dummy{num: time.Now().Nanosecond()}
}
