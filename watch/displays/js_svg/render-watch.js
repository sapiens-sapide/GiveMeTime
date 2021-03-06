const svgEl = document.querySelector("#watchcase");

/** time components **/
let civilNightLength = document.createElement("div");
let nightLength = document.createElement("div");
let timeOn = false;
let sunOn = false;
let MoonOn = false;
const hh = getHourHandle();
const mh = getMinutesHandle();
const secElems = getSecElems();
const tc = getTimeContainer();
const nm = getNoonMark();
const wdc = getWeekDayContainer();
const dc = getDateContainer();
const hc = getHoursCircle();
const wb = getWatchBorder();
const wback = getWatchBackground();
const mc = getMarkersCircle();
let mooncomp = getMoon(0);

/** sun components **/
const cc = getCompassCircle();
const sc = getSunCircle();
let sre = getSunRiseElems(0, "");
let sse = getSunSetElems(0, "");
let sne = getSunNoonElems("");

/** buttons components **/
const bSun = getButtonSun();
const bMoon = getButtonMoon();
const bSet = getButtonSet();
const txtPosition = getTxtPosition();
/** sync logic **/
let ephemDaysLeft = -1; // how many days ahead are in local storage. -1 means no data for current day.
let lastUpdate = new Date();

svgEl.appendChild(wback);
svgEl.appendChild(wb);
timeDisplayOn();

/** buttons **/
const svgButton = document.querySelector("#watchbuttons");
svgButton.appendChild(bSun);
svgButton.appendChild(bMoon);
svgButton.appendChild(bSet);
svgButton.appendChild(txtPosition);

bSun.addEventListener("click", function () {
    if (timeOn) {
        timeDisplayOff();
        sunDisplayOn();
    } else if (sunOn) {
        sunDisplayOff();
        timeDisplayOn();
    }
});

bSet.addEventListener("click", function () {
    getLocation()
});

let today = {
    d: 0
};
today.d = -1;
let hour = -1;
let min = -1;
let seconds = 60;
let secInHour = 0;
let secInDay = 0;
const position = {
    lon: 2.0,
    lat: 48.0,
    update: 0,
};
//dayRendering(now, position);
secondRendering();
setInterval(() => {
    secondRendering()
}, 1000);

function secondRendering() {
    //today.now = new Date("2018-03-02T00:00:00");
    today.now = new Date();
    if ((today.now.getSeconds() - seconds > 1) || (hour !== today.now.getHours()) || (min !== today.now.getMinutes())) {
        minuteRendering();
    } else {
        seconds = today.now.getSeconds()
    }
    const minAngle = ((secInHour + seconds) / 3600) * 360;
    const hourAngle = ((secInDay + seconds) / 86400) * 360;
    mh.setAttribute("transform", `rotate(${minAngle - 90})`);
    hh.setAttribute("transform", `rotate(${hourAngle - 90})`);
}

function minuteRendering() {
    hour = today.now.getHours();
    min = today.now.getMinutes();
    seconds = today.now.getSeconds();
    secInHour = min * 60;
    secInDay = hour * 3600 + secInHour;
    // update every 6 hours for moon accuracy.
    // TODO: get isMoonEvent boundaries for current day to prevent every 6 hours update
    if (ephemDaysLeft === -1 ||
        today.d !== today.now.getDate() ||
        ((today.now - lastUpdate) / 3600000) > 6) {
        dayRendering(position);
    }
    tc.innerHTML = `${hour < 10 ? "0" + hour : hour}   ${min < 10 ? "0" + min : min}`;

    for (let i = 1; i < 24; i += 2) {
        let h = document.getElementById("hour" + i);
        h.setAttribute("fill-opacity", 0);
        let m = document.getElementById("odd-mark-" + i);
        m.setAttribute("opacity", 0);
    }

    if (hour % 2 !== 0) {
        let h = document.getElementById("hour" + hour);
        h.setAttribute("fill-opacity", ((60 - min) / 60));
        let m = document.getElementById("odd-mark-" + hour);
        m.setAttribute("opacity", ((60 - min) / 60));
    } else {
        let h = document.getElementById("hour" + (hour + 1));
        h.setAttribute("fill-opacity", (min / 60));
        let m = document.getElementById("odd-mark-" + (hour + 1));
        m.setAttribute("opacity", (min / 60));
    }

}

function dayRendering(position) {
    today.y = today.now.getFullYear();
    today.m = today.now.getMonth();
    today.d = today.now.getDate();
    today.wd = today.now.getDay();
    wdc[1].innerHTML = weekDays[today.wd];
    dc[1].innerHTML = `${today.d < 10 ? " " + today.d : today.d}`;
    txtPosition.innerHTML=`lat: ${position.lat.toPrecision(5)}° lon: ${position.lon.toPrecision(5)}°`;
    getEphemeris(today.now, position, updateSunData);
}

function updateSunData(syncSucceeded) {
    if (sunOn) {
        sunDisplayOff();
    }
    if (timeOn) {
        timeDisplayOff();
    }
    if (syncSucceeded) {
        const rise_minutes = Math.floor(sun.rise % 3600 / 60);
        sre = getSunRiseElems(sun.riseAz, `● ${Math.floor(sun.rise / 3600)}:${rise_minutes < 10 ? "0" + rise_minutes : rise_minutes} α ${Math.floor(sun.riseAz)}°`);
        const set_minutes = Math.floor(sun.set % 3600 / 60);
        sse = getSunSetElems(sun.setAz, `${Math.floor(sun.set / 3600)}:${set_minutes < 10 ? "0" + set_minutes : set_minutes} α ${Math.floor(sun.setAz)}° ●`);
        const noon_minutes = Math.floor(Math.floor(sun.zenith % 3600 / 60));
        sne = getSunNoonElems(`${Math.floor(sun.zenith / 3600)}:${noon_minutes < 10 ? "0" + noon_minutes : noon_minutes}`);
        const noonAngle = (sun.zenith / 86400) * 360;
        nm.setAttribute("transform", `rotate(${noonAngle})`);
        civilNightLength = getNightArc(sun.civilRise, sun.civilSet);
        nightLength = getNightArc(sun.rise, sun.set);
        ephemDaysLeft = 0; // TODO: manage local storage when more days will be provided by sync service
        lastUpdate = new Date();
        mooncomp = getMoon(moon.isMoonEvent)
    } else {
        ephemDaysLeft = -1;
    }
    timeDisplayOn()
}

function getLocation() {
    if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(updatePosition);
    } else {
        x.innerHTML = "Geolocation is not supported by this browser.";
    }
}

function updatePosition(data) {
    position.lat = data.coords.latitude;
    position.lon = data.coords.longitude;
    position.update = new Date();
    dayRendering(position);
}