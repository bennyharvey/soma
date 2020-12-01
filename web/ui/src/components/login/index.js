import React from 'react'
import { connect } from 'react-redux'
import { login, INVALID_LOGIN_OR_PASSWORD, UNEXPECTED_RESPONSE_CODE, REQUEST_ERROR, UNKNOWN_ERROR } from '../../reducers/skuder'
import { Button, Card, Col, Container, Form, Row } from 'react-bootstrap'
import { setLogin, setPassword } from '../../reducers/login'
import queryString from 'query-string'
import './index.css'

const Login = props => {
    const onLogin = () => {
        props.onLogin({
            login: props.login,
            password: props.password,
            returnPath: props.returnPath
        })
    }

    const onLoginChange = (e) => {
        props.onLoginChange(e.currentTarget.value)
    }

    const onPasswordChange = (e) => {
        props.onPasswordChange(e.currentTarget.value)
    }

    let authError
    if (props.authError) {
        switch (props.authError.type) {
            case INVALID_LOGIN_OR_PASSWORD:
                authError = <div className='error' type='invalid'>Неверный логин или пароль</div>
                break
            case UNEXPECTED_RESPONSE_CODE:
                authError = <div className='error' type='invalid'>Неожиданный ответ от сервера: {props.authError.message}</div>
                break
            case REQUEST_ERROR:
                authError = <div className='error' type='invalid'>Ошибка выполнения запроса: {props.authError.message}</div>
                break
            case UNKNOWN_ERROR:
            default:
                authError = <div className='error' type='invalid'>Неизвестная ошибка: {props.authError.message}</div>
                break
        }
    }

    return (
        <Container id='login-container'>
            <Row>
                <Col md='4'>
                    <Card>
                        <Card.Body>
                            <Form>
                                <Form.Group controlId='login'>
                                    <Form.Label>Логин</Form.Label>
                                    <Form.Control placeholder='Введите логин' value={props.login} onChange={onLoginChange}/>
                                </Form.Group>
                                <Form.Group controlId='password'>
                                    <Form.Label>Пароль</Form.Label>
                                    <Form.Control type='password' placeholder='Введите пароль' value={props.password} onChange={onPasswordChange}/>
                                </Form.Group>
                                {authError}
                                <Button variant='primary' onClick={onLogin}>Войти</Button>
                            </Form>
                        </Card.Body>
                    </Card>
                </Col>
            </Row>
        </Container>
    )
}

const mapStateToProps = state => ({
    returnPath: queryString.parse(state.router.location.search).return_path,
    login: state.login.login,
    password: state.login.password,
    authError: state.skuder.authError
})

const mapDispatchToProps = dispatch => {
    return {
        onLogin: params => {
            dispatch(login(params))
        },
        onLoginChange: newLogin => {
            dispatch(setLogin(newLogin))
        },
        onPasswordChange: newPassword => {
            dispatch(setPassword(newPassword))
        }
    }
}

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(Login)