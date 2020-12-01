import * as axios from 'axios'
import { handleAuthError, eventsURL, passageNamesURL } from './skuder'
import debounce from 'debounce'
import { push } from 'connected-react-router'
import querystring from 'querystring'
import { formatDateTime } from '../utils'

const SET_EVENTS = 'events/SET_EVENTS'
const SET_FROM = 'events/SET_FROM'
const SET_TO = 'events/SET_TO'
const SET_PASSAGE_ID = 'events/SET_PASSAGE_ID'
const SET_PERSON_NAME = 'events/SET_PERSON_NAME'
const SET_PAGE = 'events/SET_PAGE'
const SET_PASSAGE_NAMES = 'events/SET_PASSAGE_NAMES'

const query = querystring.decode(window.location.search.substring(1,
    window.location.search.length))

const initialState = {
    events: null,
    from: query.from ? new Date(query.from) : null,
    to: query.to ? new Date(query.to) : null,
    passageID: query.passage_id || null,
    personName: query.person_name || '',
    orderBy: 'time',
    orderDirection: 'desc',
    page: 1,
    recordsPerPages: 100,
    passageNames: {},
}

export default (state = initialState, action) => {
    switch (action.type) {

        case SET_EVENTS:
            return {
                ...state,
                events: action.events
            }

        case SET_FROM:
            return {
                ...state,
                from: action.from
            }

        case SET_TO:
            return {
                ...state,
                to: action.to
            }

        case SET_PASSAGE_ID:
            return {
                ...state,
                passageID: action.passageID
            }

        case SET_PERSON_NAME:
            return {
                ...state,
                personName: action.personName
            }

        case SET_PAGE:
            return {
                ...state,
                page: action.page
            }

        case SET_PASSAGE_NAMES:
            return {
                ...state,
                passageNames: action.passageNames
            }

        default:
            return state
    }
}

const _loadEvents = (dispatch, getState) => {
    const state = getState()
    axios.get(eventsURL, {
        params: {
            from: state.events.from,
            to: state.events.to,
            passage_id: state.events.passageID,
            person_name: state.events.personName || null,
            order_by: state.events.orderBy,
            order_direction: state.events.orderDirection,
            limit: state.events.recordsPerPages,
            offset: (state.events.page-1)*state.events.recordsPerPages,
        }
    }).then(res => {
        dispatch({ type: SET_EVENTS, events: res.data })
    }).catch(err => {
        handleAuthError(dispatch, getState, err)
    })
}

export const loadEvents = () => {
    return _loadEvents
}

export const loadPassageNames = () => {
    return (dispatch, getState) => {
        axios.get(passageNamesURL).then(res => {
            dispatch({ type: SET_PASSAGE_NAMES, passageNames: res.data })
        })
    }
}

const loadEventsDebounced = debounce(_loadEvents, 100)

export const setFrom = (from) => {
    return (dispatch, getState) => {
        dispatch({ type: SET_FROM, from })

        const state = getState()

        const path = state.router.location.pathname
        const query = querystring.decode(state.router.location.search.substring(1,
            window.location.search.length))

        if (from) {
            query.from = formatDateTime(from)
        } else {
            delete query.from
        }

        if (from && state.events.to && from > state.events.to) {
            dispatch(setTo(from))
        }

        dispatch(push(path + '?' + querystring.stringify(query)))

        loadEventsDebounced(dispatch, getState)
    }
}

export const setTo = (to) => {
    return (dispatch, getState) => {
        dispatch({ type: SET_TO, to })

        const state = getState()

        const path = state.router.location.pathname
        const query = querystring.decode(state.router.location.search.substring(1,
            window.location.search.length))

        if (to) {
            query.to = formatDateTime(to)
        } else {
            delete query.to
        }

        if (to && state.events.from && to < state.events.from) {
            dispatch(setFrom(to))
        }

        dispatch(push(path + '?' + querystring.stringify(query)))

        loadEventsDebounced(dispatch, getState)
    }
}

export const setPassageID = (passageID) => {
    return (dispatch, getState) => {
        dispatch({ type: SET_PASSAGE_ID, passageID })

        const path = getState().router.location.pathname
        const query = querystring.decode(getState().router.location.search.substring(1,
            window.location.search.length))

        if (passageID) {
            query.passage_id = passageID
        } else {
            delete query.passage_id
        }

        dispatch(push(path + '?' + querystring.stringify(query)))

        loadEventsDebounced(dispatch, getState)
    }
}

export const setPersonName = (personName) => {
    return (dispatch, getState) => {
        dispatch({ type: SET_PERSON_NAME, personName })

        const path = getState().router.location.pathname
        const query = querystring.decode(getState().router.location.search.substring(1,
            window.location.search.length))

        if (personName) {
            query.person_name = personName
        } else {
            delete query.person_name
        }

        dispatch(push(path + '?' + querystring.stringify(query)))

        loadEventsDebounced(dispatch, getState)
    }
}
export const setPage = (page) => {
    return (dispatch, getState) => {
        if (page < 1) {
            return
        }

        const state = getState()

        const currentPage = state.events.page
        const currentEvents = state.events.events
        if (currentPage < page && currentEvents.length === 0) {
            return
        }

        dispatch({ type: SET_PAGE, page })

        const path = state.router.location.pathname
        const query = querystring.decode(state.router.location.search.substring(1,
            window.location.search.length))

        if (page && page !== 1) {
            query.page = page
        } else {
            delete query.page
        }

        dispatch(push(path + '?' + querystring.stringify(query)))

        loadEventsDebounced(dispatch, getState)
    }
}