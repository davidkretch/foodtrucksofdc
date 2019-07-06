const functions = require('firebase-functions');
const db = require('./db');

// setAvgRating updates a truck's average rating whenever 
// a user enters a new rating.
module.exports.setAvgRating = functions.firestore
    .document('ratings/{truck}/ratings/{userId}')
    .onWrite((change, context) => {

        var rating = change.after.data().rating;
        var truckRef = db.collection('trucks').doc(context.params.truck);
        
        return db.runTransaction(transaction => {
            return transaction.get(truckRef).then(truckDoc => {
                var newNumRatings;
                var newAvgRating;
                if (truckDoc.exists) {
                    var truck = truckDoc.data();
                    var oldRatingTotal = truck.avgRating * truck.numRatings;
                    newNumRatings = truck.numRatings + 1;
                    newAvgRating = (oldRatingTotal + rating) / newNumRatings;
                } else {
                    newNumRatings = 1;
                    newAvgRating = rating;
                }
                return transaction.set(truckRef, {
                    avgRating: newAvgRating,
                    numRatings: newNumRatings
                }, {merge: true});
            });
        });
    });
