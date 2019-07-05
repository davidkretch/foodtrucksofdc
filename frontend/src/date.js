// dateDisplay returns a Weekday, Month Day representation of a date.
// e.g. Sun Jun 30 2019 -> "Sunday, June 30"
function dateDisplay(date) {
    const date_display_options = {"weekday": "long", "month": "long", "day": "numeric"};
    return new Date(date).toLocaleDateString("en-US", date_display_options);
}

// dateKey returns a yyyy-mm-dd representation of a date.
// e.g. Sun Jun 30 2019 ... -> "2019-06-30"
function dateKey(date) {
    const d = new Date(date);
    const year = d.getFullYear();
    var month = d.getMonth() + 1;
    month = month < 10 ? "0" + month : month;
    var day = d.getDate();
    day = day < 10 ? "0" + day : day;
    return `${year}-${month}-${day}`;
}

// dateSequence returns a sequence of dates starting with the given date.
function dateSequence(date) {
    var result = new Array(7);
    for (var i = 0; i < result.length; i++) {
        var d = new Date(date);
        d.setDate(d.getDate() + i);
        result[i] = d;
    }
    return result;
}

export {
    dateDisplay,
    dateKey,
    dateSequence
};