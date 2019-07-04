import React from "react";

import Content from "./Content";
import Header from "./Header";
import Layout from "./Layout";
import Sidebar from "./Sidebar";

import { dateKey } from "./date";
import { firebase, getData, processData } from "./firebase";
import { status, statusError } from "./status"

// TODO: Make a variable for fixed navbar height.
const date = new Date();

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      stops: [],
      date: date,
      status: "Loading..."
    };
  }

  // TODO: Put data fetch in a single location.
  componentDidMount() {
    firebase.auth().signInAnonymously()
    .then(() => getData(dateKey(this.state.date)))
    .then(data => processData(data))
    .then(data => this.setState({stops: data, status: status(data)}))
    .catch(error => this.setState({status: statusError()}));
  }

  // TODO: Put data fetch in a single location.
  componentDidUpdate() {
    getData(dateKey(this.state.date))
    .then(data => processData(data))
    .then(data => this.setState({stops: data, status: status(data)}))
    .catch(error => this.setState({status: statusError()}));
  }
  
  render() {
    return (
      <div className="App">
      <Header
        stops={this.state.stops}
        startDate={date}
        selectedDate={this.state.date}
        setDate={(date) => this.setState({date: date})}
      />
      <Layout
        left={<Sidebar stops={this.state.stops} />}
        middle={<Content stops={this.state.stops} status={this.state.status} />}
      />
      </div>
      )
    }
}

export default App;
