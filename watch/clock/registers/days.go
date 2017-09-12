package registers

import "sync"

type Days struct {
	lock  sync.Mutex
	value uint16
}

func (d *Days) Set(value interface{}) {
	d.lock.Lock()
	d.value = uint16(value.(int))
	d.lock.Unlock()
}

func (d Days) Status() interface{} {
	return d.value
}
