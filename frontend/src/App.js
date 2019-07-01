import React from "react";

import { firebase, getData, processData } from "./Firebase";
import Header from "./Header";
import Layout from "./Layout";
import Content from "./Content";
import Sidebar from "./Sidebar";
import { dateKey } from "./Date";

// TODO: Make a variable for fixed navbar height.
const date = new Date();

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      data: [],
      date: date
    };
  }

  // TODO: Put data fetch in a single location.
  componentDidMount() {
    firebase.auth().signInAnonymously()
    .then(() => getData(dateKey(this.state.date)))
    .then(data => processData(data))
    .then(data => this.setState({data: data}))
    .catch(error => console.log(error));
  }

  // TODO: Put data fetch in a single location.
  componentDidUpdate() {
    getData(dateKey(this.state.date))
    .then(data => processData(data))
    .then(data => this.setState({data: data}))
    .catch(error => console.log(error));
  }
  
  render() {
    return (
      <div className="App">
      <Header
        data={this.state.data}
        startDate={date}
        selectedDate={this.state.date}
        setDate={(date) => this.setState({date: date})}
      />
      <Layout
        left={<Sidebar data={this.state.data} />}
        middle={<Content data={this.state.data} />}
      />
      </div>
      )
    }
}

export default App;
