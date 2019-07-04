import firebase from "firebase/app";
import "firebase/auth";
import "firebase/firestore";
import firebaseConfig from "./firebaseConfig";

firebase.initializeApp(firebaseConfig);

// shortName returns a shortened stop name, with only the first two words.
// e.g. "Farragut Square 17th St" -> "Farragut Square"
function shortName(name) {
    return name.split(/ |\//).slice(0, 2).join(" ");
}

// linkName returns a stop name suitable for use in a URL.
// e.g. "Farragut Square 17th St" -> "farragut-square"
function linkName(name) {
    return shortName(name).replace(/[^0-9a-z ]/gi, '').toLowerCase().split(" ").join("-");
}

// getData returns a promise with the stop data for a given date.
function getData(date) {
    return firebase.firestore()
    .collection("dates")
    .doc(date)
    .get();
}

// processData returns an array of stops, each with their trucks.
// The array is built from the data returned by getData.
function processData(doc) {
    if (doc.exists) {
        const data = doc.data();
        const arr = Object.keys(data).map(key => {
            return {
                "name": key,
                "abbrev": shortName(key),
                "link": linkName(key),
                "trucks": data[key]
            }
        })
        return arr;
    }
    return [];
}

export {
    firebase,
    getData,
    processData
};