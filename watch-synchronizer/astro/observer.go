package astro

import (
	"fmt"
	"github.com/soniakeys/meeus/globe"
	"github.com/soniakeys/meeus/julian"
	pp "github.com/soniakeys/meeus/planetposition"
	"github.com/soniakeys/unit"
	"os"
	"time"
)

type Observer struct {
	Date       time.Time
	JulianDate float64
	Position   *globe.Coord
	earth      *pp.V87Planet
	TOff       int // timezone offset in seconds
}

// date is a year, a month and a day, time is ignored but not timezone.
// lat is positive northward
// lon is positive westward
// lat & lon are in decimal degrees
func NewObserver(date time.Time, lat, lon float64) *Observer {
	_, offset := date.Zone()
	path := os.Getenv("VSOP87")
	if path == "" {
		path = "/usr/local/goland/src/github.com/sapiens-sapide/GiveMeTime/watch-synchronizer/astro/VSOP87"
	}
	E, err := pp.LoadPlanetPath(pp.Earth, path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &Observer{
		Date:       date,
		JulianDate: julian.TimeToJD(date),
		Position: &globe.Coord{
			Lat: unit.AngleFromDeg(lat),
			Lon: unit.AngleFromDeg(lon), // positive westward
		},
		earth: E,
		TOff:  offset,
	}
}
