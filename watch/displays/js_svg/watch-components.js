/**
 Designed for a 320px (diam) circular screen.
 The viewbox is a 2unit x 2unit square, with origin (0,0) at center.
 Top left corner is at x=-1,y=-1.
 Within this referential, 1 unit = 160px.
 **/
const radius = 1;
const unitbase = 0.1; // = 1 em = 16px
const dot = 0.625e-2; // currently 1px
const darkcolor = "#000000";
const bluecolor = "#0D66FD";
const redcolor = "#FB2733";
const greycolor = "#656363";
const fontFamily = "asapregular";
const fontBold = "asapbold";
const smallFontFamily = "robotoregular";
const weekDays = ["dim", "lun", "mar", "mer", "jeu", "ven", "sam"];
const secondsRadius = radius * 0.05;

// returns the svg that represents the watchCase
// for now, it's just the outer circle.
function getWatchCase() {
    const r = radius - (dot / 2);
    const circEl = document.createElementNS("http://www.w3.org/2000/svg", "circle");
    setAttributes(circEl, {
        id: "case",
        r: r,
        stroke: darkcolor,
        "stroke-width": dot,
        fill: "#FFFFFF"
    });
    return circEl;
}

function timeDisplayOn() {
    svgEl.appendChild(hh);
    if (ephemDaysLeft !== -1) {
        svgEl.appendChild(civilNightLength);
        svgEl.appendChild(nightLength);
        svgEl.appendChild(mooncomp);
    }
    svgEl.appendChild(mh);
    svgEl.appendChild(secElems);
    svgEl.appendChild(tc);
//svgEl.appendChild(getOuterRect());
//svgEl.appendChild(getCrossLines());
    svgEl.appendChild(wdc[0]); // weekday background
    svgEl.appendChild(wdc[1]); // weekday container
    svgEl.appendChild(dc[0]); // date background
    svgEl.appendChild(dc[1]); // date container
    if (ephemDaysLeft !== -1) {
        svgEl.appendChild(nm);
    }
    svgEl.appendChild(mc);
    svgEl.appendChild(hc);
    timeOn = true;
}

function timeDisplayOff() {
    svgEl.removeChild(hh);
    svgEl.removeChild(mh);
    svgEl.removeChild(secElems);
    svgEl.removeChild(tc);
    if (ephemDaysLeft !== -1) {
        svgEl.removeChild(civilNightLength);
        svgEl.removeChild(nightLength);
        svgEl.removeChild(nm);
        svgEl.removeChild(mooncomp);
    }
//svgEl.removeChild(getOuterRect());
//svgEl.removeChild(getCrossLines());
    svgEl.removeChild(wdc[0]); // weekday background
    svgEl.removeChild(wdc[1]); // weekday container
    svgEl.removeChild(dc[0]); // date background
    svgEl.removeChild(dc[1]); // date container
    svgEl.removeChild(hc);
    svgEl.removeChild(mc);
    timeOn = false;
}

// returns the svg for an arc of `percent` of a full circle starting at the top middle of our viewport
// NB : an arc of 100% will not render correctly. Use circle instead.
// percent is a fraction of 1.
// rad is the radius to scale to as a fraction of 1.
function getArc(percent, strokeWidth, rad) {

    const start = coordinatesForPercent(0, strokeWidth);
    const end = coordinatesForPercent(percent, strokeWidth);
    const largeArcFlag = percent > .5 ? 1 : 0;
    const r = 1 - ((strokeWidth / 2) / radius);

    const pathData = [
        `M ${start[0] * rad} ${start[1] * rad}`,
        `A ${r * rad} ${r * rad} 0 ${largeArcFlag} 1 ${end[0] * rad} ${end[1] * rad}`,
    ].join(' ');

    const arcEl = document.createElementNS("http://www.w3.org/2000/svg", "path");
    setAttributes(arcEl, {
        d: pathData,
        fill: "transparent",
        stroke: greycolor,
        "stroke-opacity": 0.3,
        "stroke-width": strokeWidth
    });
    return arcEl;
}

// returns an svg arc representing the night length
// sunrise and sunset are in minutes from midnight.
function getNightArc(sunrise, sunset) {
    const nightLength = 86400 - (sunset - sunrise);
    const nightArc = getArc(nightLength / 86400, 15 * dot, radius);
    const sunsetAngle = (sunset / 86400) * 360;
    nightArc.setAttribute("transform", `rotate(${sunsetAngle})`);
    return nightArc;
}

// returns an svg representing hours as markers around a circle
// every other hour is written down as a number
function getHoursCircle() {
    //group for all components
    const hc = document.createElementNS("http://www.w3.org/2000/svg", "g");
    //hc.appendChild(getMarkersCircle());

    //group for hours name, slightly shifted to the bottom
    const hg = document.createElementNS("http://www.w3.org/2000/svg", "g");
    setAttributes(hg, {
        transform: `translate(0, ${unitbase * 0.2})`
    });
    for (let i = 0; i < 12; i++) {
        if (i % 3 !== 0) {
            //minutes
            const txtG = document.createElementNS("http://www.w3.org/2000/svg", "g");
            const background = document.createElementNS("http://www.w3.org/2000/svg", "circle");
            const ang = 30 * i;
            const coord = coordinatesForPercent(ang / 360, 52 * dot);
            setAttributes(background, {
                cx: coord[0],
                cy: coord[1],
                r: unitbase * 0.6,
                fill: bluecolor,
                stroke: "#FFFFFF",
                "stroke-width": dot,
                transform: "translate(0, -0.022)"
            });
            const txt = document.createElementNS("http://www.w3.org/2000/svg", "text");
            setAttributes(txt, {
                "class": "hoursNum",
                x: coord[0],
                y: coord[1],
                "font-size": unitbase * 0.55,
                fill: "#FFFFFF",
                style: `text-anchor: middle;font-family: ${smallFontFamily};font-weight:bold;`,
                transform: "translate(-0.001, -0.002)"
            });
            txt.innerHTML = i * 5;
            txtG.appendChild(background);
            txtG.appendChild(txt);
            hg.appendChild(txtG);
        }
    }
    for (let i = 0; i < 24; i++) {
        // hours
        const txtG2 = document.createElementNS("http://www.w3.org/2000/svg", "g");
        //const background2 = document.createElementNS("http://www.w3.org/2000/svg", "circle");
        const ang2 = 15 * i;
        const coord2 = coordinatesForPercent(ang2 / 360, 15 * dot);
        /*setAttributes(background2, {
            cx: coord2[0],
            cy: coord2[1],
            r: unitbase * 0.35,
            fill: "#FFFFFF",
            "fill-opacity": 0.7,
            transform: "translate(0, -0.022)"
        });*/
        const txt2 = document.createElementNS("http://www.w3.org/2000/svg", "text");
        setAttributes(txt2, {
            "class": "hoursNum",
            x: coord2[0],
            y: coord2[1],
            "font-size": unitbase * 0.4,
            fill: "#000000",
            style: `text-anchor: middle;font-family: ${smallFontFamily}`,
            transform: "translate(0, -0.008)"
        });
        txt2.innerHTML = i;
        //txtG2.appendChild(background2);
        txtG2.appendChild(txt2);
        hg.appendChild(txtG2);
    }
    hc.appendChild(hg);
    return hc;
}

// returns an svg representing markers distributed around a circle
// there are 12 main markers + 12 small markers.
function getMarkersCircle() {
    const mg = document.createElementNS("http://www.w3.org/2000/svg", "g");
    for (let i = 0; i < 24; i++) {

        // minutes markers
        if (i % 2 === 0) {
            const g = document.createElementNS("http://www.w3.org/2000/svg", "g");
            const ang = (15 * i) - 90;
            setAttributes(g, {
                transform: `rotate(${ang})`
            });
            if (i % 6 === 0) {
                const l = document.createElementNS("http://www.w3.org/2000/svg", "line");
                setAttributes(l, {
                    x1: 7 * unitbase,
                    y1: 0,
                    x2: radius * 0.88,
                    y2: 0,
                    stroke: greycolor,
                    "stroke-width": dot,
                    "stroke-linecap": "round"
                });
                g.appendChild(l);
            }

            mg.appendChild(g);
        }
        // hours markers
        const g = document.createElementNS("http://www.w3.org/2000/svg", "g");
        const ang = (15 * i) - 60;
        setAttributes(g, {
            transform: `rotate(${ang})`
        });
        const r = document.createElementNS("http://www.w3.org/2000/svg", "rect");
        const w = dot * 2;
        const h = dot * 1;
        const x = radius - ( dot * 2);
        setAttributes(r, {
            x: x,
            y: -h / 2,
            width: w,
            height: h,
            fill: greycolor,
            stroke: greycolor,
            "stroke-width": dot
        });
        g.appendChild(r);
        mg.appendChild(g);
    }
    return mg;
}

// returns an svg `text` components ready to render hours & minutes within.
function getTimeContainer() {
    const txt = document.createElementNS("http://www.w3.org/2000/svg", "text");
    setAttributes(txt, {
        x: 0,
        y: radius * 0.15,
        style: `
        text-anchor: middle;
        font-family: ${fontBold};
        font-size: ${5 * unitbase};
        letter-spacing: ${-0.05 * unitbase};
        fill: ${darkcolor};`
    });
    return txt;
}

// returns an svg `text` components ready to render weekday within.
function getWeekDayContainer() {
    const txt = document.createElementNS("http://www.w3.org/2000/svg", "text");
    setAttributes(txt, {
        x: -0.01 * radius,
        y: -0.431 * radius,
        style: `
        text-anchor: end;
        font-family: ${fontFamily};
        font-size: ${1.5 * unitbase};
        fill: ${darkcolor}`
    });
    const background = document.createElementNS("http://www.w3.org/2000/svg", "rect");
    setAttributes(background, {
        x: -0.345 * radius,
        y: -0.6 * radius,
        width: unitbase * 3.34,
        height: unitbase * 1.7,
        fill: "#FFFFFF",
        "fill-opacity": 0.9,
    });
    return [background, txt];
}

// returns an svg `text` components ready to render date within.
function getDateContainer() {
    const txt = document.createElementNS("http://www.w3.org/2000/svg", "text");
    setAttributes(txt, {
        x: 0.01 * radius,
        y: -0.415 * radius,
        style: `
        text-anchor: start;
        font-family: ${fontBold};
        font-size: ${2.55 * unitbase};
        fill: ${darkcolor}`
    });
    const background = document.createElementNS("http://www.w3.org/2000/svg", "rect");
    setAttributes(background, {
        x: 0.015 * radius,
        y: -0.615 * radius,
        width: unitbase * 2.8,
        height: unitbase * 2.1,
        fill: "#FFFFFF",
        "fill-opacity": 0.9,
    });
    return [background, txt];
}

// returns the svg elements to render the seconds arc,
// ie: the container and the inner circle.
function getSecElems() {
    const g = document.createElementNS("http://www.w3.org/2000/svg", "g");
    //setAttributes(g, {
    //    transform: `translate(0, ${unitbase * 0.7}) rotate(-90)`
    //});
    const c1 = document.createElementNS("http://www.w3.org/2000/svg", "circle");
    //const dasharray = gaugeParam(secondsRadius, 0); // default settings to 0
    setAttributes(c1, {
        r: secondsRadius,
        //stroke: darkcolor,
        //"stroke-width": dot * 2,
        fill: greycolor,
        //"stroke-dasharray": dasharray[0] + " " + dasharray[1]
    });
    const c2 = document.createElementNS("http://www.w3.org/2000/svg", "circle");
    //const dasharray = gaugeParam(secondsRadius, 0); // default settings to 0
    setAttributes(c2, {
        r: secondsRadius * 13.3,
        //stroke: greycolor,
        //"stroke-width": dot,
        //"stroke-opacity": 0.1,
        fill: "#FFFFFF",
        "fill-opacity": 0.5
        //"stroke-dasharray": dasharray[0] + " " + dasharray[1]
    });
    g.appendChild(c2);
    g.appendChild(c1);
    return g;
}

function getMinutesHandle() {
    const g = document.createElementNS("http://www.w3.org/2000/svg", "g");
    //rounded rectangle
    const r = document.createElementNS("http://www.w3.org/2000/svg", "rect");
    const thickness = dot * 3;
    setAttributes(r, {
        x: 0,
        y: -thickness / 2,
        width: radius - dot * 25,
        height: thickness,
        stroke: darkcolor,
        "stroke-width": dot * 1,
        fill: "none",
    });
    // arrow
    const a = document.createElementNS("http://www.w3.org/2000/svg", "path");
    const y = 0.07 * radius;
    const x2 = 0.14 * radius;
    setAttributes(a, {
        d: `m${0.85 * radius} ${y} 
        l${x2} -${y} 
        l-${x2} -${y} 
        l-${1.3 * x2} ${y} 
        z`,
        fill: bluecolor
    });
    const l = document.createElementNS("http://www.w3.org/2000/svg", "line");
    setAttributes(l, {
        y1: 0,
        x1: 7.4 * unitbase,
        y2: 0,
        x2: radius * 0.985,
        stroke: bluecolor,
        "stroke-width": dot * 2,
        "stroke-linecap": "round"
    });
    g.appendChild(r);
    g.appendChild(a);
    return g;
}

function getHourHandle() {
    const hg = document.createElementNS("http://www.w3.org/2000/svg", "g");
    const l = document.createElementNS("http://www.w3.org/2000/svg", "line");
    setAttributes(l, {
        x1: 6.7 * unitbase,
        y1: 0,
        x2: radius * 0.97,
        y2: 0,
        stroke: redcolor,
        "stroke-width": dot * 2,
        "stroke-linecap": "round"
    });
    const c = document.createElementNS("http://www.w3.org/2000/svg", "circle");
    setAttributes(c, {
        cx: radius - unitbase * 0.5,
        cy: 0,
        r: unitbase * 0.45,
        fill: redcolor,
    });
    hg.appendChild(l);
    hg.appendChild(c);
    return hg;
}

function getNoonMark() {
    /*const g = document.createElementNS("http://www.w3.org/2000/svg", "line");
    setAttributes(g, {
        x1: 0,
        y1: -0.82,
        x2: 0,
        y2: -0.995,
        stroke: greycolor,
        "stroke-width": dot * 2
    });*/

    const g = document.createElementNS("http://www.w3.org/2000/svg", "circle");

    setAttributes(g, {
        cx: 0,
        cy: -0.905,
        r: unitbase * 0.25,
        fill: greycolor,
    });

    return g;
}

function getMoon(isMoonEvent) {
    if (isMoonEvent && isMoonEvent !== 0) {
        const moon = document.createElementNS("http://www.w3.org/2000/svg", "circle");
        const fill = isMoonEvent === 1 ? "transparent" : greycolor;
        setAttributes(moon, {
            cx: 0,
            cy: 0.35,
            r: secondsRadius * 1.3,
            fill: fill,
            stroke: "#000000",
            "stroke-width": dot,
        });
        return moon;
    }
    return document.createElement("div");
}

/** below are temporary components to help design the clock **/
function getOuterRect() {
    const recEl = document.createElementNS("http://www.w3.org/2000/svg", "rect");
    setAttributes(recEl, {
        x: -1,
        y: -1,
        width: 2,
        height: 2,
        fill: "transparent",
        stroke: darkcolor,
        "stroke-width": dot
    });
    return recEl;
}

function getCrossLines() {
    const cl = document.createElementNS("http://www.w3.org/2000/svg", "g");
    cl.appendChild(getVLine());
    cl.appendChild(getHLine());
    return cl;
}

function getVLine() {
    const vLine = document.createElementNS("http://www.w3.org/2000/svg", "line");
    setAttributes(vLine, {
        x1: 0,
        y1: -1,
        x2: 0,
        y2: 1,
        stroke: darkcolor,
        "stroke-width": dot
    });
    return vLine;
}

function getHLine() {
    const hLine = document.createElementNS("http://www.w3.org/2000/svg", "line");
    setAttributes(hLine, {
        x1: -1,
        y1: 0,
        x2: 1,
        y2: 0,
        stroke: darkcolor,
        "stroke-width": dot
    });
    return hLine;
}

function getButton1() {
    const txtG = document.createElementNS("http://www.w3.org/2000/svg", "g");
    setAttributes(txtG, {
        "class": "button"
    });
    const background = document.createElementNS("http://www.w3.org/2000/svg", "circle");
    setAttributes(background, {
        cx: -radius * 1.3,
        cy: -radius * 0.3,
        r: unitbase * 4,
        fill: "#FFFFFF",
        "class": "button"
    });
    const txt = document.createElementNS("http://www.w3.org/2000/svg", "text");
    setAttributes(txt, {
        "class": "buttonTxt",
        x: -radius * 1.3,
        y: -radius * 0.2,
        "font-size": unitbase * 3,
        style: `text-anchor: middle;font-family: ${fontFamily}`,
    });
    txt.innerHTML = "Sun";
    txtG.appendChild(background);
    txtG.appendChild(txt);
    return txtG;
}

function getButton2() {
    const txtG = document.createElementNS("http://www.w3.org/2000/svg", "g");
    const background = document.createElementNS("http://www.w3.org/2000/svg", "circle");
    setAttributes(background, {
        cx: radius * 1.3,
        cy: -radius * 0.3,
        r: unitbase * 4,
        fill: "#FFFFFF",
    });
    const txt = document.createElementNS("http://www.w3.org/2000/svg", "text");
    setAttributes(txt, {
        "class": "buttonTxt",
        x: radius * 1.3,
        y: -radius * 0.2,
        "font-size": unitbase * 3,
        style: `text-anchor: middle;font-family: ${fontFamily}`
    });
    txt.innerHTML = "Moon";
    txtG.appendChild(background);
    txtG.appendChild(txt);
    return txtG;
}