// Go-pendulum implements Oscillator interface with package time.Ticker

package clock

import (
	"container/list"
	"time"
)

type goTicker struct {
	pendulum    *time.Ticker
	subscribers *list.List
}

// add a subscriber to the oscillator. It will receive "true" every second
// subscribers follow FILO pattern
func (t *goTicker) Subscribe(c chan (bool)) *list.Element {
	return t.subscribers.PushFront(c)
}

func (t *goTicker) Unsubscribe(e *list.Element) {
	t.subscribers.Remove(e)
}

func (t *goTicker) Start() error {
	t.pendulum = time.NewTicker(time.Second)
	go func(t *goTicker) {
		for range t.pendulum.C {
			for e := t.subscribers.Front(); e != nil; e = e.Next() {
				// select clause is for non blocking writing into subscribers' channels,
				// ie. : if a channel write blocks, oscillator must ignore it and continue to drain pendulum
				select {
				case e.Value.(chan (bool)) <- true:
				default:
					continue
				}
			}
		}
	}(t)
	return nil
}

func (t *goTicker) Stop() error {
	t.pendulum.Stop()
	return nil
}
