import React from "react";
import Nav from "react-bootstrap/Nav";
import Navbar from "react-bootstrap/Navbar";

// Sidebar renders the sidebar with a list of stops, and links for jumping
// to a specific stop.
function Sidebar(props) {
    return (
        <Navbar className="d-none d-sm-block mt-1" style={{position: "fixed"}}>
            <Nav className="mr-auto flex-sm-column">
                {props.stops.map(stop => {
                    return (
                        <Nav.Link href={"#" + stop.link} key={stop.link}>{stop.abbrev}</Nav.Link>
                    );
                })}
            </Nav>
        </Navbar>
    );
}

export default Sidebar;