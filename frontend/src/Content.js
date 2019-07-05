import React from "react";
import Card from "react-bootstrap/Card";
import ListGroup from "react-bootstrap/ListGroup";
import Jumbotron from "react-bootstrap/Jumbotron";

import Rating from "./Rating";

// Truck renders a single truck's information.
function Truck(props) {
  return (
    <ListGroup.Item className="d-flex justify-content-between align-items-center">
      {props.truck.name}
      <Rating
        rating={props.truck.rating}
        setRating={rating => {props.setRating(props.truck.name, rating)}}
        key={props.truck.name}
      />
    </ListGroup.Item>
  );
}

// Stop renders a single stop's list of trucks.
function Stop(props) {
  return (
    <div id={props.stop.link} style={{paddingTop: "56px", marginTop: "-56px"}}>
      <Card className="mt-4 mb-4">
        <Card.Header as="h5">{props.stop.name}</Card.Header>
        <ListGroup variant="flush">
          {props.stop.trucks.map((truck, index) => {
            return <Truck truck={truck} setRating={props.setRating} key={truck.name} />
          })}
        </ListGroup>
      </Card>
    </div>
  );
}

// Stops renders all stops.
function Stops(props) {
  return (
    props.stops.map(stop => {
        return <Stop stop={stop} setRating={props.setRating} key={stop.name} />
    })
  );
}

// NoData renders a message when there are no stops/trucks.
function NoData(props) {
  return (
    <Jumbotron style={{background: "none"}}>
      <h5 style={{textAlign: "center"}}>{props.status}</h5>
    </Jumbotron>
  );
}

function Content(props) {
  if (props.stops.length === 0) {
    return <NoData status={props.status}/>;
  }
  return <Stops stops={props.stops} setRating={props.setRating} />;
}

export default Content;