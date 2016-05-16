package dummy

import (
	"fmt"
)

type Dummy struct {
}

func (d *Dummy) Run() error {
	fmt.Println("hello")
}
