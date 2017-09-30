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
const ste = getSunTimesElems();
/** buttons components **/
const b1 = getButton1();
const b2 = getButton2();

svgEl.appendChild(wc);
svgEl.appendChild(civilNightLength);
svgEl.appendChild(nightLength);
timeDisplayOn();
svgEl.appendChild(mc);

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

const now = new Date();
let today = {
    y: now.getYear(),
    m: now.getMonth(),
    d: now.getDate(),
    wd: now.getDay()
};
let seconds = 60;
let secInHour = 0;
let secInDay = 0;
const position = {
    lon: -2.366944,
    lat: 48.860833
};
dayRendering(now, position);
secondRendering();
setInterval(() => {
    secondRendering()
}, 1000);

function secondRendering() {
    //const t = new Date("2017-09-26T13:34:20");
    const t = new Date();
    const d = t.getDate();
    const h = t.getHours();
    const min = t.getMinutes();
    seconds = t.getSeconds();
    tc.innerHTML = `${h < 10 ? "0" + h : h} ${min < 10 ? "0" + min : min}`;
    secInHour = min * 60;
    secInDay = h * 3600 + secInHour;
    const minAngle = ((secInHour + seconds) / 3600) * 360;
    const hourAngle = ((secInDay + seconds) / 86400) * 360;
    mh.setAttribute("transform", `rotate(${minAngle - 90})`);
    hh.setAttribute("transform", `rotate(${hourAngle})`);
    if (d !== today.d) {
        dayRendering(t, position);
    }
}

function dayRendering(date, position) {
    today.y = date.getYear();
    today.m = date.getMonth();
    today.d = date.getDate();
    today.wd = date.getDay();
    getEphemeris(date, position, updateSunData);
}

function updateSunData() {
    wdc[1].innerHTML = weekDays[today.wd];
    dc[1].innerHTML = `${today.d < 10 ? "0" + today.d : today.d}`;
    ste[1].innerHTML = `${Math.floor(sun.rise / 60)}:${Math.floor(sun.rise % 60)} - ${Math.floor(sun.set / 60)}:${Math.floor(sun.set % 60)}`;
    const noonAngle = (sun.zenith / 1440) * 360;
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
}