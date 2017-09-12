package relays

type Minutes struct {
	SteppingIntRelay
}

func NewMinutesRelay(relays []Relay) *Minutes {
	return &Minutes{
		SteppingIntRelay{
			steps: 60,
			next:  relays,
		},
	}
}
