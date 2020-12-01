import React, { createRef, useEffect } from 'react'
import { connect } from 'react-redux'
import { Button, Col, Container, Row, Table, Modal, Form } from 'react-bootstrap'
import './index.css'
import {
    cancelEditPerson,
    cancelNewPerson,
    createNewPerson,
    editPerson,
    loadPersons,
    newPerson,
    removePerson,
    saveEditPerson,
    setEditPersonName,
    setEditPersonPosition,
    setNewPersonName,
    setNewPersonPosition,
    setNewPersonUnit,
    setEditPersonUnit,
    addNewPersonPhoto,
    removeNewPersonPhoto,
    uploadNewPersonPhotos,
    setEditPersonFaceToRemove,
    removeEditPersonPhoto,
    addEditPersonPhoto,
    restoreEditPersonFace
} from '../../reducers/persons'
import Dropzone from 'react-dropzone'
import { photosURL } from '../../reducers/skuder'
import { faCheck } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'

const NewPersonModal = (props) => {

    if (!props.newPerson) {
        return ''
    }

    let created = ''
    if (props.newPerson.id) {
        created = <div style={{color: 'green'}}>Персона создана</div>
    }

    let error = ''
    if (props.newPersonError) {
        error = <div style={{color: 'red'}}>{props.newPersonError}</div>
    }

    let allPhotosLoaded = !props.newPersonPhotos.some(p => !p.loaded)

    const dropzoneRef = createRef()
    const openPhotoLoader = () => {
        if (dropzoneRef.current) {
            dropzoneRef.current.open()
        }
    }

    return (
        <Modal id='person-modal' show={true} onHide={props.onCancelNewPerson} animation={false}>
            <Modal.Header closeButton>
                <Modal.Title>Новая персона</Modal.Title>
            </Modal.Header>
            <Modal.Body>
                <Form>
                    <Form.Group controlId='new-person-name'>
                        <Form.Label>Имя</Form.Label>
                        <Form.Control readOnly={props.newPerson.id} autoComplete='new-password' placeholder='Введите имя' value={props.newPerson.name} onChange={(e) => {
                            props.onNewPersonNameChange(e.currentTarget.value)
                        }}/>
                    </Form.Group>
                    <Form.Group controlId='new-person-position'>
                        <Form.Label>Должность</Form.Label>
                        <Form.Control readOnly={props.newPerson.id} autoComplete='new-password' placeholder='Введите должность' value={props.newPerson.position} onChange={(e) => {
                            props.onNewPersonPositionChange(e.currentTarget.value)
                        }}/>
                    </Form.Group>
                    <Form.Group controlId='new-person-unit'>
                        <Form.Label>Подразделение</Form.Label>
                        <Form.Control readOnly={props.newPerson.id} autoComplete='new-password' placeholder='Введите подразделение' value={props.newPerson.unit} onChange={(e) => {
                            props.onNewPersonUnitChange(e.currentTarget.value)
                        }}/>
                    </Form.Group>
                    <Form.Group>
                        {created}
                    </Form.Group>
                    <Form.Group>
                        <Form.Label>Фотографии</Form.Label>
                        <div className='photos'>
                            {props.newPersonPhotos.map((p, i) => (
                                <div key={i} className='photo'>
                                    <div className='image' style={{backgroundImage: 'url('+p.url+')'}}>
                                        {p.error
                                            ? <div className='error'>
                                                {p.error}
                                            </div>
                                            : ''
                                        }
                                    </div>
                                    {p.loaded
                                        ? <Button disabled={true} size='sm' variant='outline-success'><FontAwesomeIcon icon={faCheck}/></Button>
                                        : <Button size='sm' variant='outline-danger'
                                                  onClick={() => props.onNewPersonPhotoRemove(p)}>Удалить</Button>
                                    }
                                </div>
                            ))}
                            <Dropzone accept='image/jpeg' multiple={true} onDrop={files => files.forEach(f => props.onNewPersonPhotoAdd(f))} ref={dropzoneRef}>
                                {({getRootProps, getInputProps}) => (
                                    <div className='photo'>
                                        <div className='new-photo-dropzone' {...getRootProps()}>
                                            <div>Перетащите фотографии в эту область</div>
                                            <input {...getInputProps()} />
                                        </div>
                                        <Button size='sm' variant='outline-primary' onClick={openPhotoLoader}>Добавить фото</Button>
                                    </div>
                                )}
                            </Dropzone>
                        </div>
                    </Form.Group>
                    <Form.Group>
                        {error}
                    </Form.Group>
                </Form>
            </Modal.Body>
            {props.newPerson.id
                ? <Modal.Footer>
                    <Button variant='secondary'
                            onClick={props.onCancelNewPerson}>Отмена</Button>
                    {allPhotosLoaded
                        ? <Button variant='primary'
                                  onClick={props.onCancelNewPerson}>ОК</Button>
                        : <Button variant='primary'
                                  onClick={props.onNewPersonPhotosUpload}>Загрузить</Button>
                    }
                </Modal.Footer>
                : <Modal.Footer>
                    <Button variant='secondary'
                            onClick={props.onCancelNewPerson}>Отмена</Button>
                    <Button variant='primary'
                            onClick={props.onCreateNewPerson}>Создать</Button>
                </Modal.Footer>
            }
        </Modal>
    )
}

const EditPersonModal = (props) => {

    if (!props.editPerson) {
        return ''
    }

    let error = ''
    if (props.editPersonError) {
        error = <div style={{color: 'red'}}>{props.editPersonError}</div>
    }

    const dropzoneRef = createRef()
    const openPhotoLoader = () => {
        if (dropzoneRef.current) {
            dropzoneRef.current.open()
        }
    }

    return (
        <Modal id='person-modal' show={true} animation={false} onHide={props.onCancelEditPerson}>
            <Modal.Header closeButton>
                <Modal.Title>Персона №{props.editPerson.id}</Modal.Title>
            </Modal.Header>
            <Modal.Body>
                <Form>
                    <Form.Group controlId='edit-person-name'>
                        <Form.Label>Имя</Form.Label>
                        <Form.Control autoComplete='new-password' placeholder='Введите имя' value={props.editPerson.name} onChange={(e) => {
                            props.onEditPersonNameChange(e.currentTarget.value)
                        }}/>
                    </Form.Group>
                    <Form.Group controlId='edit-person-position'>
                        <Form.Label>Должность</Form.Label>
                        <Form.Control autoComplete='new-password' placeholder='Введите должность' value={props.editPerson.position} onChange={(e) => {
                            props.onEditPersonPositionChange(e.currentTarget.value)
                        }}/>
                    </Form.Group>
                    <Form.Group controlId='edit-person-unit'>
                        <Form.Label>Подразделение</Form.Label>
                        <Form.Control autoComplete='new-password' placeholder='Введите подразделение' value={props.editPerson.unit} onChange={(e) => {
                            props.onEditPersonUnitChange(e.currentTarget.value)
                        }}/>
                    </Form.Group>
                    <Form.Group>
                        <Form.Label>Фотографии</Form.Label>
                        <div className='photos'>
                            {props.editPersonFaces.map((f) => (
                                <div key={f.id} className='photo'>
                                    <div className='image' style={{backgroundImage: `url(${photosURL}/${f.photo_id})`}}>
                                        {f.error
                                            ? <div className='error'>
                                                {f.error}
                                            </div>
                                            : ''
                                        }
                                    </div>
                                    {f.toRemove
                                        ? <Button size='sm' variant='outline-success' onClick={() => props.onEditPersonFaceRestore(f.id)}>Восстановить</Button>
                                        : <Button size='sm' variant='outline-danger' onClick={() => props.onEditPersonFaceRemove(f.id)}>Удалить</Button>
                                    }
                                </div>
                            ))}
                            {props.editPersonPhotos.map((p, i) => (
                                <div key={'photo'+i} className='photo'>
                                    <div className='image' style={{backgroundImage: 'url('+p.url+')'}}>
                                        {p.error
                                            ? <div className='error'>
                                                {p.error}
                                            </div>
                                            : ''
                                        }
                                    </div>
                                    <Button size='sm' variant='outline-danger'
                                            onClick={() => props.onEditPersonPhotoRemove(p)}>Удалить</Button>
                                </div>
                            ))}
                            <Dropzone accept='image/jpeg' multiple={true} onDrop={files => files.forEach(f => props.onEditPersonPhotoAdd(f))} ref={dropzoneRef}>
                                {({getRootProps, getInputProps}) => (
                                    <div className='photo'>
                                        <div className='new-photo-dropzone' {...getRootProps()}>
                                            <div>Перетащите фотографию в эту область</div>
                                            <input {...getInputProps()} />
                                        </div>
                                        <Button size='sm' variant='outline-primary' onClick={openPhotoLoader}>Добавить фото</Button>
                                    </div>
                                )}
                            </Dropzone>
                        </div>
                    </Form.Group>
                    <Form.Group>
                        {error}
                    </Form.Group>
                </Form>
            </Modal.Body>
            <Modal.Footer>
                <Button variant='secondary' onClick={props.onCancelEditPerson}>Отмена</Button>
                <Button variant='primary' onClick={props.onSaveEditPerson}>Сохранить</Button>
            </Modal.Footer>
        </Modal>
    )
}

const Persons = ({ loadPersons, ...props }) => {

    useEffect(() => {
        loadPersons()
    }, [loadPersons])

    if (props.persons === null) {
        return ''
    }

    return (
        <Container id='persons-container'>
            <Row>
                <Col md='12'>
                    <Button onClick={props.onNewPerson}>Новая персона</Button>
                </Col>
            </Row>
            <Row>
                <Col md='12'>
                    <Table size='sm' striped bordered>
                        <thead>
                            <tr>
                                <th>ID</th>
                                <th>Имя</th>
                                <th>Должность</th>
                                <th>Подразделение</th>
                                <th className='manage-col'>Управление</th>
                            </tr>
                        </thead>
                        <tbody>
                            {props.persons.map(p =>
                                <tr key={p.id}>
                                    <td>{p.id}</td>
                                    <td>{p.name}</td>
                                    <td>{p.position}</td>
                                    <td>{p.unit}</td>
                                    <td className='manage-col'>
                                        <Button size='sm' variant='primary' onClick={() => props.onEditPerson(p)}>Редактировать</Button>
                                        <Button size='sm' variant='danger' onClick={() => props.onRemovePerson(p.id)}>Удалить</Button>
                                    </td>
                                </tr>
                            )}
                        </tbody>
                    </Table>
                </Col>
            </Row>
            <NewPersonModal {...props} />
            <EditPersonModal {...props} />
        </Container>
    )
}

const mapStateToProps = state => ({
    ...state.persons
})

const mapDispatchToProps = dispatch => {
    return {
        loadPersons: () => {
            dispatch(loadPersons())
        },
        onNewPerson: () => {
            dispatch(newPerson())
        },
        onCreateNewPerson: () => {
            dispatch(createNewPerson())
        },
        onCancelNewPerson: () => {
            dispatch(cancelNewPerson())
        },
        onNewPersonNameChange: (name) => {
            dispatch(setNewPersonName(name))
        },
        onNewPersonPositionChange: (position) => {
            dispatch(setNewPersonPosition(position))
        },
        onNewPersonUnitChange: (unit) => {
            dispatch(setNewPersonUnit(unit))
        },
        onNewPersonPhotoAdd: (photoFile) => {
            dispatch(addNewPersonPhoto(photoFile))
        },
        onNewPersonPhotoRemove: (photo) => {
            dispatch(removeNewPersonPhoto(photo))
        },
        onNewPersonPhotosUpload: () => {
            dispatch(uploadNewPersonPhotos())
        },
        onEditPerson: (person) => {
            dispatch(editPerson(person))
        },
        onSaveEditPerson: () => {
            dispatch(saveEditPerson())
        },
        onCancelEditPerson: () => {
            dispatch(cancelEditPerson())
        },
        onEditPersonNameChange: (name) => {
            dispatch(setEditPersonName(name))
        },
        onEditPersonPositionChange: (position) => {
            dispatch(setEditPersonPosition(position))
        },
        onEditPersonUnitChange: (unit) => {
            dispatch(setEditPersonUnit(unit))
        },
        onEditPersonFaceRemove: (personFaceID) => {
            dispatch(setEditPersonFaceToRemove(personFaceID))
        },
        onEditPersonFaceRestore: (personFaceID) => {
            dispatch(restoreEditPersonFace(personFaceID))
        },
        onEditPersonPhotoAdd: (photoFile) => {
            dispatch(addEditPersonPhoto(photoFile))
        },
        onEditPersonPhotoRemove: (photo) => {
            dispatch(removeEditPersonPhoto(photo))
        },
        onRemovePerson: (id) => {
            dispatch(removePerson(id))
        },
    }
}

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(Persons)