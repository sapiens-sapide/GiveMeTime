package astro

import (
	"testing"
	"time"
)

func TestSunBody_ComputeMoment(t *testing.T) {
	T := time.Now()
	//T := time.Date(2017, time.Month(10), 3, 13, 39, 30, 0, time.Local)
	o := NewObserver(T, 48.860833, -2.366944)
	sun := NewSun(o)
	t.Log(sun.ApparentPosition())
	Tmidnight := time.Date(T.Year(),T.Month(), T.Day(), 0, 0, 0, 0, time.Local)
	day_o := NewObserver(Tmidnight, 48.860833, -2.366944)
	sun_day := NewSun(day_o)
	rise, transit, set, _ := sun_day.ComputeTransit(SunStdAlt)
	t.Log(Tmidnight.Add(time.Duration(uint64(rise)) * time.Second))
	t.Log(Tmidnight.Add(time.Duration(uint64(transit)) * time.Second))
	t.Log(Tmidnight.Add(time.Duration(uint64(set)) * time.Second))

	o_rise := NewObserver(Tmidnight.Add(time.Duration(uint64(rise)) * time.Second), 48.860833, -2.366944)
	sun_rise := NewSun(o_rise)
	t.Log(sun_rise.ApparentPosition())
}