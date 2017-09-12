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
	reg.Sun = relays.NewSunStepRelay([]relays.Relay{computers.SunPosition})
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
	reg.SetDate(now.Year(), int(now.Month()), now.Day())
	reg.SetTime(now.Hour(), now.Minute(), now.Second())
	//reg.SetTime(20, 25, 0)

	computers.DayComputation(now)
	computers.SunPosition.Trigger()
	osc, _ := clock.NewOscillator()
	tic := make(chan (bool))
	osc.Subscribe(tic)
	osc.Start()

	go func() {
		for range tic {
			web_svg.Disp.Reset()
			web_svg.WriteSVG(&web_svg.Disp)
			reg.Second.Trigger()
			reg.Sun.Trigger()
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


