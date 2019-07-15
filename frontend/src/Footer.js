import React from "react";

import Container from "react-bootstrap/Container"
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";

import constants from "./constants";

function Footer(props) {
    const footerStyle = {
        width: "100%",
        height: "60px",
        lineHeight: "60px"
    };
    return (
        <footer className="footer" style={footerStyle}>
            <Container className="container-fluid m-0">
                <Row>
                    <Col {...constants.display.widthTotal} style={{textAlign: "center"}}>
                        <a href={"mailto:" + constants.email} className="text-secondary">Contact</a>
                    </Col>
                </Row>
            </Container>
        </footer>
    );
}

export default Footer;
