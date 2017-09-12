package relays

type SunStep struct {
	SteppingIntRelay
}

func NewSunStepRelay(relays []Relay) *SunStep {
	return &SunStep{
		SteppingIntRelay{
			steps: 20,
			next: relays,
		},
	}
}