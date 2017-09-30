// returns an svg representing hours as markers around a circle
// every other hour is written down as a number

function sunDisplayOn() {
    svgEl.appendChild(cc);
    svgEl.appendChild(ste[0]);
    svgEl.appendChild(ste[1]);
    sunOn = true;
}

function sunDisplayOff() {
    svgEl.removeChild(cc);
    svgEl.removeChild(ste[0]);
    svgEl.removeChild(ste[1]);
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

// returns svg elements to write sunrise & sunset to
function getSunTimesElems() {
    const txt = document.createElementNS("http://www.w3.org/2000/svg", "text");
    setAttributes(txt, {
        x: - radius * 0.45,
        y: radius * 0.05,
        style: `
        text-anchor: center;
        font-family: ${fontFamily};
        font-size: ${2 * unitbase};
        fill: ${darkcolor}`
    });
    const background = document.createElementNS("http://www.w3.org/2000/svg", "rect");
    setAttributes(background, {
        x: 0,
        y: 0,
        width: unitbase * 5.4,
        height: unitbase * 0.9,
        fill: "#FFFFFF",
        "fill-opacity": 0.7,
    });
    return [background, txt];
}