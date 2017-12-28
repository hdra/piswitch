import * as React from 'react';
import * as ReactDOM from 'react-dom';
import * as ReactModal from 'react-modal';
import App from './App';
import registerServiceWorker from './registerServiceWorker';
import './index.css';

ReactModal.setAppElement(document.getElementById('root') as HTMLElement);

ReactDOM.render(
  <App />,
  document.getElementById('root') as HTMLElement
);
registerServiceWorker();
