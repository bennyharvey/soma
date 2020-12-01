import React, { useEffect } from 'react'
import { connect } from 'react-redux'
import {
    loadUsers,
    newUser,
    createNewUser,
    cancelNewUser,
    setNewUserLogin,
    setNewUserPassword,
    setNewUserRole,
    editUser,
    saveEditUser,
    cancelEditUser,
    setEditUserPassword,
    setEditUserRole,
    removeUser
} from '../../reducers/users'
import { Button, Col, Container, Row, Table, Modal, Form } from 'react-bootstrap'
import './index.css'

const NewUserModal = (props) => {

    if (!props.newUser) {
        return ''
    }

    const onLoginChange = (e) => {
        props.onNewUserLoginChange(e.currentTarget.value)
    }

    const onPasswordChange = (e) => {
        props.onNewUserPasswordChange(e.currentTarget.value)
    }

    const onRoleChange = (e) => {
        props.onNewUserRoleChange(e.currentTarget.value)
    }

    let roleRadios = []
    for (const roleID in props.roles) {
        roleRadios.push(
            <Form.Check custom type='radio' id={'radio-role-'+roleID} key={roleID} label={props.roles[roleID]}
                        value={roleID} checked={props.newUser.role === roleID} onChange={onRoleChange}/>
        )
    }

    let error = ''
    if (props.newUserError) {
        error = <div style={{color: 'red'}}>{props.newUserError}</div>
    }

    return (
        <Modal show={true} onHide={props.onCancelNewUser} animation={false}>
            <Modal.Header closeButton>
                <Modal.Title>Новый пользователь</Modal.Title>
            </Modal.Header>
            <Modal.Body>
                <Form>
                    <Form.Group controlId='new-user-login'>
                        <Form.Label>Логин</Form.Label>
                        <Form.Control autoComplete='new-password' placeholder='Введите логин' value={props.newUser.login} onChange={onLoginChange}/>
                        <Form.Text className='text-muted'>
                            Логин должен быть уникальным.
                        </Form.Text>
                    </Form.Group>
                    <Form.Group controlId='new-user-password'>
                        <Form.Label>Пароль</Form.Label>
                        <Form.Control autoComplete='new-password' type='password' placeholder='Введите пароль' value={props.newUser.password} onChange={onPasswordChange}/>
                        <Form.Text className='text-muted'>
                            Пароль должен состоять из не менее чем 16 символов.
                        </Form.Text>
                    </Form.Group>
                    <Form.Group controlId='new-user-role'>
                        <Form.Label>Роль</Form.Label>
                        {roleRadios}
                        <Form.Text className='text-muted'>
                            Нужно выбрать одну из ролей.
                        </Form.Text>
                    </Form.Group>
                    {error}
                </Form>
            </Modal.Body>
            <Modal.Footer>
                <Button variant='secondary' onClick={props.onCancelNewUser}>Отмена</Button>
                <Button variant='primary' onClick={props.onCreateNewUser}>Создать</Button>
            </Modal.Footer>
        </Modal>
    )
}

const EditUserModal = (props) => {

    if (!props.editUser) {
        return ''
    }

    const onPasswordChange = (e) => {
        props.onEditUserPasswordChange(e.currentTarget.value)
    }

    const onRoleChange = (e) => {
        props.onEditUserRoleChange(e.currentTarget.value)
    }

    let roleRadios = []
    for (const roleID in props.roles) {
        roleRadios.push(
            <Form.Check custom type='radio' id={'radio-role-'+roleID} key={roleID} label={props.roles[roleID]}
                        value={roleID} checked={props.editUser.role === roleID} onChange={onRoleChange}/>
        )
    }

    let error = ''
    if (props.editUserError) {
        error = <div style={{color: 'red'}}>{props.editUserError}</div>
    }

    return (
        <Modal show={true} onHide={props.onCancelEditUser} animation={false}>
            <Modal.Header closeButton>
                <Modal.Title>Пользователь {props.editUser.login}</Modal.Title>
            </Modal.Header>
            <Modal.Body>
                <Form>
                    <Form.Group controlId='edit-user-password'>
                        <Form.Label>Пароль</Form.Label>
                        <Form.Control autoComplete='new-password' type='password' placeholder='Введите пароль' value={props.editUser.password} onChange={onPasswordChange}/>
                        <Form.Text className='text-muted'>
                            Пароль должен состоять из не менее чем 16 символов. Оставить пустым чтобы оставить прежний пароль.
                        </Form.Text>
                    </Form.Group>
                    <Form.Group controlId='edit-user-role'>
                        <Form.Label>Роль</Form.Label>
                        {roleRadios}
                        <Form.Text className='text-muted'>
                            Нужно выбрать одну из ролей.
                        </Form.Text>
                    </Form.Group>
                    {error}
                </Form>
            </Modal.Body>
            <Modal.Footer>
                <Button variant='secondary' onClick={props.onCancelEditUser}>Отмена</Button>
                <Button variant='primary' onClick={props.onSaveEditUser}>Сохранить</Button>
            </Modal.Footer>
        </Modal>
    )
}

const Users = ({ loadUsers, ...props }) => {

    useEffect(() => {
        loadUsers()
    }, [loadUsers])

    if (props.users === null) {
        return ''
    }

    return (
        <Container id='users-container'>
            <Row>
                <Col md='12'>
                    <Button onClick={props.onNewUser}>Новый пользователь</Button>
                </Col>
            </Row>
            <Row>
                <Col md='12'>
                    <Table size='sm' striped bordered>
                        <thead>
                            <tr>
                                <th>Логин</th>
                                <th>Роль</th>
                                <th className='manage-col'>Управление</th>
                            </tr>
                        </thead>
                        <tbody>
                            {props.users.map(u =>
                                <tr key={u.login}>
                                    <td>{u.login}</td>
                                    <td>{props.roles[u.role] ? props.roles[u.role] : u.role}</td>
                                    <td className='manage-col'>
                                        <Button size='sm' variant='primary' onClick={props.onEditUser.bind(props, u.login)}>Редактировать</Button>
                                        <Button size='sm' variant='danger' onClick={props.onRemoveUser.bind(props, u.login)}>Удалить</Button>
                                    </td>
                                </tr>
                            )}
                        </tbody>
                    </Table>
                </Col>
            </Row>
            <NewUserModal {...props} />
            <EditUserModal {...props} />
        </Container>
    )
}

const mapStateToProps = state => ({
    ...state.users,
    roles: state.skuder.roles,
})

const mapDispatchToProps = dispatch => {
    return {
        loadUsers: () => {
            dispatch(loadUsers())
        },
        onNewUser: () => {
            dispatch(newUser())
        },
        onCreateNewUser: () => {
            dispatch(createNewUser())
        },
        onCancelNewUser: () => {
            dispatch(cancelNewUser())
        },
        onNewUserLoginChange: (login) => {
            dispatch(setNewUserLogin(login))
        },
        onNewUserPasswordChange: (password) => {
            dispatch(setNewUserPassword(password))
        },
        onNewUserRoleChange: (role) => {
            dispatch(setNewUserRole(role))
        },
        onEditUser: (login) => {
            dispatch(editUser(login))
        },
        onSaveEditUser: () => {
            dispatch(saveEditUser())
        },
        onCancelEditUser: () => {
            dispatch(cancelEditUser())
        },
        onEditUserPasswordChange: (password) => {
            dispatch(setEditUserPassword(password))
        },
        onEditUserRoleChange: (role) => {
            dispatch(setEditUserRole(role))
        },
        onRemoveUser: (login) => {
            dispatch(removeUser(login))
        }
    }
}

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(Users)