package web_svg

import (
	"fmt"
	"github.com/ajstarks/svgo"
	reg "github.com/sapiens-sapide/GiveMeTime/watch/clock/registers"
	"github.com/soniakeys/unit"
	"io"
	"math"
	"time"
)

const (
	S             = 320 //320px
	fontSize      = 16  //default body font size in px. = 1 em
	dayLengthDiff = "-210"
	nextMoon      = "20.09"
	waxing_moon   = false
)

func WriteSVG(fd io.Writer) {
	width := S
	height := S
	canvas := svg.New(fd)
	canvas.Start(width, height)
	//canvas.Roundrect(0, 0, width, width, 20, 20, "fill:none;stroke:black;stroke-width:1")
	//canvas.Rect(0, 0, S, 390, "fill:white;stroke:black;stroke-width:1") //large applewatch screen size

	//watch outline
	canvas.Circle(S/2, S/2, S/2, "fill:#C4C4C4")
	canvas.Circle(S/2, S/2, int(math.Floor(S/2-1.2*fontSize)), "fill:white")
	var circleStyle string

	/*
	//corners
	// - weekday
	canvas.Text(int(math.Floor(S*0.027)), int(math.Floor(S*0.09)),
		reg.Weekday.Display("fr-short", reg.Weekday.Status().(uint8)),
		"font-family:courier;font-weight:bold;font-size:2em;fill:black")
	// - day
	canvas.Text(int(math.Floor(S*0.82)), int(math.Floor(S*0.113)),
		fmt.Sprintf("%02d", reg.Day.Status()),
		"font-family:courier;font-weight:bold;font-size:3em;letter-spacing: -.05em;fill:black")
	// - month
	canvas.Text(int(math.Floor(S*0.027)), int(math.Floor(S*0.93)),
		reg.Month.Display("fr-short", reg.Month.Status().(uint8)),
		"font-family:courier;font-size:1.2em;fill:black")
	// - annual gauge
	canvas.TranslateRotate(int(math.Floor(float64(width)*0.915)), int(math.Floor(float64(height)*0.915)), -90)
	yearPercent := float64(reg.YearDay.Status().(uint16)) / float64(reg.YearLength.Status().(uint16))
	circleStyle := fmt.Sprintf("stroke-width:0.25em;stroke-dasharray:%s;stroke:black;fill:white", gaugeParam(1.25*fontSize, yearPercent))
	canvas.Circle(0, 0, 1.25*fontSize, circleStyle)
	canvas.Gend()
	canvas.Text(int(math.Floor(float64(width)*0.915)), int(math.Floor(float64(height)*0.93)),
		fmt.Sprintf("%d", reg.Year.Status()),
		"text-anchor:middle;font-family:courier;font-size:0.75em;fill:black")
*/


	//time
	hoursAngle := float64(int(reg.Hour.Status().(uint8))*3600+int(reg.Minute.Status().(uint8))*60+int(reg.Second.Status().(uint8))) / 240.0
	canvas.TranslateRotate(S/2, S/2, -90.0+hoursAngle)
	canvas.Translate(0, -int(math.Floor(0.23*fontSize)))
	canvas.Rect(int(math.Floor(float64(width)*0.42)), 0, int(math.Floor(float64(width)*0.075)), int(math.Floor(0.5*fontSize)), "fill:red")
	canvas.Gend()
	canvas.Gend()

	minutesAngle := float64(int(reg.Minute.Status().(uint8))*60+int(reg.Second.Status().(uint8))) / 10.0
	canvas.TranslateRotate(S/2, S/2, -90.0+minutesAngle)
	canvas.Translate(0, -int(math.Floor(0.06*fontSize)))
	canvas.Rect(int(math.Floor(float64(width)*0.03)), 0, int(math.Floor(float64(width)*0.47)), int(math.Floor(0.125*fontSize)), "fill:blue")
	canvas.Gend()
	canvas.Gend()

	timeStr := fmt.Sprintf("%02d %02d", reg.Hour.Status(), reg.Minute.Status())
	canvas.Text(int(math.Floor(float64(width)*0.5)), int(math.Floor(float64(height)*0.58)), timeStr, "text-anchor:middle;font-family:Courier;font-weight:bold;font-size:6.5em;letter-spacing: -0.1em;fill:black")

	secondStr := fmt.Sprintf("%02d", reg.Second.Status())
	canvas.TranslateRotate(int(math.Floor(float64(width)*0.5)), int(math.Floor(float64(width)*0.53)), -90)
	circleStyle = fmt.Sprintf("stroke-width:0.125em;stroke-dasharray:%s;stroke:black;fill:white", gaugeParam(0.875*fontSize, float64(reg.Second.Status().(uint8))/60.0))
	canvas.Circle(0, 0, 0.875*fontSize, circleStyle)
	canvas.Gend()
	canvas.Text(int(math.Floor(float64(width)*0.5)), int(math.Floor(float64(height)*0.55)), secondStr, "text-anchor:middle;font-family:courier;font-size:1.25em;fill:black")
	
	// - weekday
	canvas.Rect(int(math.Floor(float64(width)*0.33)), int(math.Floor(float64(height)*0.2)), int(math.Floor(float64(width)*0.18)), int(math.Floor(float64(height)*0.06)), "fill:white;opacity:0.9")
	canvas.Text(int(math.Floor(S*0.33)), int(math.Floor(S*0.25)),
		reg.Weekday.Display("fr-short", reg.Weekday.Status().(uint8)),
		"font-family:courier;font-weight:bold;font-size:2em;fill:black")
	// - day
	canvas.Rect(int(math.Floor(float64(width)*0.54)), int(math.Floor(float64(height)*0.17)), int(math.Floor(float64(width)*0.16)), int(math.Floor(float64(height)*0.11)), "fill:white;opacity:0.9")
	canvas.Text(int(math.Floor(S*0.53)), int(math.Floor(S*0.27)),
		fmt.Sprintf("%02d", reg.Day.Status()),
		"font-family:courier;font-weight:bold;font-size:3em;letter-spacing: -.05em;fill:black")

	//moon's azimuth
	if reg.MoonRise.Before(reg.MoonSet) {
		if reg.Now().After(reg.MoonRise) && reg.Now().Before(reg.MoonSet) {
			canvas.TranslateRotate(S/2, S/2, -90+reg.MoonAz.Deg())
			canvas.Circle((S/2)-10, 0, 10, "fill:black;opacity:0.5")
			canvas.Gend()
		}
	} else {
		if reg.Now().Before(reg.MoonSet) || reg.Now().After(reg.MoonRise) {
			canvas.TranslateRotate(S/2, S/2, -90+reg.MoonAz.Deg())
			canvas.Circle((S/2)-10, 0, 10, "fill:black;opacity:0.5")
			canvas.Gend()
		}
	}


	//day length arc
	riseHourAngle := float64(reg.SunRiseTime.Hour()*3600+reg.SunRiseTime.Minute()*60) / 240.0
	canvas.TranslateRotate(int(math.Floor(float64(width)*0.5)), int(math.Floor(float64(width)*0.5)), -90+riseHourAngle)
	dayLengthPercent := float64(reg.SunlightDuration/time.Second) / 86400.0
	circleStyle = fmt.Sprintf("stroke-width:0.15em;stroke-dasharray:%s;stroke:red;fill:none", gaugeParam(float64(width)*0.497, dayLengthPercent))
	canvas.Circle(0, 0, int(float64(width)*0.497), circleStyle)
	canvas.Gend()

	//noon mark
	noonHourAngle := float64(reg.SunZenithTime.Hour()*3600+reg.SunZenithTime.Minute()*60) / 240.0
	canvas.TranslateRotate(int(math.Floor(float64(width)*0.5)), int(math.Floor(float64(width)*0.5)), -90+noonHourAngle)
	canvas.Rect(int(math.Floor(S/2-1.3*fontSize)), 0, int(math.Floor(0.6*fontSize)), 0.125*fontSize, "stroke: yellow; fill:yellow")
	canvas.Gend()

	//sunrise & sunset
	drawSunMarker(canvas, reg.SunRiseAz)
	drawSunMarker(canvas, reg.SunSetAz)

	//sun's azimuth
	if reg.Now().After(reg.SunRiseTime) && reg.Now().Before(reg.SunSetTime) {
		canvas.TranslateRotate(S/2, S/2, -90+reg.SunAz.Deg())
		canvas.Circle((S/2)-10, 0, 10, "fill:yellow")
		canvas.Gend()
		canvas.Circle(S/2, int(math.Floor(S/1.35)), int(math.Floor(0.8*fontSize)), "fill:yellow")
		sun_az_str := fmt.Sprintf("%03.0f", reg.SunAz.Deg())
		canvas.Text(S/2, int(math.Floor(S/1.315)), sun_az_str, "text-anchor:middle;font-family:courier;font-size:0.875em;fill:black")
	}

	//moon
	/*
	canvas.TranslateRotate(int(math.Floor(float64(width)*0.5)), int(math.Floor(float64(height)*0.2)), -90.0)
	if waxing_moon {
		circleStyle = fmt.Sprintf("stroke-width:0.15em;stroke-dasharray:%s;stroke:black;fill:white;opacity:0.9", gaugeParam(23, reg.MoonPercent))
	} else {
		canvas.Circle(0, 0, 23, "stroke-width:0.15em;stroke:black;fill:none")
		circleStyle = fmt.Sprintf("stroke-width:0.25em;stroke-dasharray:%s;stroke:white;fill:white", gaugeParam(23, 1-reg.MoonPercent))
	}
	canvas.Circle(0, 0, 23, circleStyle)
	canvas.Gend()
	canvas.Rect(int(math.Floor(float64(width)*0.25)), int(math.Floor(float64(height)*0.233)), int(math.Floor(float64(width)*0.16)), int(math.Floor(float64(height)*0.045)), "fill:white;opacity:0.9")
	moonrise_str := fmt.Sprintf("➚%02d:%02d", reg.MoonRise.Hour(), reg.MoonRise.Minute())
	canvas.Text(int(math.Floor(float64(width)*0.25)), int(math.Floor(float64(height)*0.27)), moonrise_str, "text-anchor:left;font-family:courier;font-size:0.875em;fill:black")
	//canvas.Text(int(math.Floor(float64(width)*0.25)), int(math.Floor(float64(height)*0.27)), moonSet, "text-anchor:left;font-family:courier;font-size:0.875em;fill:black")
	if waxing_moon {
		canvas.Text(int(math.Floor(float64(width)*0.5)), int(math.Floor(float64(height)*0.17)), "●", "text-anchor:middle;font-family:courier;font-size:0.875em;fill:black")
	} else {
		canvas.Text(int(math.Floor(float64(width)*0.5)), int(math.Floor(float64(height)*0.17)), "○", "text-anchor:middle;font-family:courier;font-size:0.875em;fill:black")
	}
	canvas.Text(int(math.Floor(float64(width)*0.5)), int(math.Floor(float64(height)*0.21)), nextMoon, "text-anchor:middle;font-family:courier;font-size:0.875em;fill:black")
	canvas.Rect(int(math.Floor(float64(width)*0.6)), int(math.Floor(float64(height)*0.233)), int(math.Floor(float64(width)*0.16)), int(math.Floor(float64(height)*0.045)), "fill:white;opacity:0.9")
	//canvas.Text(int(math.Floor(float64(width)*0.6)), int(math.Floor(float64(height)*0.27)), moonRise, "text-anchor:left;font-family:courier;font-size:0.875em;fill:black")
	moonset_str := fmt.Sprintf("➘%02d:%02d", reg.MoonSet.Hour(), reg.MoonSet.Minute())
	canvas.Text(int(math.Floor(float64(width)*0.6)), int(math.Floor(float64(height)*0.27)), moonset_str, "text-anchor:left;font-family:courier;font-size:0.875em;fill:black")
*/
	//day data
	//canvas.Rect(int(math.Floor(float64(width)*0.31)), int(math.Floor(float64(height)*0.71)), int(math.Floor(float64(width)*0.4)), int(math.Floor(float64(height)*0.205)), "fill:white;opacity:0.9")

	sunrise_str := fmt.Sprintf("%02d:%02d", reg.SunRiseTime.Hour(), reg.SunRiseTime.Minute())
	sunset_str := fmt.Sprintf("%02d:%02d", reg.SunSetTime.Hour(), reg.SunSetTime.Minute())
	canvas.Rect(int(math.Floor(float64(width)*0.27)), int(math.Floor(float64(height)*0.76)), int(math.Floor(float64(width)*0.16)), int(math.Floor(float64(height)*0.045)), "fill:white;opacity:0.9")
	canvas.Rect(int(math.Floor(float64(width)*0.57)), int(math.Floor(float64(height)*0.76)), int(math.Floor(float64(width)*0.16)), int(math.Floor(float64(height)*0.045)), "fill:white;opacity:0.9")
	canvas.Text(int(math.Floor(float64(width)*0.5)), int(math.Floor(float64(height)*0.8)), sunrise_str+"     "+sunset_str, "text-anchor:middle;font-family:courier;font-size:1em;fill:black")
	/*
	sunlight_str := fmt.Sprintf("%02dh%02dm", (reg.SunlightDuration/time.Minute)/60, reg.SunlightDuration/time.Minute-((reg.SunlightDuration/time.Minute)/60)*60)
	canvas.Text(int(math.Floor(float64(width)*0.5)), int(math.Floor(float64(height)*0.83)), sunlight_str+" ["+dayLengthDiff+"s] ", "text-anchor:middle;font-family:courier;font-size:0.875em;fill:black")
	sunnoon_str := fmt.Sprintf("%02d:%02d α %02d°", reg.SunZenithTime.Hour(), reg.SunZenithTime.Minute(), int(reg.SunZenithAlt.Deg()))
	canvas.Text(int(math.Floor(float64(width)*0.5)), int(math.Floor(float64(height)*0.9)), sunnoon_str, "text-anchor:middle;font-family:courier;font-size:0.875em;fill:black")
	*/

	//cardinal markers
	canvas.TranslateRotate(S/2, S/2, -90)
	for i := 0; i < 12; i++ {
		canvas.Rotate(30)
		canvas.Rect(int(math.Floor(S/2-0.8*fontSize)), 0, int(math.Floor(0.8*fontSize)), 0.125*fontSize)
	}
	for i := 0; i < 12; i++ {
		canvas.Gend()
	}
	canvas.Gend()
	canvas.TranslateRotate(S/2, S/2, -75)
	for i := 0; i < 12; i++ {
		canvas.Rotate(30)
		canvas.Rect(int(math.Floor(S/2-0.3*fontSize)), 0, int(math.Floor(0.3*fontSize)), 0.125*fontSize)
	}
	for i := 0; i < 12; i++ {
		canvas.Gend()
	}
	canvas.Gend()

	canvas.End()
}

//returns string numbers to put to stroke-dasharray property
func gaugeParam(radius, percent float64) (strokeDasharray string) {
	circ := 2.0 * math.Pi * radius
	gaugeLength := circ * percent
	return fmt.Sprintf("%f,%f", gaugeLength, circ)
}

func drawSunMarker(canvas *svg.SVG, az unit.Angle) {
	canvas.TranslateRotate(int(math.Floor(float64(S/2))), int(math.Floor(float64(S/2))), -90.0+az.Deg())
	canvas.Line(int(math.Floor(S/2-0.95*fontSize)), 0, S/2, 0, "stroke-width:0.1875em;stroke:yellow")
	canvas.Gend()
}
