const svgEl = document.querySelector("#watchcase");
svgEl.appendChild(getWatchCase());
let civilNightLength = document.createElement("div");
svgEl.appendChild(civilNightLength);
let nightLength = document.createElement("div");
svgEl.appendChild(nightLength);
const tc = getTimeContainer();
svgEl.appendChild(tc);
const hh = getHourHandle();
svgEl.appendChild(hh);
const mh = getMinutesHandle();
svgEl.appendChild(mh);
const nm = getNoonMark();
svgEl.appendChild(nm);
const se = getSunTimesElems();
svgEl.appendChild(se[0]); // sun times background
svgEl.appendChild(se[1]); // sun times container
//svgEl.appendChild(getOuterRect());
//svgEl.appendChild(getCrossLines());
svgEl.appendChild(getHoursCircle());

const wdc = getWeekDayContainer();
svgEl.appendChild(wdc[0]); // weekday background
svgEl.appendChild(wdc[1]); // weekday container
const dc = getDateContainer();
svgEl.appendChild(dc[0]); // date background
svgEl.appendChild(dc[1]); // date container
const secElems = getSecElems();
svgEl.appendChild(secElems[0]);
const sc = getSecContainer();
svgEl.appendChild(sc);
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
    //const t = new Date("2017-09-26T10:50:00");
    const t = new Date();
    const d = t.getDate();
    const h = t.getHours();
    const min = t.getMinutes();
    seconds = t.getSeconds();
    tc.innerHTML = `${h < 10 ? "0" + h : h}  ${min < 10 ? "0" + min : min}`;
    secInHour = min * 60;
    secInDay = h * 3600 + secInHour;
    sc.innerHTML = `${seconds < 10 ? "0" + seconds : seconds}`;
    const dasharray = gaugeParam(secondsRadius, seconds / 60);
    secElems[1].setAttribute("stroke-dasharray", dasharray[0] + " " + dasharray[1]);
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
    getEphemeris(date, position);
}

