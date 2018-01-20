package ephemeris

import (
	"testing"
	"time"
)

func TestEphemerisForDay(t *testing.T) {
	eph, err := EphemerisForDay(time.Now(), 48.860833, -2.366944)
	t.Logf("%+v, %+v", eph, err)
}
