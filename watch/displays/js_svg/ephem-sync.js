/*
    Logic to retrieve ephemeris from network
    and store data to local storage
 */

// all times are given in minutes from midnight
const sun = {
    rise: 0,
    riseAz: 0,
    civilRise: 0,
    set: 0,
    setAz: 0,
    civilSet: 0,
    zenith: 0
};

// date is a js date object
// position is a simple obj with 2 properties :
//  lat=latitude in decimal degrees (positive northward) (Paris = 48.860833)
//  lon=longitude in decimal degrees (positive westward) (Paris = -2.366944)
function getEphemeris(date, position) {
    const xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function () {
        if (this.readyState === 4) {
            if (this.status === 200) {
                resp = JSON.parse(this.response);
                sun.rise = resp.Sun.Rise;
                sun.civilRise = resp.Sun.CivilRise;
                sun.zenith = resp.Sun.Zenith;
                sun.set = resp.Sun.Set;
                sun.civilSet = resp.Sun.CivilSet;

                wdc[1].innerHTML = weekDays[today.wd];
                dc[1].innerHTML = `${today.d < 10 ? "0" + today.d : today.d}`;
                se[1].innerHTML = `${Math.floor(sun.rise / 60)}:${Math.floor(sun.rise % 60)} - ${Math.floor(sun.set / 60)}:${Math.floor(sun.set % 60)}`;
                const noonAngle = (sun.zenith / 1440) * 360;
                nm.setAttribute("transform", `rotate(${noonAngle})`);

                svgEl.removeChild(civilNightLength);
                civilNightLength = getNightArc(sun.civilRise, sun.civilSet);
                svgEl.appendChild(civilNightLength);
                svgEl.removeChild(nightLength);
                nightLength = getNightArc(sun.rise, sun.set);
                svgEl.appendChild(nightLength);
            } else {
                console.log(this.response);
                sun.rise = 0;
                sun.civilRise = 0;
                sun.zenith = 0;
                sun.set = 0;
                sun.civilSet = 0;
            }
        }
    };
    const host = document.location.hostname;
    const enc_date = encodeURIComponent(toRFC3339string(date));
    xhttp.open("GET", `http://${host}/ephemeris?lat=${position.lat}&lon=${position.lon}&date=${enc_date}`, true);
    xhttp.send();
}