package cmd

const appComponentContent = (
`import React, { Component } from 'react';
import './App.css';

class App extends Component {

  render() {
    return (
      <div></div>
    );
  }
}

export default App;
`)

const componentContent =(
`import React, { Component } from 'react';
import './{$1}.css';

class {$1} extends Component {

  constructor(props) {
    super(props);
  }

  render() {
    return ();
  }
}

export default {$1};
`)

const testContent = (
`import React from 'react';
import ReactDOM from 'react-dom';
import {$1} from './{$1}';

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(< />, div);
  ReactDOM.unmountComponentAtNode(div);
});
`)