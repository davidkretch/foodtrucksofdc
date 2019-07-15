import React from "react";

import Content from "./Content";
import Footer from "./Footer";
import Header from "./Header";
import Layout from "./Layout";
import Sidebar from "./Sidebar";

import "./App.css";

import { dateKey } from "./date";
import { firebase, getData, setRating } from "./firebase";
import { status, statusError } from "./status";

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
        const appStyle = {
            display: "flex",
            flexDirection: "column",
            minHeight: "100vh"
        };
        const sidebar = (
            <Sidebar stops={this.state.stops} />
        );
        const content = (
            <Content
                stops={this.state.stops}
                status={this.state.status}
                setRating={setRating}
            />
        );
        return (
            <div className="App" style={appStyle}>
                <Header
                    stops={this.state.stops}
                    startDate={date}
                    selectedDate={this.state.date}
                    setDate={(date) => this.setState({date: date})}
                />
                <Layout
                    left={sidebar}
                    center={content}
                />
                <Footer />
            </div>
        )
    }
}

export default App;
