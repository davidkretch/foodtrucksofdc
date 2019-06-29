import React from "react";
import Nav from "react-bootstrap/Nav";
import Navbar from "react-bootstrap/Navbar";

function Sidebar(props) {
    return (
        <Navbar className="d-none d-sm-block mt-1" style={{position: "fixed"}}>
            <Nav className="mr-auto flex-sm-column">
                {props.data.map(stop => {
                    return (
                        <Nav.Link href={"#" + stop.link} key={stop.link}>{stop.abbrev}</Nav.Link>
                    );
                })}
            </Nav>
        </Navbar>
    );
}

export default Sidebar;