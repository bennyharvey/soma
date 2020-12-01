import React from 'react'
import { Route, Switch } from 'react-router-dom'
import { Navbar, Nav } from 'react-bootstrap'
import { connect } from 'react-redux'
import Login from '../login'
import { LinkContainer } from 'react-router-bootstrap'
import './index.css'
import { logout, ADMIN, SECURITY } from '../../reducers/skuder'
import Users from '../users'
import Persons from '../persons'
import Events from '../events'
import NotFound from './not-found'

const pages = [
    {
        menu: <LinkContainer to='/users' key='users'>
                <Nav.Link>Пользователи</Nav.Link>
            </LinkContainer>,
        route: <Route exact path='/users' component={Users} key='users'/>,
        roles: [ADMIN]
    },
    {
        menu: <LinkContainer to='/persons' key='persons'>
            <Nav.Link>Персоны</Nav.Link>
        </LinkContainer>,
        route: <Route exact path='/persons' component={Persons} key='persons'/>,
        roles: [ADMIN, SECURITY]
    },
    {
        menu: <LinkContainer to='/events' key='events'>
            <Nav.Link>События</Nav.Link>
        </LinkContainer>,
        route: <Route exact path='/events' component={Events} key='events'/>,
        roles: [ADMIN, SECURITY]
    }
]

const App = ({ user, onLogout }) => {
    let authorization
    let menu
    let routes

    if (user) {
        authorization = (
            <Nav activeKey='/'>
                <Nav.Link onClick={onLogout}>Выход</Nav.Link>
            </Nav>
        )

        const rolePages = pages.filter(mi => mi.roles.includes(user.role))

        menu = (
            <Nav activeKey='/' className='mr-auto'>
                {rolePages.map(rp => rp.menu)}
            </Nav>
        )

        routes = rolePages.map(rp => rp.route)
    } else {
        authorization = (
            <Nav activeKey='/'>
                <LinkContainer to='/login'>
                    <Nav.Link>Вход</Nav.Link>
                </LinkContainer>
            </Nav>
        )

        menu = (
            <Nav activeKey='/' className='mr-auto'>
            </Nav>
        )
    }

    return (<div id='app'>
        <Navbar bg='light' expand='lg'>
            <Navbar.Brand>Skuder</Navbar.Brand>
            <Navbar.Toggle aria-controls='basic-navbar-nav'/>
            <Navbar.Collapse id='basic-navbar-nav'>
                {menu}
                {authorization}
            </Navbar.Collapse>
        </Navbar>
        <Switch>
            {routes}
            <Route exact path='/login' component={Login}/>
            <Route render={(props) => <NotFound {...props} redirectOnRoot={routes ? routes[0].props.path : '/login'} />} />
        </Switch>
    </div>)
}

const mapStateToProps = ({ skuder }) => ({
    user: skuder.user
})

const mapDispatchToProps = dispatch => {
    return {
        onLogout: () => {
            dispatch(logout())
        },
    }
}

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(App)