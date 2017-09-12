package registers

import "sync"

type Months struct {
	lock    sync.Mutex
	value   uint8 // range from 1 to 12
	strings map[string][]string
}

func NewMonths() *Months {
	return &Months{
		strings: map[string][]string{
			"fr-short": {"Jan", "Fév", "Mars", "Avr", "Mai", "Juin", "Juil", "Août", "Sept", "Oct", "Nov", "Déc"},
		},
	}
}

func (m Months) Display(lang string, index uint8) string {
	if s, ok := m.strings[lang]; ok {
		if int(index-1) < len(s) {
			return s[index-1]
		} else {
			return "not available"
		}
	} else {
		return "not available"
	}
}

// month #1 is january, month #12 is december
func (m *Months) Status() interface{} {
	return m.value
}

// month #1 is january, month #12 is december
func (m *Months) Set(value interface{}) {
	if value.(int) > 0 && value.(int) < 13 {
		m.lock.Lock()
		m.value = uint8(value.(int))
		m.lock.Unlock()
	}
}
