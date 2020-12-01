import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './components/app';
import * as serviceWorker from './serviceWorker';
import { Provider } from 'react-redux'
import { ConnectedRouter } from 'connected-react-router'
import store from './store';
import { history } from './reducers'
import 'bootstrap/dist/css/bootstrap.min.css';
import * as axios from 'axios'

axios.defaults.withCredentials = true

ReactDOM.render(
    <Provider store={store}>
        <ConnectedRouter history={history}>
            <App />
        </ConnectedRouter>
    </Provider>,
    document.getElementById('root'));

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more login service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
