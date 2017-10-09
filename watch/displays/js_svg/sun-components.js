// returns an svg representing hours as markers around a circle
// every other hour is written down as a number

function sunDisplayOn() {
    svgEl.appendChild(cc);
    svgEl.appendChild(sre);
    svgEl.appendChild(sse);
    svgEl.appendChild(sne);
    svgEl.appendChild(sc);
    svgEl.appendChild(mc);
    sunOn = true;
}

function sunDisplayOff() {
    svgEl.removeChild(cc);
    svgEl.removeChild(sre);
    svgEl.removeChild(sse);
    svgEl.removeChild(sne);
    svgEl.removeChild(sc);
    svgEl.removeChild(mc);
    sunOn = false;
}

function getCompassCircle() {
    //group for all components
    const hc = document.createElementNS("http://www.w3.org/2000/svg", "g");
    //hc.appendChild(getMarkersCircle());

    //group for hours name, slightly shifted to the bottom
    const hg = document.createElementNS("http://www.w3.org/2000/svg", "g");
    setAttributes(hg, {
        transform: `translate(0, ${unitbase * 0.3})`
    });
    const cardinals = ["N", "E", "S", "W"];
    for (let i = 0; i < 12; i++) {
        const txtG = document.createElementNS("http://www.w3.org/2000/svg", "g");
        const background = document.createElementNS("http://www.w3.org/2000/svg", "circle");
        const ang = 30 * i;
        const coord = coordinatesForPercent(ang / 360, 48 * dot);
        setAttributes(background, {
            cx: coord[0] * 0.95,
            cy: coord[1] * 0.95,
            r: unitbase * 0.35,
            fill: "#FFFFFF",
            transform: "translate(0, -0.022)"
        });
        const txt = document.createElementNS("http://www.w3.org/2000/svg", "text");

        if (i % 3 === 0) {
            setAttributes(txt, {
                "class": "hoursNum",
                x: coord[0] * 0.95,
                y: coord[1] * 0.95,
                "font-size": unitbase,
                style: `text-anchor: middle;font-family: ${fontFamily}`
            });
            txt.innerHTML = cardinals[i / 3]
        } else {
            setAttributes(txt, {
                "class": "hoursNum",
                x: coord[0] * 0.95,
                y: coord[1] * 0.95,
                "font-size": unitbase * 0.55,
                style: `text-anchor: middle;font-family: ${smallFontFamily}`
            });
            txt.innerHTML = i * 30
        }

        txtG.appendChild(background);
        txtG.appendChild(txt);
        hg.appendChild(txtG);
    }
    hc.appendChild(hg);
    return hc;
}

// returns svg elements with sun noon time
function getSunNoonElems(string) {
    const txt = document.createElementNS("http://www.w3.org/2000/svg", "text");
    setAttributes(txt, {
        x: 0,
        y: radius * 0.29 + (unitbase * 3),
        "font-size": unitbase * 0.7,
        style: `text-anchor: middle;font-family: ${smallFontFamily};`
    });
    txt.innerHTML = string;
    return txt;
}

// returns svg elements with sunrise marker and text
// positioned on sunCircle
function getSunRiseElems(sunriseAz, string) {
    const g = document.createElementNS("http://www.w3.org/2000/svg", "g");
    setAttributes(g, {
        transform: `translate (0 ${radius * 0.29})`
    });
    const coord = coordinatesForPercent(sunriseAz / 360, dot);
    const txt = document.createElementNS("http://www.w3.org/2000/svg", "text");
    setAttributes(txt, {
        x: coord[0] * unitbase * 2.3,
        y: coord[1] * unitbase * 2.3,
        "font-size": unitbase * 0.7,
        style: `text-anchor: start;font-family: ${smallFontFamily};`
    });
    txt.innerHTML = string;
    g.appendChild(txt);

    return g;
}

// returns svg elements with sunrise marker and text
// positioned on sunCircle
function getSunSetElems(sunsetAz, string) {
    const g = document.createElementNS("http://www.w3.org/2000/svg", "g");
    setAttributes(g, {
        transform: `translate (0 ${radius * 0.29})`
    });
    const coord = coordinatesForPercent(sunsetAz / 360, dot);
    const txt = document.createElementNS("http://www.w3.org/2000/svg", "text");
    setAttributes(txt, {
        x: coord[0] * unitbase * 2.3,
        y: coord[1] * unitbase * 2.3,
        "font-size": unitbase * 0.7,
        style: `text-anchor: end;font-family: ${smallFontFamily};`
    });
    txt.innerHTML = string;
    g.appendChild(txt);

    return g;
}

function getSunCircle() {
    const sc = document.createElementNS("http://www.w3.org/2000/svg", "g");
    setAttributes(sc, {
        transform: `translate (0 ${radius * 0.27})`
    });
    for (let i = 0; i < 4; i++) {
        const g = document.createElementNS("http://www.w3.org/2000/svg", "g");
        const ang = (90 * i) - 90;
        setAttributes(g, {
            transform: `rotate(${ang})`
        });
        const r = document.createElementNS("http://www.w3.org/2000/svg", "rect");
        const w = unitbase * 0.5;
        const x = (unitbase * 2.5) - w;
        setAttributes(r, {
            x: x - dot,
            y: -dot / 2,
            width: w,
            height: dot,
            fill: darkcolor,
            stroke: darkcolor,
            "stroke-width": dot
        });
        g.appendChild(r);
        sc.appendChild(g);
    }
    const c = document.createElementNS("http://www.w3.org/2000/svg", "circle");
    setAttributes(c, {
        cx: 0,
        cy: 0,
        r: unitbase * 2.5,
        fill: "none",
        stroke: darkcolor,
        "stroke-width": dot,
    });
    sc.appendChild(c);
    return sc;
}