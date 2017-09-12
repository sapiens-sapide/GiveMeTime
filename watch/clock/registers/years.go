package registers

import (
	"sync"
)

type Years struct {
	lock sync.Mutex
	value uint16
}

func (y *Years) Set(value interface{}) {
	y.lock.Lock()
	y.value = uint16(value.(int))
	y.lock.Unlock()
}

func (y Years) Status() interface{} {
	return y.value
}

// borrowed from time std pkg
func (y Years) IsLeapYear() bool {
	year := int(y.value)
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}