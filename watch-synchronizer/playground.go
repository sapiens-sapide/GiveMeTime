package main

import (
	"fmt"
	"github.com/sapiens-sapide/GiveMeTime/watch-synchronizer/astro"
	"github.com/soniakeys/meeus/coord"
	"github.com/soniakeys/meeus/eqtime"
	"github.com/soniakeys/meeus/globe"
	"github.com/soniakeys/meeus/julian"
	pp "github.com/soniakeys/meeus/planetposition"
	"github.com/soniakeys/meeus/sidereal"
	"github.com/soniakeys/meeus/solar"
	"github.com/soniakeys/unit"
	"os"
	"time"
)

func main() {
	os.Setenv("VSOP87", "./astro")
	p := globe.Coord{
		Lat: unit.NewAngle('+', 48, 51, 39),
		Lon: unit.NewAngle('-', 2, 22, 1), //positive westward
	}
	now := time.Now()//.Add(-16*time.Hour).Add(-48*time.Minute)
	_, tz := now.Zone()
	fmt.Printf("time zone : %d\n", tz)
	jd := julian.TimeToJD(now)
	//t0 is time at beginning of the day for astronomical calculation. By convention, it is 3 in the morning.
	t0 := time.Date(now.Year(), now.Month(), now.Day(), 3, 0, 0, 0, now.Location())
	jd0 := julian.TimeToJD(t0)

	earth, err := pp.LoadPlanet(pp.Earth)
	if err != nil {
		fmt.Println(err)
		return
	}
	//computation for the day
	eot0 := eqtime.E(jd0, earth)
	α0, δ0 := solar.TrueEquatorial(jd0)
	SunRiseHA := astro.SunRiseHA(p.Lat.Deg(), δ0.Deg())
	SolarNoon := astro.SolarNoon(t0, -p.Lon.Deg(), eot0.Min(), tz/3600)

	SunRiseTime := astro.TimeFromDay(astro.DayFromTime(SolarNoon) - SunRiseHA*4/1440)
	fmt.Printf("Sunrise : %s\n", SunRiseTime)
	sunrise_jd := julian.TimeToJD(SunRiseTime)
	RiseAz, RiseAlt := coord.EqToHz(α0, δ0, p.Lat, p.Lon, sidereal.Apparent(sunrise_jd))
	fmt.Printf("Sunrise azimuth : %f\n", 180+RiseAz.Deg())
	fmt.Printf("Sunrise Alt : %f\n", RiseAlt.Deg())

	fmt.Printf("Zenith local time : %s\n", SolarNoon)
	sunzenith_jd := julian.TimeToJD(SolarNoon)
	ZenithAz, ZenithAlt := coord.EqToHz(α0, δ0, p.Lat, p.Lon, sidereal.Apparent(sunzenith_jd))
	fmt.Printf("Sun zenith azimuth : %f\n", 180+ZenithAz.Deg())
	fmt.Printf("Sun zenith Alt : %f\n", ZenithAlt.Deg())

	SunSetTime := astro.TimeFromDay(astro.DayFromTime(SolarNoon) + SunRiseHA*4/1440)
	fmt.Printf("Sunset : %s\n", SunSetTime)
	sunset_jd := julian.TimeToJD(SunSetTime)
	SetAz, SetAlt := coord.EqToHz(α0, δ0, p.Lat, p.Lon, sidereal.Apparent(sunset_jd))
	fmt.Printf("Sunset azimuth : %f\n", 180+SetAz.Deg())
	fmt.Printf("Sunset Alt : %f\n", SetAlt.Deg())

	SunlightDuration := time.Duration(int64(SunRiseHA*8*60)) * time.Second
	fmt.Printf("Sunlight duration : %s\n", SunlightDuration.String())

	//computation for current time
	α, δ := solar.TrueEquatorial(jd)
	Az, Alt := coord.EqToHz(α, δ, p.Lat, p.Lon, sidereal.Apparent(jd))
	fmt.Printf("Azimuth : %f\n", 180+Az.Deg())
	fmt.Printf("Alt : %f\n", Alt.Deg())
	tst := astro.TrueSolarTime(now, p)
	fmt.Printf("True solar time : %s\n", tst)
	sunHA := astro.SolarHourAngle(tst)
	fmt.Printf("Sun Hour Angle : %f\n", sunHA)


}
