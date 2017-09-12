package astro

import (
	"github.com/soniakeys/meeus/globe"
	pp "github.com/soniakeys/meeus/planetposition"
	"github.com/soniakeys/unit"
)

const (
	Rise = iota
	Transit
	Set
	Yesterday = iota - 1
	Today
	Tomorrow
)

/*
TODO : refactor these var affectations
*/
var Position *globe.Coord
var Sun *SunBody
var Moon *MoonBody
var Earth *pp.V87Planet

type EclipticPosition struct {
	long, lat unit.Angle
	distance  float64 // distance between centers of the Earth and body, in AU
	jd        float64 // julian date at which computation has been done
}

type EquatorialPosition struct {
	RA  unit.RA
	Dec unit.Angle
	jd  float64 // julian date at which computation has been done
}

type CelestialBody interface {
	EclipticPosition() EclipticPosition
	EquatorialPosition() EquatorialPosition
}

type TransitEvent struct {
	Type int8 // rise=0, transit=1, set=2
	Day  int8 // -1 for yesterday, 0 for today, 1 for tomorrow
	Time unit.Time
	Az   unit.Angle
	Elev unit.Angle
}
