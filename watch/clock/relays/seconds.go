package relays

type Seconds struct {
	SteppingIntRelay
}

func NewSecondsRelay(relays []Relay) *Seconds {
	return &Seconds{
		SteppingIntRelay{
			steps: 60,
			next: relays,
		},
	}
}