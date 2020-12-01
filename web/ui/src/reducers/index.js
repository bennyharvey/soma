import { combineReducers } from 'redux'
import skuder from './skuder'
import { createBrowserHistory } from 'history'
import { connectRouter } from 'connected-react-router'
import login from './login'
import users from './users'
import persons from './persons'
import events from './events'

export const history = createBrowserHistory()

export default combineReducers({
    router: connectRouter(history),
    skuder,
    login,
    users,
    persons,
    events
})