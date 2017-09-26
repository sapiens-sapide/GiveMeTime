// helper to set attributes from an object.
function setAttributes(elem, attr) {
    Object.entries(attr).forEach(([key, value]) => {
        elem.setAttribute(key, value);
    })
}

// returns coordinates within our referential to draw an arc of `percent` of a circle
// coordinates for 0 is at top middle of our viewport
// stroke-width is taken account to shift the arc from viewport border.
function coordinatesForPercent(percent, width) {
    let x = Math.cos(2 * Math.PI * percent);
    let y = Math.sin(2 * Math.PI * percent);
    // adjust coordinates to take account of the stroke-width
    x *= 1 - ((width / 2) / radius);
    y *= 1 - ((width / 2) / radius);
    return [y, -x]; // rotate -90°
}

//returns numbers to put to stroke-dasharray property
function gaugeParam(radius, percent) {
    const circ = 2 * Math.PI * radius;
    const gaugeLength = circ * percent;
    return [gaugeLength, circ];
}

// output a RFC3339 string from a js date object
// with timezone.
// example output : 2006-01-02T15:04:05+02:00
function toRFC3339string(date) {
    const tz = -date.getTimezoneOffset(); // offset to add to local time to get UTC time in minutes (for ex., Paris' offset is -120 at DST)
    const abs_offset = Math.abs(tz);
    const tz_h = Math.floor(abs_offset / 60) < 10 ? "0" + Math.floor(abs_offset / 60) : Math.floor(abs_offset / 60);
    const tz_m = abs_offset % 60;
    const date_shifted = new Date(Date.parse(date) + (tz * 60000));
    const tz_str = date_shifted.toISOString().substr(0, 19); // remove trailing 'Z' and nano seconds.
    return tz_str + (tz > 0 ? "+" : "-") + tz_h + ":" + (tz_m < 10 ? "0" + tz_m : tz_m);
}