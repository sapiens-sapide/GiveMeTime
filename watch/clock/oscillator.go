// package clock dispatches a signal every second to its registered subscribers.
// For now, it simply relies on time.Ticker.
// In future, it should get its second ticks from a specialized & accurate device.
//
// Create an oscillator with clock.NewOscillator().
// Add subscribers, then start/stop the oscillator.

package clock

import (
	"container/list"
)

type Oscillator interface {
	Start() error
	Stop() error
	Subscribe(chan (bool)) *list.Element
	Unsubscribe(*list.Element)
}

// for now, oscillator relies on GOlang's time.Ticker
func NewOscillator() (o Oscillator, err error) {
	o = &goTicker{
		subscribers: list.New(),
	}

	return o, nil
}
