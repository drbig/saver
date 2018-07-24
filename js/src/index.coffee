require './index.less'
import React from 'react'
import {render} from 'react-dom'


class App extends React.Component
  render: ->
    <div>
      "Hello world!"
    </div>

render <App />, document.getElementById('app')
