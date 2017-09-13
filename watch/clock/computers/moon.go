package computers

import (
	"github.com/sapiens-sapide/GiveMeTime/watch-synchronizer/astro"
	r "github.com/sapiens-sapide/GiveMeTime/watch/clock/registers"
	"github.com/soniakeys/meeus/coord"
	"github.com/soniakeys/meeus/julian"
	"github.com/soniakeys/meeus/sidereal"
	"github.com/soniakeys/unit"
)

type moonPosition struct{}

var MoonPosition = new(moonPosition)

func (mp moonPosition) Trigger() {
	now := r.Now()
	jd := julian.TimeToJD(now)
	moon := astro.NewMoon(now)
	moon.SetPositions(now)

	moon_az, _ := coord.EqToHz(moon.EquatorialPosition().RA, moon.EquatorialPosition().Dec, astro.Position.Lat, astro.Position.Lon, sidereal.Apparent(jd))
	r.MoonAz = unit.AngleFromDeg(moon_az.Deg() + 180)
}

func (mp moonPosition) Status() interface{} { return nil }

func (mp moonPosition) Set(value interface{}) {}