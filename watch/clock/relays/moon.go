package relays

type MoonStep struct {
	SteppingIntRelay
}

func NewMoonStepRelay(relays []Relay) *MoonStep {
	return &MoonStep{
		SteppingIntRelay{
			steps: 20,
			next:  relays,
		},
	}
}
