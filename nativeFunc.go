package mygolox

import "time"

type clock struct {
}

func NewClock() *clock {
	return &clock{}
}

func (c *clock) Arity() int {
	return 0
}

func (c *clock) Call(interpreter Interpreter, arguments []any) any {
	return float64(time.Now().Unix())
}

func (c *clock) String() string {
	return "<native fn>"
}
