package main

const componentContent = (
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

const tsComponentContent = (
`import * as React from 'react';
import './{$1}.css';

interface Props {}

interface State {}

class {$1} extends React.Component<Props, State> {

  constructor(props: Props) {
    super(props)
    this.state = {};
  }

  render() {
    return ();
  }
}

export default {$1};
`)

const sfcComponentContent = (
`import * as React from 'react';
import './{$1}.css';

const {$1} = props => {
  return ();
}

export default {$1};
`)

const tsSfcComponentContent = (
`import * as React from 'react';
import './{$1}.css';

interface Props {}

const {$1}: React.SFC<Props> = props => {
  return ();
}

export default {$1};
`)

const testContent = (
`import React from 'react';
import ReactDOM from 'react-dom';
import {$1} from './{$1}';

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(<{$1} />, div);
  ReactDOM.unmountComponentAtNode(div);
});
`)