import React from "react";
import Card from "react-bootstrap/Card"
import ListGroup from "react-bootstrap/ListGroup"
import Jumbotron from "react-bootstrap/Jumbotron"

function Stops(props) {
  return (
    props.stops.map(stop => {
        return (
          <div id={stop.link} style={{paddingTop: "56px", marginTop: "-56px"}} key={stop.name}>
            <Card className="mt-4 mb-4">
              <Card.Header as="h5">{stop.name}</Card.Header>
              <ListGroup variant="flush">
                {stop.trucks.map((truck, index) => {
                  return (
                    <ListGroup.Item key={index}>{truck}</ListGroup.Item>
                  )
                })}
              </ListGroup>
            </Card>
          </div>
        )
    })
  );
}

function NoData() {
  return (
    <Jumbotron style={{background: "none"}}>
      <h5 style={{textAlign: "center"}}>Sorry, we don't have any info!</h5>
    </Jumbotron>
  );
}

function Content(props) {
  return (props.stops.length > 0 ? <Stops stops={props.stops} /> : <NoData />);
}

export default Content;