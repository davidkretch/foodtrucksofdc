import React from "react";

import Container from "react-bootstrap/Container";
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";

import constants from "./constants";
const height = constants.display.heightNavbar;

function Layout(props) {
    const layoutStyle = {flex: 1, paddingTop: height+"px"};
    return (
        <Container className="container-fluid m-0" style={layoutStyle}>
            <Row>
                <Col {...constants.display.widthLeft}>
                    {props.left}
                </Col>
                <Col {...constants.display.widthCenter}>
                    {props.center}
                </Col>
            </Row>
        </Container>
    );
}

export default Layout;