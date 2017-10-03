package astro

import (
	"github.com/soniakeys/unit"
)

const (
	Rise = iota
	Transit
	Set
	Yesterday = iota - 1
	Today
	Tomorrow
	SecInDay    = 86400.0
	MinInDay    = 1440.0
	SunCivilAlt = -350 // min of angle to compute civil sunrise/set
	SunStdAlt   = -50  // min of angle to compute sunrise/set
)

type EclipticPosition struct {
	long, lat unit.Angle
	distance  float64 // distance between centers of the Earth and body, in AU
}

type EquatorialPosition struct {
	RA  unit.RA
	Dec unit.Angle
}

//	Az: azimuth of observed point, measured westward from the South.
//	Elev: elevation, or height of observed point above horizon.
type ApparentPosition struct {
	Az, Elev float64
}

type CelestialBody interface {
	EclipticPosition() EclipticPosition
	EquatorialPosition() EquatorialPosition
	ApparentPosition() ApparentPosition
	JD() float64 //returns julian date at which celestial body has been set
}

type TransitEvent struct {
	Type int8 // rise=0, transit=1, set=2
	Day  int8 // -1 for yesterday, 0 for today, 1 for tomorrow
	Time unit.Time
	Az   unit.Angle
	Elev unit.Angle
}
