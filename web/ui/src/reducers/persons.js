import * as axios from 'axios'
import { personFacesURL, handleAuthError, personsURL } from './skuder'

const SET_PERSONS = 'persons/SET_PERSONS'
const ADD_PERSON = 'persons/ADD_PERSON'

const NEW_PERSON = 'persons/NEW_PERSON'
const SET_NEW_PERSON = 'persons/SET_NEW_PERSON'
const CLEAR_NEW_PERSON = 'persons/CLEAR_NEW_PERSON'
const SET_NEW_PERSON_NAME = 'persons/SET_NEW_PERSON_NAME'
const SET_NEW_PERSON_POSITION = 'persons/SET_NEW_PERSON_POSITION'
const SET_NEW_PERSON_UNIT = 'persons/SET_NEW_PERSON_UNIT'
const ADD_NEW_PERSON_PHOTO = 'persons/ADD_NEW_PERSON_PHOTO'
const SET_NEW_PERSON_PHOTO_LOADED = 'persons/SET_NEW_PERSON_PHOTO_LOADED'
const SET_NEW_PERSON_PHOTO_ERROR = 'persons/SET_NEW_PERSON_PHOTO_ERROR'
const REMOVE_NEW_PERSON_PHOTO = 'persons/REMOVE_NEW_PERSON_PHOTO'
const SET_NEW_PERSON_ERROR = 'persons/SET_NEW_PERSON_ERROR'

const EDIT_PERSON = 'persons/EDIT_PERSON'
const CLEAR_EDIT_PERSON = 'persons/CLEAR_EDIT_PERSON'
const SET_EDIT_PERSON_NAME = 'persons/SET_EDIT_PERSON_NAME'
const SET_EDIT_PERSON_POSITION = 'persons/SET_EDIT_PERSON_POSITION'
const SET_EDIT_PERSON_UNIT = 'persons/SET_EDIT_PERSON_UNIT'
const SET_EDIT_PERSON_FACE_TO_REMOVE = 'persons/SET_EDIT_PERSON_FACE_TO_REMOVE'
const RESTORE_EDIT_PERSON_FACE = 'persons/RESTORE_EDIT_PERSON_FACE'
const ADD_EDIT_PERSON_FACE = 'persons/ADD_EDIT_PERSON_FACE'
const REMOVE_EDIT_PERSON_FACE = 'persons/REMOVE_EDIT_PERSON_FACE'
const SET_EDIT_PERSON_FACE_ERROR = 'persons/SET_EDIT_PERSON_FACE_ERROR'
const ADD_EDIT_PERSON_PHOTO = 'persons/ADD_EDIT_PERSON_PHOTO'
const REMOVE_EDIT_PERSON_PHOTO = 'persons/REMOVE_EDIT_PERSON_PHOTO'
const SET_EDIT_PERSON_PHOTO_ERROR = 'persons/SET_EDIT_PERSON_PHOTO_ERROR'
const SET_EDIT_PERSON_ERROR = 'persons/SET_EDIT_PERSON_ERROR'


const initialState = {
    persons: null,
    newPerson: null,
    newPersonPhotos: null,
    newPersonError: null,
    editPerson: null,
    editPersonFaces: null,
    editPersonPhotos: null,
    editPersonError: null,
}

export default (state = initialState, action) => {
    switch (action.type) {

        case SET_PERSONS:
            return {
                ...state,
                persons: action.persons
            }

        case ADD_PERSON:
            return {
                ...state,
                persons: [...state.persons, action.person]
            }

        case NEW_PERSON:
            return {
                ...state,
                newPerson: {
                    name: '',
                    position: '',
                    unit: ''
                },
                newPersonPhotos: []
            }

        case SET_NEW_PERSON:
            return {
                ...state,
                newPerson: action.newPerson
            }

        case CLEAR_NEW_PERSON:
            return {
                ...state,
                newPerson: null,
                newPersonPhotos: null,
                newPersonError: null
            }

        case SET_NEW_PERSON_NAME:
            return {
                ...state,
                newPerson: {
                    ...state.newPerson,
                    name: action.name
                }
            }

        case SET_NEW_PERSON_POSITION:
            return {
                ...state,
                newPerson: {
                    ...state.newPerson,
                    position: action.position
                }
            }

        case SET_NEW_PERSON_UNIT:
            return {
                ...state,
                newPerson: {
                    ...state.newPerson,
                    unit: action.unit
                }
            }

        case ADD_NEW_PERSON_PHOTO:
            return {
                ...state,
                newPersonPhotos: [...state.newPersonPhotos, action.photo]
            }

        case SET_NEW_PERSON_PHOTO_LOADED:
            return {
                ...state,
                newPersonPhotos: state.newPersonPhotos.map(p => {
                    if (p === action.photo) {
                        return {
                            ...p,
                            loaded: true
                        }
                    }
                    return p
                })
            }

        case SET_NEW_PERSON_PHOTO_ERROR:
            return {
                ...state,
                newPersonPhotos: state.newPersonPhotos.map(p => {
                    if (p === action.photo) {
                        return {
                            ...p,
                            error: action.error
                        }
                    }
                    return p
                })
            }

        case REMOVE_NEW_PERSON_PHOTO:
            return {
                ...state,
                newPersonPhotos: state.newPersonPhotos.filter(
                    p => p !== action.photo)
            }

        case SET_NEW_PERSON_ERROR:
            return {
                ...state,
                newPersonError: action.newPersonError,
            }

        case EDIT_PERSON:
            return {
                ...state,
                editPerson: action.editPerson,
                editPersonFaces: action.editPersonFaces,
                editPersonPhotos: []
            }

        case CLEAR_EDIT_PERSON:
            return {
                ...state,
                editPerson: null,
                editPersonFaces: null,
                editPersonPhotos: null,
                editPersonError: null
            }

        case SET_EDIT_PERSON_NAME:
            return {
                ...state,
                editPerson: {
                    ...state.editPerson,
                    name: action.name
                }
            }

        case SET_EDIT_PERSON_POSITION:
            return {
                ...state,
                editPerson: {
                    ...state.editPerson,
                    position: action.position
                }
            }

        case SET_EDIT_PERSON_UNIT:
            return {
                ...state,
                editPerson: {
                    ...state.editPerson,
                    unit: action.unit
                }
            }

        case SET_EDIT_PERSON_FACE_TO_REMOVE:
            return {
                ...state,
                editPersonFaces: state.editPersonFaces.map(f => {
                    if (f.id === action.personFaceID) {
                        return {
                            ...f,
                            toRemove: true
                        }
                    }
                    return f
                })
            }

        case RESTORE_EDIT_PERSON_FACE:
            return {
                ...state,
                editPersonFaces: state.editPersonFaces.map(f => {
                    if (f.id === action.personFaceID) {
                        return {
                            ...f,
                            toRemove: false
                        }
                    }
                    return f
                })
            }

        case REMOVE_EDIT_PERSON_FACE:
            return {
                ...state,
                editPersonFaces: state.editPersonFaces.filter(f => f.id !== action.personFaceID)
            }

        case SET_EDIT_PERSON_FACE_ERROR:
            return {
                ...state,
                editPersonFaces: state.editPersonFaces.map(f => {
                    if (f.id === action.personFaceID) {
                        return {
                            ...f,
                            error: action.error
                        }
                    }
                    return f
                })
            }

        case ADD_EDIT_PERSON_FACE:
            return {
                ...state,
                editPersonFaces: [...state.editPersonFaces, action.personFace]
            }

        case ADD_EDIT_PERSON_PHOTO:
            return {
                ...state,
                editPersonPhotos: [
                    ...state.editPersonPhotos,
                    action.photo
                ]
            }

        case REMOVE_EDIT_PERSON_PHOTO:
            return {
                ...state,
                editPersonPhotos: state.editPersonPhotos.filter(p => p !== action.photo)
            }

        case SET_EDIT_PERSON_PHOTO_ERROR:
            return {
                ...state,
                editPersonPhotos: state.editPersonPhotos.map(p => {
                    if (p === action.photo) {
                        return {
                            ...p,
                            error: action.error
                        }
                    }
                    return p
                })
            }

        case SET_EDIT_PERSON_ERROR:
            return {
                ...state,
                editPersonError: action.editPersonError,
            }

        default:
            return state
    }
}

export const loadPersons = () => {
    return (dispatch, getState) => {
        axios.get(personsURL).then(res => {
            dispatch({ type: SET_PERSONS, persons: res.data })
        }).catch(err => {
            handleAuthError(dispatch, getState, err)
        })
    }
}

const addPerson = (person) => ({
    type: ADD_PERSON,
    person
})

export const newPerson = () => ({
    type: NEW_PERSON,
})

const setNewPerson = (newPerson) => ({
    type: SET_NEW_PERSON,
    newPerson
})

export const createNewPerson = () => {
    return (dispatch, getState) => {
        const newPerson = getState().persons.newPerson

        newPerson.name = newPerson.name.trim()
        if (newPerson.name === '') {
            dispatch(setNewPersonError('Пустое имя'))
            return
        }

        newPerson.position = newPerson.position.trim()
        if (newPerson.position === '') {
            dispatch(setNewPersonError('Пустая должность'))
            return
        }

        newPerson.unit = newPerson.unit.trim()
        if (newPerson.unit === '') {
            dispatch(setNewPersonError('Пустое подразделение'))
            return
        }

        axios.post(personsURL, newPerson).then((res) => {
            dispatch(addPerson(res.data))
            dispatch(setNewPerson(res.data))
            dispatch(uploadNewPersonPhotos())
        }).catch(err => {
            if (handleAuthError(dispatch, getState, err)) {
                return
            }
            dispatch(setNewPersonError(err.toString()))
        })
    }
}

export const uploadNewPersonPhotos = () => {
    return (dispatch, getState) => {
        const newPerson = getState().persons.newPerson
        const newPersonPhotos = getState().persons.newPersonPhotos

        let authErrorOccurred = false
        let errorOccurred = false

        let photosAdded = 0

        const requests = []

        newPersonPhotos.forEach(p => {
            if (p.loaded) {
                photosAdded++
                return
            }
            requests.push(axios.post(`${personsURL}/${newPerson.id}/faces`, p.file).then(() => {
                photosAdded++
                dispatch(setNewPersonPhotoLoaded(p))
            }).catch(err => {

                if (authErrorOccurred) {
                    dispatch(setNewPersonPhotoError(p, 'Отсутствовала авторизация'))
                    return
                } else if(handleAuthError(dispatch, getState, err)) {
                    authErrorOccurred = true
                }

                if (!errorOccurred) {
                    errorOccurred = true
                    dispatch(setNewPersonError('Не удалось загрузить все фото'))
                }

                if (err.response) {
                    switch (err.response.status) {
                        case 401:
                            dispatch(setNewPersonPhotoError(p, 'Отсутствовала авторизация'))
                            break
                        case 400:
                            dispatch(setNewPersonPhotoError(p, 'Не удалось распознать лицо'))
                            break
                        default:
                            dispatch(setNewPersonPhotoError(p, 'Неизвестная ошибка'))
                            console.error(err)
                            break
                    }
                }
            }))
        })

        axios.all(requests).then(() => {
            if (photosAdded === newPersonPhotos.length) {
                dispatch(cancelNewPerson())
            }
        })
    }
}

export const cancelNewPerson = () => ({
    type: CLEAR_NEW_PERSON
})

export const setNewPersonName = (name) => ({
    type: SET_NEW_PERSON_NAME,
    name
})

export const setNewPersonPosition = (position) => ({
    type: SET_NEW_PERSON_POSITION,
    position
})

export const setNewPersonUnit = (unit) => ({
    type: SET_NEW_PERSON_UNIT,
    unit
})

export const addNewPersonPhoto = (photoFile) => {
    return (dispatch) => {
        var reader = new FileReader()
        reader.onload = (e) => {
            dispatch({
                type: ADD_NEW_PERSON_PHOTO,
                photo: {
                    url: e.target.result,
                    file: photoFile,
                    loaded: false,
                    error: null
                }
            })
        }
        reader.readAsDataURL(photoFile)
    }
}

const setNewPersonPhotoLoaded = (photo) => ({
    type: SET_NEW_PERSON_PHOTO_LOADED,
    photo
})

export const removeNewPersonPhoto = (photo) => ({
    type: REMOVE_NEW_PERSON_PHOTO,
    photo
})

const setNewPersonPhotoError = (photo, error) => ({
    type: SET_NEW_PERSON_PHOTO_ERROR,
    photo, error
})

const setNewPersonError = (newPersonError) => ({
    type: SET_NEW_PERSON_ERROR,
    newPersonError
})

export const editPerson = (person) => {
    return (dispatch, getState) => {
        axios.get(`${personsURL}/${person.id}/faces`).then(res => {
            dispatch({ type: EDIT_PERSON, editPerson: person, editPersonFaces: res.data })
        }).catch(err => {
            handleAuthError(dispatch, getState, err)
        })
    }
}

export const saveEditPerson = () => {
    return (dispatch, getState) => {

        const editPerson = getState().persons.editPerson

        editPerson.name = editPerson.name.trim()
        if (editPerson.name === '') {
            dispatch(setEditPersonError('Пустое имя'))
            return
        }

        editPerson.position = editPerson.position.trim()
        if (editPerson.position === '') {
            dispatch(setEditPersonError('Пустая должность'))
            return
        }

        axios.put(personsURL, editPerson).then(() => {
            dispatch(syncEditPersonFacesAndPhotos())
        }).catch(err => {
            handleAuthError(dispatch, getState, err)
            dispatch(setEditPersonError('Не удалось сохранить персону'))
        })
    }
}

export const syncEditPersonFacesAndPhotos = () => {
    return (dispatch, getState) => {

        const editPerson = getState().persons.editPerson
        const editPersonFaces = getState().persons.editPersonFaces
        const editPersonPhotos = getState().persons.editPersonPhotos

        const facesToRemove = editPersonFaces.filter(f => f.toRemove)

        let facesRemoved = 0
        let photosAdded = 0

        let authErrorOccurred = false
        let errorOccurred = false

        const requests = []

        const errMsg = 'Не удалось синхронизировать все фотографии'

        facesToRemove.forEach((f) => {
            requests.push(axios.delete(`${personFacesURL}/${f.id}`).then(() => {
                facesRemoved++
                dispatch(removeEditPersonFace(f.id))
            }).catch(err => {

                if (authErrorOccurred) {
                    dispatch(setEditPersonFaceError(f.id, 'Отсутствовала авторизация'))
                    return
                } else if(handleAuthError(dispatch, getState, err)) {
                    authErrorOccurred = true
                }

                if (!errorOccurred) {
                    errorOccurred = true
                    dispatch(setEditPersonError(errMsg))
                }

                if (err.response) {
                    switch (err.response.status) {
                        case 401:
                            dispatch(setEditPersonFaceError(f.id, 'Отсутствовала авторизация'))
                            break
                        default:
                            dispatch(setEditPersonFaceError(f.id, 'Неизвестная ошибка'))
                            console.error(err)
                            break
                    }
                }

            }))
        })

        editPersonPhotos.forEach(p => {
            requests.push(axios.post(`${personsURL}/${editPerson.id}/faces`, p.file).then((res) => {
                photosAdded++
                dispatch(removeEditPersonPhoto(p))
                dispatch(addEditPersonFace(res.data))
            }).catch(err => {

                if (authErrorOccurred) {
                    dispatch(setEditPersonPhotoError(p, 'Отсутствовала авторизация'))
                    return
                } else if(handleAuthError(dispatch, getState, err)) {
                    authErrorOccurred = true
                }

                if (!errorOccurred) {
                    errorOccurred = true
                    dispatch(setEditPersonError(errMsg))
                }

                if (err.response) {
                    switch (err.response.status) {
                        case 401:
                            dispatch(setEditPersonPhotoError(p, 'Отсутствовала авторизация'))
                            break
                        case 400:
                            dispatch(setEditPersonPhotoError(p, 'Не удалось распознать лицо'))
                            break
                        default:
                            dispatch(setEditPersonPhotoError(p, 'Неизвестная ошибка'))
                            console.error(err)
                            break
                    }
                }
            }))
        })

        axios.all(requests).then(() => {
            if (facesRemoved === facesToRemove.length && photosAdded === editPersonPhotos.length) {
                dispatch(cancelEditPerson())
            }
        })

        // if (facesRemoved === facesToRemove.length && photosAdded === editPersonPhotos.length) {
        //     dispatch(cancelEditPerson())
        // }
    }
}

export const cancelEditPerson = () => ({
    type: CLEAR_EDIT_PERSON,
})

export const setEditPersonName = (name) => ({
    type: SET_EDIT_PERSON_NAME,
    name
})

export const setEditPersonPosition = (position) => ({
    type: SET_EDIT_PERSON_POSITION,
    position
})

export const setEditPersonUnit = (unit) => ({
    type: SET_EDIT_PERSON_UNIT,
    unit
})

export const setEditPersonFaceToRemove = (personFaceID) => ({
    type: SET_EDIT_PERSON_FACE_TO_REMOVE,
    personFaceID,
})

export const restoreEditPersonFace = (personFaceID) => ({
    type: RESTORE_EDIT_PERSON_FACE,
    personFaceID,
})

const removeEditPersonFace = (personFaceID) => ({
    type: REMOVE_EDIT_PERSON_FACE,
    personFaceID
})

const setEditPersonFaceError = (personFaceID, error) => ({
    type: SET_EDIT_PERSON_FACE_ERROR,
    personFaceID, error
})

const addEditPersonFace = (personFace) => ({
    type: ADD_EDIT_PERSON_FACE,
    personFace
})

export const addEditPersonPhoto = (photoFile) => {
    return (dispatch, getState) => {
        var reader = new FileReader()
        reader.onload = (e) => {
            dispatch({
                type: ADD_EDIT_PERSON_PHOTO,
                photo: {
                    url: e.target.result,
                    file: photoFile
                }
            })
        }
        reader.readAsDataURL(photoFile)
    }
}

export const removeEditPersonPhoto = (photo) => ({
    type: REMOVE_EDIT_PERSON_PHOTO,
    photo
})

const setEditPersonPhotoError = (photo, error) => ({
    type: SET_EDIT_PERSON_PHOTO_ERROR,
    photo, error
})

const setEditPersonError = (editPersonError) => ({
    type: SET_EDIT_PERSON_ERROR,
    editPersonError
})

export const removePerson = (id) => {
    return (dispatch, getState) => {
        axios.delete(`${personsURL}/${id}`).then(() => {
            dispatch(loadPersons())
        }).catch(err => {
            handleAuthError(dispatch, getState, err)
        })
    }
}