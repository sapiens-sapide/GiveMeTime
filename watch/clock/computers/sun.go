package computers

import (
	"github.com/sapiens-sapide/GiveMeTime/watch-synchronizer/astro"
	r "github.com/sapiens-sapide/GiveMeTime/watch/clock/registers"
	"github.com/soniakeys/meeus/coord"
	"github.com/soniakeys/meeus/julian"
	"github.com/soniakeys/meeus/sidereal"
	"github.com/soniakeys/meeus/solar"
	"github.com/soniakeys/unit"
	"time"
)

type sunPosition struct{}

var SunPosition = new(sunPosition)

func (sp sunPosition) Trigger() {
	now := time.Now()
	jd := julian.TimeToJD(now)
	α, δ := solar.ApparentEquatorial(jd)
	sun_az, _ := coord.EqToHz(α, δ, astro.Position.Lat, astro.Position.Lon, sidereal.Apparent(jd))
	r.SunAz = unit.AngleFromDeg(sun_az.Deg() + 180)
}

func (sp sunPosition) Status() interface{} { return nil }

func (sp sunPosition) Set(value interface{}) {}
