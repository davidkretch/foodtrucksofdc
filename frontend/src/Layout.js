import React from "react";

import Container from 'react-bootstrap/Container'
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'

function Layout(props) {
    return (
        <Container className="container-fluid m-0" style={{paddingTop: "56px"}}>
          <Row>
            <Col xs={12} md={3} xl={2}>
                {props.left}
            </Col>
            <Col xs={12} md={9} xl={8}>
                {props.middle}
            </Col>
            <Col xl={2}>
                {props.right}
            </Col>
          </Row>
        </Container>
    );
}

export default Layout;