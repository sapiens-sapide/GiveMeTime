package relays

type Hours struct {
	SteppingIntRelay
}

func NewHoursRelay(relays []Relay) *Hours {
	return &Hours{
		SteppingIntRelay{
			steps: 24,
			next:  relays,
		},
	}
}
