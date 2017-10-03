package ephemeris

import (
	"testing"
	"time"
)

func TestEphemerisForDay(t *testing.T) {
	t.Log(EphemerisForDay(time.Now(), 48.860833, -2.366944))
}
