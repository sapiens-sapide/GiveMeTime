package astro

import (
	"fmt"
	"github.com/soniakeys/meeus/eqtime"
	"github.com/soniakeys/meeus/globe"
	"github.com/soniakeys/meeus/julian"
	pp "github.com/soniakeys/meeus/planetposition"
	"math"
	"os"
	"time"
)

// returns solar hour angle at sunrise
// lat is observer latitude (positive northward)
// sun_dec is declination of the sun for the day
// algorithm from NOAA spreadsheet (https://www.esrl.noaa.gov/gmd/grad/solcalc/NOAA_Solar_Calculations_day.xls)
func SunRiseHA(lat, sun_dec float64) float64 {
	return Degrees(math.Acos(
		math.Cos(Radians(90.833))/
			(math.Cos(Radians(lat))*math.Cos(Radians(sun_dec))) -
			math.Tan(Radians(lat))*math.Tan(Radians(sun_dec))))
}

func SolarNoon(day time.Time, long, eot float64, tz int) time.Time {
	noon := (720 - 4*long - eot + float64(tz)*60) / 1440
	second_shift := noon * 24 * 60 * 60
	return time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location()).Add(time.Duration(second_shift) * time.Second)
}

func TrueSolarTime(t time.Time, pos globe.Coord) time.Time {
	os.Setenv("VSOP87", "./astro")
	jd := julian.TimeToJD(t)
	earth, err := pp.LoadPlanet(pp.Earth)
	if err != nil {
		fmt.Println(err)
		return time.Time{}
	}
	eq := eqtime.E(jd, earth)
	_, offset := t.Zone()
	tst := math.Mod(float64(t.Hour())*60.0+float64(t.Minute())+(float64(t.Second())/60.0)+eq.Min()+(4*(-pos.Lon.Deg())-60.0*float64(offset)/3600.0), 1440)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC).Add(time.Duration(tst*60) * time.Second)
}

// tst = True Solar Time
func SolarHourAngle(tst time.Time) float64 {
	tst_minutes := float64(tst.Hour()*60+tst.Minute()) + float64(tst.Second())/60.0
	t := tst_minutes / 4.0
	if t < 0 {
		return t + 180
	}
	return t - 180
}
