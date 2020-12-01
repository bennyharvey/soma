import { connect } from 'react-redux'
import { Col, Container, Row } from 'react-bootstrap'
import React from 'react'
import { Redirect } from 'react-router-dom'
import './not-found.css'

const NotFound = ({ path, redirectOnRoot }) => {
    if (redirectOnRoot && path === '/') {
        return <Redirect to={redirectOnRoot} />
    }
    return (
        <Container id='not-found-container'>
            <Row>
                <Col className='text-center'>
                    <i className='text-muted'>Страница не найдена</i>
                </Col>
            </Row>
        </Container>
    )
}

const mapStateToProps = ({ router }) => ({
    path: router.location.pathname
})

export default connect(
    mapStateToProps
)(NotFound)