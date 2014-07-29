// bad implement. but you can pass the test.
package calc

type Calc struct {
	vals []int64
	ans  int64
}

func NewCalc() *Calc {
	c := &Calc{}
	return c
}

func (c *Calc) Push(val int64) {
	c.vals = append(c.vals, val)
}

func (c *Calc) Add() {
	for _, v := range c.vals {
		c.ans += v
	}
}

func (c *Calc) Result() int64 {
	return c.ans
}
