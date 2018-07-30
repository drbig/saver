require './index.less'
import React from 'react'
import {render} from 'react-dom'


class App extends React.Component
  constructor: (props) ->
    super props
    this.state = {
      isLoaded: false,
      error: null,
      cfg: {},
    }

  componentDidMount: ->
    fetch('/api/list')
      .then((response) => response.json())
      .then(
        (ok) =>
          this.setState({
            isLoaded: true,
            cfg: ok,
          })
        ,
        (error) =>
          this.setState({
            isLoaded: true,
            error: error,
          })
      )

  render: ->
    if this.state.error
      <div class="error">Error: {this.state.error.message}</div>
    else if !this.state.isLoaded
      <div>Loading...</div>
    else
      <div>
        Root: {this.state.cfg.Root}
        {if this.state.cfg.Games.length < 1
          <div class="Error">No games defined. Please use CLI.</div>
        else
          <ul>
          {this.state.cfg.Games.map((game) =>
            <li>{game.Name}</li>
          )}
          </ul>
        }
      </div>

render <App />, document.getElementById('app')
