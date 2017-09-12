package relays

import "sync"

type Relay interface {
	Set(interface{})
	Status() interface{} // func to call to get current relay's value
	Trigger()            // func to call to send a signal to the relay
}

// a stepping relay receives a signal from its predecessor, increments its steps counter until limit,
// then its reset its counter and forward signal to its successor(s)
type SteppingIntRelay struct {
	lock  sync.Mutex
	floor uint8 // count starts from floor
	next  []Relay
	steps uint8 // how many steps to count before cascading ingress signal to next relay(s)
	value uint8
}

func (sr *SteppingIntRelay) Trigger() {
	sr.lock.Lock()
	defer sr.lock.Unlock()
	if sr.value < sr.steps-(sr.floor+1) {
		sr.value++
	} else {
		sr.value = sr.floor
		for _, relay := range sr.next {
			relay.Trigger()
		}
	}
}

func (sr *SteppingIntRelay) Status() interface{} {
	return sr.value
}

func (sr *SteppingIntRelay) Set(value interface{}) {
	sr.lock.Lock()
	sr.value = uint8(value.(int))
	sr.lock.Unlock()
}
