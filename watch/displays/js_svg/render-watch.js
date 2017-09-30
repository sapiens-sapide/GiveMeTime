const svgEl = document.querySelector("#watchcase");
svgEl.appendChild(getWatchCase());
let civilNightLength = document.createElement("div");
svgEl.appendChild(civilNightLength);
let nightLength = document.createElement("div");
svgEl.appendChild(nightLength);
const hh = getHourHandle();
svgEl.appendChild(hh);
const mh = getMinutesHandle();
svgEl.appendChild(mh);
const secElems = getSecElems();
svgEl.appendChild(secElems);
const tc = getTimeContainer();
svgEl.appendChild(tc);

const nm = getNoonMark();
svgEl.appendChild(nm);
//svgEl.appendChild(getOuterRect());
//svgEl.appendChild(getCrossLines());
svgEl.appendChild(getHoursCircle());
const wdc = getWeekDayContainer();
svgEl.appendChild(wdc[0]); // weekday background
svgEl.appendChild(wdc[1]); // weekday container
const dc = getDateContainer();
svgEl.appendChild(dc[0]); // date background
svgEl.appendChild(dc[1]); // date container

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
    //const t = new Date("2017-09-27T07:04:00");
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
    const noonAngle = (sun.zenith / 1440) * 360;
    nm.setAttribute("transform", `rotate(${noonAngle})`);
    svgEl.removeChild(civilNightLength);
    civilNightLength = getNightArc(sun.civilRise, sun.civilSet);
    svgEl.appendChild(civilNightLength);
    svgEl.removeChild(nightLength);
    nightLength = getNightArc(sun.rise, sun.set);
    svgEl.appendChild(nightLength);
}