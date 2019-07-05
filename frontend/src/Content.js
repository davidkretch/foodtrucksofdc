import React from "react";
import Card from "react-bootstrap/Card";
import ListGroup from "react-bootstrap/ListGroup";
import Jumbotron from "react-bootstrap/Jumbotron";

import Rating from "./Rating";

function Truck(props) {
  return (
    <ListGroup.Item className="d-flex justify-content-between align-items-center">
      {props.truck.name}
      <Rating rating={props.truck.rating} key={props.truck.name} />
    </ListGroup.Item>
  );
}

function Stop(props) {
  return (
    <div id={props.stop.link} style={{paddingTop: "56px", marginTop: "-56px"}}>
      <Card className="mt-4 mb-4">
        <Card.Header as="h5">{props.stop.name}</Card.Header>
        <ListGroup variant="flush">
          {props.stop.trucks.map((truck, index) => {
            return <Truck truck={truck} key={index} />
          })}
        </ListGroup>
      </Card>
    </div>
  );
}

function Stops(props) {
  return (
    props.stops.map(stop => {
        return <Stop stop={stop} key={stop.name} />
    })
  );
}

function NoData(props) {
  return (
    <Jumbotron style={{background: "none"}}>
      <h5 style={{textAlign: "center"}}>{props.status}</h5>
    </Jumbotron>
  );
}

function Content(props) {
  return (props.stops.length > 0 ? <Stops stops={props.stops} /> : <NoData status={props.status}/>);
}

export default Content;