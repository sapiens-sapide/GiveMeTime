package computers

import (
	"fmt"
	reg "github.com/sapiens-sapide/GiveMeTime/watch/clock/registers"
	"github.com/soniakeys/meeus/coord"
	"github.com/soniakeys/meeus/deltat"
	"github.com/soniakeys/meeus/globe"
	"github.com/soniakeys/meeus/julian"
	"github.com/soniakeys/meeus/rise"
	"github.com/soniakeys/meeus/sidereal"
	"github.com/soniakeys/meeus/solar"
	"github.com/soniakeys/unit"
	"time"
	"os"
	//pp "github.com/soniakeys/meeus/planetposition"
	"github.com/soniakeys/meeus/moonphase"
	"github.com/exploded/riseset"
	"github.com/soniakeys/meeus/moonposition"
	"github.com/soniakeys/sexagesimal"
	"github.com/soniakeys/meeus/base"
)

type nextDay struct{}

var NextDay = new(nextDay)

func (nd nextDay) Trigger() {
	ny, nm, nday, _ := NextDate(int(reg.Year.Status().(uint16)), int(reg.Month.Status().(uint8)), int(reg.Day.Status().(uint16)))
	reg.SetDate(ny, nm, nday)
	DayComputation(time.Now())
}

func (nd nextDay) Status() interface{} { return nil }

func (nd nextDay) Set(value interface{}) {}

// returns next logical date (day, weekday, month, year) following the given date (ie the "tomorrow")
func NextDate(y, m, d int) (ny, nm, nd, nwd int) {
	next := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1)
	ny = next.Year()
	nm = int(next.Month())
	nd = next.Day()
	wd := int(next.Weekday())
	if wd == 0 {
		nwd = 7
	} else {
		nwd = wd
	}
	return
}

func DayComputation(now time.Time) {
	os.Setenv("VSOP87", "/usr/local/goland/src/github.com/sapiens-sapide/GiveMeTime/watch-synchronizer/astro")
	/*earth, err := pp.LoadPlanet(pp.Earth)
	if err != nil {
		fmt.Println(err)
		return
	}*/
	reg.Position = globe.Coord{
		Lat: unit.NewAngle('+', 48, 51, 39),
		Lon: unit.NewAngle('-', 2, 22, 1), //positive westward
	}

	// t0 is time at beginning of the day for astronomical calculation.
	// By convention, we set it at 3:00 in the morning to improve accuracy at middle of the day.
	t0 := time.Date(now.Year(), now.Month(), now.Day(), 3, 0, 0, 0, now.Location())
	jd0 := julian.TimeToJD(t0)

	// sun computations
	α := make([]unit.RA, 3)
	δ := make([]unit.Angle, 3)
	α[0], δ[0] = solar.ApparentEquatorial(jd0 - 1)
	α[1], δ[1] = solar.ApparentEquatorial(jd0)
	α[2], δ[2] = solar.ApparentEquatorial(jd0 + 1)

	Th0 := sidereal.Apparent0UT(jd0 - 0.125)
	ΔT := deltat.PolyAfter2000(float64(now.Year()))
	h0 := rise.Stdh0Solar
	tRise, tTransit, tSet, err := rise.Times(reg.Position, ΔT, h0, Th0, α, δ)
	if err != nil {
		fmt.Println(err)
		return
	}

	reg.SunRiseTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(time.Duration(uint64(tRise)) * time.Second)
	reg.SunZenithTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(time.Duration(uint64(tTransit)) * time.Second)
	reg.SunSetTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(time.Duration(uint64(tSet)) * time.Second)

	sunrise_jd := julian.TimeToJD(reg.SunRiseTime)
	sunrise_az, _ := coord.EqToHz(α[1], δ[1], reg.Position.Lat, reg.Position.Lon, sidereal.Apparent(sunrise_jd))
	reg.SunRiseAz = unit.AngleFromDeg(sunrise_az.Deg() + 180)

	sunzenith_jd := julian.TimeToJD(reg.SunZenithTime)
	_, reg.SunZenithAlt = coord.EqToHz(α[1], δ[1], reg.Position.Lat, reg.Position.Lon, sidereal.Apparent(sunzenith_jd))

	sunset_jd := julian.TimeToJD(reg.SunSetTime)
	sunset_az, _ := coord.EqToHz(α[1], δ[1], reg.Position.Lat, reg.Position.Lon, sidereal.Apparent(sunset_jd))
	reg.SunSetAz = unit.AngleFromDeg(sunset_az.Deg() + 180)
	reg.SunlightDuration = reg.SunSetTime.Sub(reg.SunRiseTime)

	// moon computations
	fmt.Println(julian.JDToCalendar(moonphase.Full(base.JDEToJulianYear(jd0-0.125)))) //plus proche pleine lune
	fmt.Println(julian.JDToCalendar(moonphase.New(base.JDEToJulianYear(jd0-0.125)))) //plus proche nouvelle lune

	got := riseset.Riseset(riseset.Moon, now, -reg.Position.Lon.Deg(), reg.Position.Lat.Deg(), 2)
	fmt.Printf("%+v\n", got)

/*
	//geocentric position of the moon
	jdm := julian.TimeToJD(time.Date(2017, time.Month(9), 12, 0,0,0,0,time.Local))
	λ, β, Δ := moonposition.Position(jdm) // (λ without nutation)
	Δψ, Δε := nutation.Nutation(jdm)
	fmt.Println(Δψ, Δε )
	obl := coord.NewObliquity(nutation.MeanObliquityLaskar(jdm))
	αm, δm := coord.EclToEq(λ, β, obl.S, obl.C)
	E := globe.Earth76
	ρsφʹ, ρcφʹ := E.ParallaxConstants(reg.Position.Lat, E.RadiusAtLatitude(reg.Position.Lat)*1000)
	αm, δm = parallax.Topocentric(αm, δm, Δ/base.AU, ρsφʹ, ρcφʹ, reg.Position.Lon, jdm)
	fmt.Printf("%s %s\n", sexa.FmtTime(αm.Time()), sexa.FmtAngle(δm))
	fmt.Println(Δ)
*/
	α = make([]unit.RA, 3)
	δ = make([]unit.Angle, 3)
	Δ := make([]float64, 3)
	α[0], δ[0], Δ[0] = TopoMoonPosition(jd0 - 1, reg.Position)
	α[1], δ[1], Δ[1] = TopoMoonPosition(jd0, reg.Position)
	α[2], δ[2], Δ[2] = TopoMoonPosition(jd0 + 1, reg.Position)
	Th0 = sidereal.Apparent0UT(jd0)
	h0 = rise.Stdh0Lunar(moonposition.Parallax(Δ[1]))

	tRise, tTransit, tSet, err = rise.Times(reg.Position, ΔT, h0, Th0, α, δ)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("rising:  %+.5f %02s\n", tRise/86400, sexa.FmtTime(tRise))
	fmt.Printf("transit: %+.5f %02s\n", tTransit/86400, sexa.FmtTime(tTransit))
	fmt.Printf("seting:  %+.5f %02s\n", tSet/86400, sexa.FmtTime(tSet))
/*
	obl := nutation.MeanObliquityLaskar(jd0-0.125)
	αm, δm := coord.EclToEq(λ, β, math.Sin(obl.Deg()), math.Cos(obl.Deg()))
	//coordinate of the sun
	αs, δs , R := solar.ApparentEquatorialVSOP87(earth, jd0-0.125)
	i := moonillum.PhaseAngleEq(αm, δm , Δ, αs, δs , R )
	fmt.Println(i)
	k := base.Illuminated(i)
	fmt.Printf("k = %.4f\n", k) // % illuminé
*/

}
