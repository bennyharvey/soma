import React, { useEffect } from 'react'
import { connect } from 'react-redux'
import {
    Button,
    Col,
    Container,
    Form, FormControl,
    InputGroup,
    Row,
    Table
} from 'react-bootstrap'
import {
    loadEvents, loadPassageNames,
    setFrom, setPage,
    setPassageID,
    setPersonName, setTo
} from '../../reducers/events'
import './index.css'
import DateTimePicker from 'react-datetime-picker'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faArrowLeft, faArrowRight } from '@fortawesome/free-solid-svg-icons'
import { formatDateTime } from '../../utils'
import { photosURL } from '../../reducers/skuder'

const eventTypeNames = {
    face_recognize: 'Лицо распознано',
    person_recognize: 'Персона распознана',
    passage_open: 'Открытие прохода'
}

const Event = ({ type, passageNames, data }) => {
    switch (type) {
        case 'face_recognize':
            return <div className="event-info">
                <img className="event-info-photo" src={`${photosURL}/${data.photo_id}`} alt={data.photo_id}/>
                <div className="event-info-data">
                    <div className="event-info-row">Точность детектинга: {data.detect_confidence.toFixed(2)}</div>
                </div>
            </div>
        case 'person_recognize':
            return <div className="event-info">
                <img className="event-info-photo" src={`${photosURL}/${data.photo_id}`} alt={data.photo_id}/>
                <div className="event-info-data">
                    <div className="event-info-row">Имя персоны: {data.person_name}</div>
                    <div className="event-info-row">Точность детектинга: {data.detect_confidence.toFixed(2)}</div>
                    <div className="event-info-row">Расстояние совпадения: {data.descriptors_distance.toFixed(2)}</div>
                </div>
            </div>
        case 'passage_open':
            return <div className="event-info">
                <div className="event-info-data">
                    <div className="event-info-row">Название прохода: {passageNames[data.passage_id] || data.passage_id}</div>
                    <div className="event-info-row">Имя персоны: {data.person_name}</div>
                </div>
            </div>
        default:
            return <code>{JSON.stringify(data, null, 2)}</code>
    }
}

const Events = ({ loadEvents, loadPassageNames, ...props }) => {

    useEffect(() => {
        loadPassageNames()
        loadEvents()
    }, [loadEvents, loadPassageNames])

    if (props.events === null) {
        return ''
    }

    let passagesOptions = []
    for (const passageID in props.passageNames) {
        passagesOptions.push(
            <option key={passageID} value={passageID}>{props.passageNames[passageID]}</option>
        )
    }

    // const csvQueryString = queryString.stringify({
    //     from: props.from || undefined,
    //     to: props.to || undefined,
    //     passage_id: props.passageID || undefined,
    //     person_name: props.personName || undefined,
    //     order_by: props.orderBy || undefined,
    //     order_direction: props.orderDirection || undefined,
    //     format: 'csv'
    // })

    return (
        <div>
            <Row>
                <Col>
                    <Form.Row>
                        <Form.Group as={Col} xl='3' sm='6'>
                            <Form.Label>С</Form.Label>
                            <DateTimePicker
                                className='form-control'
                                format='yyyy-MM-dd HH:mm:ss'
                                calendarIcon={null}
                                yearPlaceholder='____'
                                monthPlaceholder='__'
                                dayPlaceholder='__'
                                hourPlaceholder='__'
                                minutePlaceholder='__'
                                secondPlaceholder='__'
                                renderNumbers={true}
                                renderSecondHand={true}
                                onChange={props.onFromChange}
                                locale='ru'
                                value={props.from} />
                        </Form.Group>
                        <Form.Group as={Col} xl='3' sm='6'>
                            <Form.Label>По</Form.Label>
                            <DateTimePicker
                                className='form-control'
                                format='yyyy-MM-dd HH:mm:ss'
                                calendarIcon={null}
                                yearPlaceholder='____'
                                monthPlaceholder='__'
                                dayPlaceholder='__'
                                hourPlaceholder='__'
                                minutePlaceholder='__'
                                secondPlaceholder='__'
                                renderNumbers={true}
                                renderSecondHand={true}
                                onChange={props.onToChange}
                                locale='ru'
                                value={props.to} />
                        </Form.Group>
                        <Form.Group as={Col} xl='3' sm='6'>
                            <Form.Label>Проход</Form.Label>
                            <Form.Control as='select' value={props.passage} onChange={(e) => props.onPassageIDChange(e.currentTarget.value || null)}>
                                <option></option>
                                {passagesOptions}
                            </Form.Control>
                        </Form.Group>
                        <Form.Group as={Col} xl='3' sm='6'>
                            <Form.Label>Имя</Form.Label>
                            <Form.Control value={props.personName} onChange={(e) => props.onPersonNameChange(e.currentTarget.value)}/>
                        </Form.Group>
                    </Form.Row>
                    {/*<Button as='a' href={eventsURL+'?'+csvQueryString} variant='primary'>*/}
                    {/*    Скачать CSV*/}
                    {/*</Button>*/}
                </Col>
            </Row>
            <Row>
                {props.events.length
                    ? <Col md='12'>
                        <Table size='sm' striped bordered>
                            <thead>
                            <tr>
                                <th>ID</th>
                                <th>Дата и время</th>
                                <th>Тип</th>
                            </tr>
                            </thead>
                            <tbody>
                                {props.events.map(e => [
                                    <tr key={e.id}>
                                        <td>{e.id}</td>
                                        <td>{formatDateTime(new Date(e.time))}</td>
                                        <td>{eventTypeNames[e.type] || e.type}</td>
                                    </tr>,
                                    <tr key={`${e.id}-data`}>
                                        <td colspan='3'>
                                            <Event passageNames={props.passageNames} type={e.type} data={e.data}/>
                                        </td>
                                    </tr>
                                ])}
                            </tbody>
                        </Table>
                    </Col>
                    : <Col md='12'><i>Нет записей</i></Col>
                }
                <Col md='12'>
                    <InputGroup className='pager'>
                        <InputGroup.Prepend>
                            {props.page > 1
                                ? <Button variant='outline-secondary' onClick={() => props.onPageChange(props.page-1)}><FontAwesomeIcon icon={faArrowLeft}/></Button>
                                : ''}
                            <InputGroup.Text as='label' htmlFor='events-page'>Страница</InputGroup.Text>
                        </InputGroup.Prepend>
                        <FormControl id='events-page' type='number' value={props.page} onChange={(e) => props.onPageChange(e.currentTarget.value)}/>
                        {props.events.length
                            ? <InputGroup.Append>
                                <Button variant='outline-secondary' onClick={() => props.onPageChange(props.page+1)}><FontAwesomeIcon icon={faArrowRight}/></Button>
                            </InputGroup.Append>
                            : ''}
                    </InputGroup>
                </Col>
            </Row>
        </div>
    )
}

const mapStateToProps = state => ({
    ...state.events,
})

const mapDispatchToProps = dispatch => {
    return {
        loadEvents: () => {
            dispatch(loadEvents())
        },
        loadPassageNames: () => {
            dispatch(loadPassageNames())
        },
        onFromChange: (from) => {
            dispatch(setFrom(from))
        },
        onToChange: (to) => {
            dispatch(setTo(to))
        },
        onPassageIDChange: (passageID) => {
            dispatch(setPassageID(passageID))
        },
        onPersonNameChange: (personName) => {
            dispatch(setPersonName(personName))
        },
        onPageChange: (page) => {
            dispatch(setPage(page))
        },
    }
}

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(Events)