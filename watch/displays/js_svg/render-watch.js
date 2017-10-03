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
const wc = getWatchCase();
const mc = getMarkersCircle();
/** sun components **/
const cc = getCompassCircle();
const sc = getSunCircle();
let sre = getSunRiseElems(0, "");
let sse = getSunSetElems(0, "");
let sne = getSunNoonElems("");

/** buttons components **/
const b1 = getButton1();
const b2 = getButton2();
/** sync logic **/
let ephemDaysLeft = -1; // how many days ahead are in local storage. -1 means no data for current day.

svgEl.appendChild(wc);
svgEl.appendChild(civilNightLength);
svgEl.appendChild(nightLength);
timeDisplayOn();
svgEl.appendChild(mc);
svgEl.appendChild(nm);

/** buttons **/
const svgButton = document.querySelector("#watchbuttons");
svgButton.appendChild(b1);
svgButton.appendChild(b2);

b1.addEventListener("click", function () {
    if (timeOn) {
        timeDisplayOff();
        sunDisplayOn();
    } else if (sunOn) {
        sunDisplayOff();
        timeDisplayOn();
    }
});

let today = {
    d: 0
};
today.d = 0;
let hour = 0;
let min = 0;
let seconds = 60;
let secInHour = 0;
let secInDay = 0;
const position = {
    lon: -2.366944,
    lat: 48.860833
};
//dayRendering(now, position);
secondRendering();
setInterval(() => {
    secondRendering()
}, 1000);

function secondRendering() {
    //today.now = new Date("2017-06-24T03:34:20");
    today.now = new Date();
    if (seconds > 59) {
        minuteRendering();
    } else {
        seconds++
    }
    const minAngle = ((secInHour + seconds) / 3600) * 360;
    const hourAngle = ((secInDay + seconds) / 86400) * 360;
    mh.setAttribute("transform", `rotate(${minAngle - 90})`);
    hh.setAttribute("transform", `rotate(${hourAngle})`);
}

function minuteRendering() {
    hour = today.now.getHours();
    min = today.now.getMinutes();
    seconds = today.now.getSeconds();
    secInHour = min * 60;
    secInDay = hour * 3600 + secInHour;
    if (ephemDaysLeft === -1 || today.d !== today.now.getDate()) {
        dayRendering(position);
    }
    tc.innerHTML = `${hour < 10 ? "0" + hour : hour} ${min < 10 ? "0" + min : min}`;
}

function dayRendering(position) {
    today.y = today.now.getFullYear();
    today.m = today.now.getMonth();
    today.d = today.now.getDate();
    today.wd = today.now.getDay();
    wdc[1].innerHTML = weekDays[today.wd];
    dc[1].innerHTML = `${today.d < 10 ? " " + today.d : today.d}`;
    getEphemeris(today.now, position, updateSunData);
}

function updateSunData(syncSucceeded) {
    if (syncSucceeded) {
        ephemDaysLeft = 0; // TODO: manage local storage when more days will be provided by sync service
        //ste[1].innerHTML = `${Math.floor(sun.rise / 3600)}:${Math.floor(sun.rise % 3600 / 60)} - ${Math.floor(sun.set / 3600)}:${Math.floor(sun.set % 3600 / 60)}`;
        if (sunOn) {
            svgEl.removeChild(sre);
            svgEl.removeChild(sse);
            svgEl.removeChild(sne);
        }
        sre = getSunRiseElems(sun.riseAz, `● ${Math.floor(sun.rise / 3600)}:${Math.floor(sun.rise % 3600 / 60)} α ${Math.floor(sun.riseAz)}°`);
        sse = getSunSetElems(sun.setAz, `${Math.floor(sun.set / 3600)}:${Math.floor(sun.set % 3600 / 60)} α ${Math.floor(sun.setAz)}° ●`);
        sne = getSunNoonElems(`${Math.floor(sun.zenith / 3600)}:${Math.floor(sun.zenith % 3600 / 60)}`);
        if (sunOn) {
            svgEl.appendChild(sre);
            svgEl.appendChild(sse);
            svgEl.appendChild(sne);
        }
        const noonAngle = (sun.zenith / 86400) * 360;
        nm.setAttribute("transform", `rotate(${noonAngle})`);
        if (timeOn) {
            svgEl.removeChild(civilNightLength);
            svgEl.removeChild(nightLength);
        }
        civilNightLength = getNightArc(sun.civilRise, sun.civilSet);
        nightLength = getNightArc(sun.rise, sun.set);
        if (timeOn) {
            svgEl.appendChild(civilNightLength);
            svgEl.appendChild(nightLength);
        }
    } else {
        ephemDaysLeft = -1;
        ste[1].innerHTML = "";
        svgEl.removeChild(nm);
    }
}