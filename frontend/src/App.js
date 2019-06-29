import React from "react";

import {firebase, getData, processData} from "./Firebase";
import Header from "./Header";
import Layout from "./Layout";
import Content from "./Content";
import Sidebar from "./Sidebar";

// TODO: Make a variable for fixed navbar height.

// Determine which day's data to fetch.
// TODO: Fix time zone.
const date = new Date();
const date_display_options = {"weekday": "long", "month": "long", "day": "numeric","year": "numeric"};
const date_display = date.toLocaleDateString("en-US", date_display_options);
const date_key = date.toISOString().slice(0, 10);

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {data: []};
  }
  
  componentDidMount() {
    firebase.auth().signInAnonymously()
    .then(() => getData(date_key))
    .then(data => processData(data))
    .then(data => this.setState({data: data}))
    .catch(error => console.log(error));
  }
  
  render() {
    return (
      <div className="App">
      <Header data={this.state.data} date={date_display} />
      <Layout
        left={<Sidebar data={this.state.data} />}
        middle={<Content data={this.state.data} />}
      />
      </div>
      )
    }
}

export default App;
