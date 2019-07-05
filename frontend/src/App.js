import React from "react";

import Content from "./Content";
import Header from "./Header";
import Layout from "./Layout";
import Sidebar from "./Sidebar";

import "./App.css";

import { dateKey } from "./date";
import { firebase, getData } from "./firebase";
import { status, statusError } from "./status";


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

  fetch() {
    getData(dateKey(this.state.date))
    .then(stops => this.setState({stops: stops, status: status(stops)}))
    .catch(err => {
      this.setState({status: statusError()});
      console.log(err);
    });
  }

  componentDidMount() {
    firebase.auth().signInAnonymously()
    .then(() => this.fetch());
  }

  componentDidUpdate(prepProps, prevState) {
    if (this.state.date !== prevState.date) {
      this.fetch();
    }
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
        middle={<Content stops={this.state.stops} trucks={this.state.trucks} status={this.state.status} />}
      />
      </div>
      )
    }
}

export default App;
