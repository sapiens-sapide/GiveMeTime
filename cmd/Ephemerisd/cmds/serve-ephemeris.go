package cmds

import (
	"encoding/json"
	"github.com/sapiens-sapide/GiveMeTime/watch-synchronizer/ephemeris"
	"github.com/soniakeys/meeus/globe"
	"github.com/soniakeys/unit"
	"log"
	"net/http"
	"time"
	"strconv"
)

func StartEphemerisServer() {

	http.HandleFunc("/ephemeris", returnEphemeris)
	log.Fatal(http.ListenAndServe(":1971", nil))

}

// returns computed ephemeris for given date & position
// query params :
//  lat=latitude in decimal degrees (positive northward) (Paris = 48.860833)
//  lon=longitude in decimal degrees (positive westward) (Paris = -2.366944)
//  date=RFC3339 date string
func returnEphemeris(w http.ResponseWriter, req *http.Request) {
	var lat float64
	var lon float64
	var date time.Time
	var err error
	q := req.URL.Query()

	if lt, ok := q["lat"]; !ok {
		http.Error(w, "'lat' param is missing", http.StatusBadRequest)
		return
	} else {
		lat, err = strconv.ParseFloat(lt[0], 32)
		if err != nil {
			http.Error(w, "unable to parse latitude : " + err.Error(), http.StatusBadRequest)
		}
	}

	if ln, ok := q["lon"]; !ok {
		http.Error(w, "'lon' param is missing", http.StatusBadRequest)
		return
	} else {
		lon, err = strconv.ParseFloat(ln[0], 32)
		if err != nil {
			http.Error(w, "unable to parse longitude : " + err.Error(), http.StatusBadRequest)
		}
	}

	if d, ok := q["date"]; !ok {
		http.Error(w, "'date' param is missing", http.StatusBadRequest)
		return
	} else {
		date, err = time.Parse(time.RFC3339, d[0])
		if err != nil {
			http.Error(w, "unable to parse date string : " + err.Error(), http.StatusBadRequest)
			return
		}
	}

	eph, err := ephemeris.EphemerisForDay(date, globe.Coord{
		Lat: unit.AngleFromDeg(lat),
		Lon: unit.AngleFromDeg(lon), //positive westward
	})
	if err != nil {
		log.Println(err)
	}
	j_eph, err := json.Marshal(eph)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(j_eph)
}

func returnError(w http.ResponseWriter, err error) {

}
