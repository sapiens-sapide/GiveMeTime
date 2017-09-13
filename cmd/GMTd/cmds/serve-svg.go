package cmds

import (
	"github.com/sapiens-sapide/GiveMeTime/watch/clock"
	"github.com/sapiens-sapide/GiveMeTime/watch/clock/computers"
	reg "github.com/sapiens-sapide/GiveMeTime/watch/clock/registers"
	"github.com/sapiens-sapide/GiveMeTime/watch/clock/relays"
	"github.com/sapiens-sapide/GiveMeTime/watch/displays/web-svg"
	"io"
	"log"
	"net/http"
	"time"
)

func init() {
	reg.Month = reg.NewMonths()
	reg.Weekday = relays.NewWeekdaysRelay([]relays.Relay{})
	reg.Hour = relays.NewHoursRelay([]relays.Relay{reg.Weekday, computers.NextDay})
	reg.Minute = relays.NewMinutesRelay([]relays.Relay{reg.Hour})
	reg.Second = relays.NewSecondsRelay([]relays.Relay{reg.Minute})
	reg.SunRelay = relays.NewSunStepRelay([]relays.Relay{computers.SunPosition})
	reg.MoonRelay = relays.NewMoonStepRelay([]relays.Relay{computers.MoonPosition})
	reg.Day = new(reg.Days)
	reg.Day.Set(0)
	reg.YearDay = new(reg.Days)
	reg.YearDay.Set(0)
	reg.Year = new(reg.Years)
	reg.Year.Set(0)
	reg.YearLength = new(reg.Days)
	reg.YearLength.Set(0)
	reg.Tz = new(reg.TimeZones)
	reg.Tz.Set(0)
	reg.Dst = new(reg.DST)
	reg.Dst.Set(false)
}

func StartSVGServer() {
	now := time.Now()
	//now := time.Date(2017, 8, 21, 11, 0, 0, 0, time.UTC)
	reg.SetDate(now.Year(), int(now.Month()), now.Day())
	//reg.SetDate(2017, 8, 21)
	reg.SetTime(now.Hour(), now.Minute(), now.Second())
	//reg.SetTime(18, 0, 0)
	reg.Tz.Set(2)

	computers.DayComputation(reg.Now())
	computers.SunPosition.Trigger()
	computers.MoonPosition.Trigger()
	osc, _ := clock.NewOscillator()
	tic := make(chan (bool))
	osc.Subscribe(tic)
	osc.Start()

	go func() {
		for range tic {
			web_svg.Disp.Reset()
			web_svg.WriteSVG(&web_svg.Disp)
			reg.Second.Trigger()
			reg.SunRelay.Trigger()
			reg.MoonRelay.Trigger()
		}
	}()

	http.HandleFunc("/", IndexServer)
	http.HandleFunc("/watch-screen", SVGserver)
	log.Fatal(http.ListenAndServe(":1971", nil))

}

func IndexServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, `
	<!DOCTYPE html>
		<html lang="en">
			<head>
                <meta charset="UTF-8">
                <title>GiveMeTime</title>
			</head>
			<body style="font-size:16px">
				<img id="watch-screen" src="watch-screen.svg">
			</body>
			<script>
                setInterval(function(){
                    document.getElementById('watch-screen').src = "watch-screen?random="+new Date().getTime();
                }, 1000);
			</script>
		</html>
	`)
}

func SVGserver(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	web_svg.Disp.CopyTo(w)
}
