import React from "react";
import Nav from "react-bootstrap/Nav";
import Navbar from "react-bootstrap/Navbar";
import NavDropdown from "react-bootstrap/NavDropdown";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faTruck } from "@fortawesome/free-solid-svg-icons";
import { dateDisplay, dateKey, dateSequence } from "./date";

// DateDropdown renders a dropdown for selecting which date's food truck
// schedule the user would like to see.
function DateDropdown(props) {
  return (
    <Navbar.Collapse className="justify-content-end">
      <Nav className="d-none d-sm-block">
        <NavDropdown alignRight title={dateDisplay(props.selectedDate)}>
          {dateSequence(props.startDate).map(date => {
            return (
              <NavDropdown.Item
                key={dateKey(date)}
                eventKey={dateKey(date)}
                onClick={() => {props.setDate(date)}}>
                {dateDisplay(date)}
              </NavDropdown.Item>
            );
          })}
        </NavDropdown>
      </Nav>
    </Navbar.Collapse>
  );
}

// StopDropdown renders a dropdown for selecting which stop to jump to.
function StopDropdown(props) {
  return (
    <Nav className="d-block d-sm-none">
      <NavDropdown alignRight title="Stops">
        {props.stops.map(stop => {
          return (
            <NavDropdown.Item href={"#" + stop.link} key={stop.link}>{stop.abbrev}</NavDropdown.Item>
          );
        })}
      </NavDropdown>
    </Nav>
  );
}

// OptionalStopDropdown renders the stop dropdown, or nothing if there are
// no trucks.
function OptionalStopDropdown(props) {
  return (
    props.stops.length > 0 ? <StopDropdown stops={props.stops} /> : <div />
  );
}

// Header renders the header at the top of the page.
function Header(props) {
  return (
    <Navbar fixed="top" bg="dark" variant="dark">
      <Navbar.Brand href="/">
        <FontAwesomeIcon icon={faTruck} color="#a569bd"/>
        {" Food Trucks of DC "}
      </Navbar.Brand>
      <DateDropdown
        startDate={props.startDate}
        selectedDate={props.selectedDate}
        setDate={props.setDate}
      />
      <OptionalStopDropdown stops={props.stops} />
    </Navbar>
  );
}

export default Header;
