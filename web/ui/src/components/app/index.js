import React from 'react'
import { BrowserRouter, Route, Switch, Link } from 'react-router-dom'
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
import 'semantic-ui-css/semantic.min.css'
import { 
    Button,
    Checkbox,
    Grid,
    Header,
    Icon,
    Image,
    Menu,
    Segment,
    Sidebar, 
} from 'semantic-ui-react'

const ButtonExampleButton = () => (
    <div>
    <Button animated>
      <Button.Content visible>Next</Button.Content>
      <Button.Content hidden>
        <Icon name='arrow right' />
      </Button.Content>
    </Button>
    <Button animated='vertical'>
      <Button.Content hidden>Shop</Button.Content>
      <Button.Content visible>
        <Icon name='shop' />
      </Button.Content>
    </Button>
  </div>
)


const pages = [
    {
        menu: <LinkContainer to='/users' key='users'>
                <Nav.Link>Пользователи системы</Nav.Link>
            </LinkContainer>,
        route: <Route exact path='/users' component={Users} key='users'/>,
        roles: [ADMIN]
    },
    {
        menu: <LinkContainer to='/persons' key='persons'>
            <Nav.Link>Досье</Nav.Link>
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
    const [sidebarVisible, setSidebarVisible] = React.useState(true)

    if (user) {
        authorization = (
            <Nav activeKey='/'>
                <Nav.Link onClick={onLogout}>Выход</Nav.Link>
            </Nav>
        )

        const rolePages = pages.filter(mi => mi.roles.includes(user.role))

        menu = (
            <Sidebar
                as={Menu}
                animation='push'
                // icon='labeled'
                inverted
                //   onHide={() => setVisible(false)}
                vertical
                visible={sidebarVisible}
                width='wide'
                className="soma-sidebar"
            >
                <Menu.Item as={Link} to='/users' >
                    <Icon name='home' />
                    Пользователи
                </Menu.Item>
                <Menu.Item as={Link} to='/persons' >
                    <Icon name='gamepad' />
                    Досье
                </Menu.Item>
                <Menu.Item as={Link} to='/events' >
                    <Icon name='camera' />
                    События
                </Menu.Item>
            </Sidebar>
        )

        routes = rolePages.map(rp => rp.route)

        routes = (
            <div id='app'>
                <Route path="/users" component={Users} />
                <Route path="/persons" component={Persons} />
                <Route path="/events" component={Events} />
            </div>
        )
    } else {
        authorization = (
            <Nav activeKey='/'>
                <LinkContainer to='/login'>
                    <Nav.Link>Вход</Nav.Link>
                </LinkContainer>
            </Nav>
        )

        menu = (
            <Sidebar>
            </Sidebar>
        )

        routes = (
            <div id='app'>
                <Route path="/login" component={Login} />
            </div>
        )
    }

    return (
        <div>
        
        <div className="soma-wrapper">
        <BrowserRouter>
            
           {menu}

            <Sidebar.Pusher>
                <Navbar bg='light' expand='lg'>
                    <Button secondary icon onClick={(e, data) => setSidebarVisible(sidebarVisible ? false : true)}>
                        <Icon name='align justify' />
                    </Button>
                    <Navbar.Brand>TTK SOMA UI</Navbar.Brand>
                    <Navbar.Toggle aria-controls='basic-navbar-nav'/>
                    <Navbar.Collapse id='basic-navbar-nav'>
                        {authorization}
                    </Navbar.Collapse>
                </Navbar>
               
               {routes}
               
            </Sidebar.Pusher>
        </BrowserRouter>
        </div>
        </div>
    )
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