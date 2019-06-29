# Food Trucks of DC

[Food Trucks of DC](https://foodtrucksofdc.com) is a daily food truck tracker
for Washington DC.

## Frontend

The frontend uses Firebase to authenticate users and fetch data. It uses React
to draw the page based on the fetched data.

### Firebase credentials

The frontend requires a Firebase config in `frontend/src/FirebaseConfig.js`.

### Firebase credential API requirements

Firebase uses an API key which must have access to the following APIs:

* Cloud Firestore API
* Identity Toolkit API
* Token Service API

## Backend

The backend:
1. Checks the [DC food truck lottery results](https://dcra.dc.gov/mrv).
2. Downloads any new PDFs to Cloud Storage.
3. Converts them to CSV.
4. Processes them and uploads them to Firestore.
5. Makes daily data available through Firestore.
