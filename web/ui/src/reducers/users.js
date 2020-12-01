import * as axios from 'axios'
import { handleAuthError, usersURL } from './skuder'

const SET_USERS = 'users/SET_USERS'

const NEW_USER = 'users/NEW_USER'
const CLEAR_NEW_USER = 'users/CLEAR_NEW_USER'
const SET_NEW_USER_LOGIN = 'users/SET_NEW_USER_LOGIN'
const SET_NEW_USER_PASSWORD = 'users/SET_NEW_USER_PASSWORD'
const SET_NEW_USER_ROLE = 'users/SET_NEW_USER_ROLE'
const SET_NEW_USER_ERROR = 'users/SET_NEW_USER_ERROR'

const EDIT_USER = 'users/EDIT_USER'
const CLEAR_EDIT_USER = 'users/CLEAR_EDIT_USER'
const SET_EDIT_USER_PASSWORD = 'users/SET_EDIT_USER_PASSWORD'
const SET_EDIT_USER_ROLE = 'users/SET_EDIT_USER_ROLE'
const SET_EDIT_USER_ERROR = 'users/SET_EDIT_USER_ERROR'

const initialState = {
    users: null,
    newUser: null,
    newUserError: null,
    editUser: null,
    editUserError: null
}

export default (state = initialState, action) => {
    switch (action.type) {

        case SET_USERS:
            return {
                ...state,
                users: action.users
            }

        case NEW_USER:
            return {
                ...state,
                newUser: {
                    login: '',
                    password: '',
                    role: null
                }
            }

        case CLEAR_NEW_USER:
            return {
                ...state,
                newUser: null,
            }

        case SET_NEW_USER_LOGIN:
            return {
                ...state,
                newUser: {
                    ...state.newUser,
                    login: action.login
                }
            }

        case SET_NEW_USER_PASSWORD:
            return {
                ...state,
                newUser: {
                    ...state.newUser,
                    password: action.password
                }
            }

        case SET_NEW_USER_ROLE:
            return {
                ...state,
                newUser: {
                    ...state.newUser,
                    role: action.role
                }
            }

        case SET_NEW_USER_ERROR:
            return {
                ...state,
                newUserError: action.newUserError,
            }

        case EDIT_USER:
            return {
                ...state,
                editUser: action.editUser
            }

        case CLEAR_EDIT_USER:
            return {
                ...state,
                editUser: null
            }

        case SET_EDIT_USER_PASSWORD:
            return {
                ...state,
                editUser: {
                    ...state.editUser,
                    password: action.password
                }
            }

        case SET_EDIT_USER_ROLE:
            return {
                ...state,
                editUser: {
                    ...state.editUser,
                    role: action.role
                }
            }

        case SET_EDIT_USER_ERROR:
            return {
                ...state,
                editUserError: action.editUserError,
            }

        default:
            return state
    }
}

export const loadUsers = () => {
    return (dispatch, getState) => {
        axios.get(usersURL).then(res => {
            dispatch({ type: SET_USERS, users: res.data })
        }).catch(err => {
            handleAuthError(dispatch, getState, err)
        })
    }
}

export const newUser = () => ({
    type: NEW_USER,
})

export const createNewUser = () => {
    return (dispatch, getState) => {
        const newUser = getState().users.newUser

        newUser.login = newUser.login.trim()
        if (newUser.login === '') {
            dispatch(setNewUserError('Пустой логин'))
            return
        }

        if (newUser.password === '') {
            dispatch(setNewUserError('Пустой пароль'))
            return
        }

        if (newUser.password.length < 16) {
            dispatch(setNewUserError('Пароль слишком короткий'))
            return
        }

        if (!newUser.role) {
            dispatch(setNewUserError('Не выбрана роль'))
            return
        }

        axios.post(usersURL, newUser).then(() => {
            dispatch(cancelNewUser())
            dispatch(loadUsers())
        }).catch(err => {
            if (handleAuthError(dispatch, getState, err)) {
                return
            }
            dispatch(setNewUserError(err.toString()))
        })
    }
}

export const cancelNewUser = () => ({
    type: CLEAR_NEW_USER
})

export const setNewUserLogin = (login) => ({
    type: SET_NEW_USER_LOGIN,
    login
})

export const setNewUserPassword = (password) => ({
    type: SET_NEW_USER_PASSWORD,
    password
})

export const setNewUserRole = (role) => ({
    type: SET_NEW_USER_ROLE,
    role
})

const setNewUserError = (newUserError) => ({
    type: SET_NEW_USER_ERROR,
    newUserError
})

export const editUser = login => {
    return (dispatch, getState) => {
        axios.get(`${usersURL}/${login}`).then(res => {
            dispatch({ type: EDIT_USER, editUser: res.data })
        }).catch(err => {
            handleAuthError(dispatch, getState, err)
        })
    }
}

export const saveEditUser = () => {
    return (dispatch, getState) => {

        const editUser = getState().users.editUser

        if (editUser.password.length > 0 && editUser.password.length < 16) {
            dispatch(setEditUserError('Пароль слишком короткий'))
            return
        }

        axios.put(usersURL, editUser).then(() => {
            dispatch(cancelEditUser())
            dispatch(loadUsers())
        }).catch(err => {
            handleAuthError(dispatch, getState, err)
        })
    }
}

export const cancelEditUser = () => ({
    type: CLEAR_EDIT_USER,
})

export const setEditUserPassword = (password) => ({
    type: SET_EDIT_USER_PASSWORD,
    password
})

export const setEditUserRole = (role) => ({
    type: SET_EDIT_USER_ROLE,
    role
})

const setEditUserError = (editUserError) => ({
    type: SET_EDIT_USER_ERROR,
    editUserError
})

export const removeUser = (login) => {
    return (dispatch, getState) => {
        axios.delete(`${usersURL}/${login}`).then(() => {
            dispatch(loadUsers())
        }).catch(err => {
            handleAuthError(dispatch, getState, err)
        })
    }
}