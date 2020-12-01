import * as axios from 'axios'
import { push } from 'connected-react-router'
import queryString from 'query-string'
import Cookie from 'js-cookie'

export const ADMIN = 'admin'
export const SECURITY = 'security'

export const INVALID_LOGIN_OR_PASSWORD = 'invalid_login_or_password'
export const UNEXPECTED_RESPONSE_CODE = 'unexpected_response_code'
export const REQUEST_ERROR = 'request_error'
export const UNKNOWN_ERROR = 'unknown_error'

const baseURL = 'https://localhost'
const apiBaseURL = `${baseURL}/api`
const loginURL = `${apiBaseURL}/login`
export const usersURL = `${apiBaseURL}/users`
export const personsURL = `${apiBaseURL}/persons`
export const photosURL = `${apiBaseURL}/photos`
export const personFacesURL = `${apiBaseURL}/person_faces`
export const eventsURL = `${apiBaseURL}/events`
export const passageNamesURL = `${apiBaseURL}/passage_names`

const LOGIN = 'skuder/LOGIN'
const LOGOUT = 'skuder/LOGOUT'

const CHANGE_AUTH_ERROR = 'skuder/CHANGE_AUTH_ERROR'

const SKUDER_USER = 'skuder_user'

const AUTH_COOKIE = 'auth'

let user = null
if (Cookie.get(AUTH_COOKIE)) {
    user = JSON.parse(localStorage.getItem(SKUDER_USER))
}

const roles = {}
roles[ADMIN] = 'Админ'
roles[SECURITY] = 'Безопасность'

const initialState = {
    user: user,
    authError: null,
    roles: roles
}

export default (state = initialState, action) => {
    switch (action.type) {

        case LOGIN:
            return {
                ...state,
                user: action.user,
            }

        case LOGOUT:
            return {
                ...state,
                user: null,
            }

        case CHANGE_AUTH_ERROR:
            return {
                ...state,
                authError: action.authError
            }

        default:
            return state
    }
}

export const login = ({ login, password, returnPath }) => {
    return dispatch => {
        axios.post(loginURL, {
            login,
            password,
        }).then(res => {
            dispatch({ type: LOGIN, user: res.data.user })
            dispatch(push(returnPath || ''))
            Cookie.set(AUTH_COOKIE, res.data.token)
            localStorage.setItem(SKUDER_USER, JSON.stringify(res.data.user))
        }).catch(error => {
            if (error.response) {
                if (error.response.status === 401) {
                    dispatch(onAuthErrorChange({
                        type: INVALID_LOGIN_OR_PASSWORD,
                        message: error.message
                    }))
                } else {
                    dispatch(onAuthErrorChange({
                        type: UNEXPECTED_RESPONSE_CODE,
                        message: error.message
                    }))
                }
            } else if (error.request) {
                dispatch(onAuthErrorChange({
                    type: REQUEST_ERROR,
                    message: error.message
                }))
            } else {
                dispatch(onAuthErrorChange({
                    type: UNKNOWN_ERROR,
                    message: error.message
                }))
            }
        })
    }
}

export const logout = () => {
    return (dispatch, getState) => {
        const { router } = getState()
        const queryParams = queryString.stringify({
            return_path: router.location.pathname + router.location.search
        })

        Cookie.remove(AUTH_COOKIE)
        localStorage.removeItem(SKUDER_USER)
        dispatch({ type: LOGOUT })
        dispatch(push('/login?'+queryParams))
    }
}

const onAuthErrorChange = authError => {
    return {
        type: CHANGE_AUTH_ERROR,
        authError
    }
}

export const handleAuthError = (dispatch, state, err) => {
    if (err.response) {
        if (err.response.status === 401) {
            dispatch(logout())
            return true
        }
    }
    return false
}