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

// getData returns a promise with the truck schedule for a given date.
function getStops(date) {
    return firebase.firestore()
    .collection("schedules")
    .doc(date)
    .get()
    .then(doc => {
        if (doc.exists) {
            const data = doc.data();
            const stops = Object.keys(data).map(key => {
                return {
                    "name": key,
                    "abbrev": shortName(key),
                    "link": linkName(key),
                    "trucks": data[key]
                }
            })
            return stops;
        }
        return [];
    })
}

// getTrucks returns a promise with data on each truck (e.g. rating).
function getTrucks() {
    return firebase.firestore()
    .collection("trucks")
    .get()
    .then(query => {
        if (!query.empty) {
            var trucks = {};
            query.forEach(doc => {
                trucks[doc.id] = doc.data();
            });
            return trucks;
        }
        return {};
    })
}

// getData returns all stops for a given day, with data about each truck
// (e.g. average rating) merged onto it.
function getData(date) {
    var s = getStops(date);
    var t = getTrucks();
    return Promise.all([s, t]).then(([stops, trucks]) => {
        for (var [i, stop] of stops.entries()) {
            stop.trucks = Object.keys(stop.trucks).map(id => {
                var data = {
                    name: trucks[id]['displayName'],
                    ...trucks[id]
                };
                data.avgRating = data.avgRating || null;
                return data;
            })
            stop[i] = stop;
        }
        return stops;
    });
}

// setRating sets a user's rating for a specific truck. If the user has
// already rated the truck before, their rating will be updated.
function setRating(truck, rating) {
    const userId = firebase.auth().currentUser.uid;
    return firebase.firestore()
    .collection("ratings")
    .doc(truck)
    .collection("ratings")
    .doc(userId)
    .set({
        rating: rating,
        datetime: firebase.firestore.FieldValue.serverTimestamp()
    });
}

export {
    firebase,
    getData,
    setRating
};