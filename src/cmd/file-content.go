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
import './.css';

class  extends Component {

  constructor(props) {
    super(props);
  }

  render() {
    return ();
  }
}

export default;
`)

const testContent = (
`import React from 'react';
import ReactDOM from 'react-dom';
import  from './';

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(< />, div);
  ReactDOM.unmountComponentAtNode(div);
});
`)