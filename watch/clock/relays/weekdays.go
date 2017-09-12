package relays

type Weekdays struct {
	SteppingIntRelay
	strings map[string][]string // map holds languages strings for Weekdayss. Example {"fr-short":["lun", "mar", "mer", "jeu", "ven", "sam", "dim"]
								// first weekday is monday = 1, sunday = 7.
}

func NewWeekdaysRelay(relays []Relay) *Weekdays {
	return &Weekdays{
		SteppingIntRelay: SteppingIntRelay{
			steps: 7,
			next: relays,
		},
		strings: map[string][]string{
			"fr-short":{"lun", "mar", "mer", "jeu", "ven", "sam", "dim"},
		},
	}
}

func (wd Weekdays) Display(lang string, index uint8) string {
	if s, ok := wd.strings[lang]; ok {
		if int(index-1) < len(s) {
			return s[index-1]
		} else {
			return "not available"
		}
	} else {
		return "not available"
	}
}