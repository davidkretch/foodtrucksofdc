const OK = ""
const NO_DATA = "Sorry, we don't have any data!"
const ERROR = "Oops, something went wrong!"

function status(stops) {
    if (stops.length === 0) {
        return NO_DATA;
    }
    return OK;
}

function statusError() {
    return ERROR;
}

export {
    status,
    statusError
}