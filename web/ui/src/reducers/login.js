const SET_LOGIN = 'login/SET_LOGIN'
const SET_PASSWORD = 'login/SET_PASSWORD'

const initialState = {
    login: '',
    password: ''
}

export default (state = initialState, action) => {
    switch (action.type) {

        case SET_LOGIN:
            return {
                ...state,
                login: action.login
            }

        case SET_PASSWORD:
            return {
                ...state,
                password: action.password
            }

        default:
            return state
    }
}

export const setLogin = (login) => ({
    type: SET_LOGIN,
    login
})

export const setPassword = (password) => ({
    type: SET_PASSWORD,
    password
})