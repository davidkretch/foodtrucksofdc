import React from "react";
import Nav from "react-bootstrap/Nav";
import Navbar from "react-bootstrap/Navbar";
import NavDropdown from "react-bootstrap/NavDropdown";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faTruck } from "@fortawesome/free-solid-svg-icons";

function Date(props) {
  return (
    <Navbar.Collapse className="justify-content-end">
      <Navbar.Text className="d-none d-sm-block">
        {props.date}
      </Navbar.Text>
    </Navbar.Collapse>
  );
}

function StopDropdown(props) {
  return (
    <Nav className="d-block d-sm-none">
      <NavDropdown alignRight title="Stops">
        {props.data.map(stop => {
          return (
            <NavDropdown.Item href={"#" + stop.link} key={stop.link}>{stop.abbrev}</NavDropdown.Item>
          );
        })}
      </NavDropdown>
    </Nav>
  );
}

function OptionalStopDropdown(props) {
  return (
    props.data.length > 0 ? <StopDropdown data={props.data} /> : <div />
  );
}

function Header(props) {
  return (
    <Navbar fixed="top" bg="dark" variant="dark">
      <Navbar.Brand href="/">
        <FontAwesomeIcon icon={faTruck} color="#a569bd"/>
        {" Food Trucks of DC "}
      </Navbar.Brand>
      <Date date={props.date} />
      <OptionalStopDropdown data={props.data} />
    </Navbar>
  );
}

export default Header;
