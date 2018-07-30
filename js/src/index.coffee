require './index.less'
import React from 'react'
import {render} from 'react-dom'


class Save extends React.Component
  constructor: (props) ->
    super props
    this.state = {
      isDetailed: false,
    }

  toggleDetailed: ->
    this.setState((prevState) => {isDetailed: !prevState.isDetailed})

  render: ->
    <li>
      <a className="knob">[del]</a>
      <a className="saveMain">{this.props.save.Stamp}</a>
      <span className="gameStamp">{this.props.save.Note}</span>
    </li>


class Game extends React.Component
  constructor: (props) ->
    super props
    this.state = {
      isExpanded: false,
      isDetailed: false,
    }

  toggleExpanded: ->
    this.setState((prevState) => {isExpanded: !prevState.isExpanded})

  toggleDetailed: ->
    this.setState((prevState) => {isDetailed: !prevState.isDetailed})

  render: ->
    knob = if this.state.isExpanded
      '[ - ]'
    else
      '[ + ]'

    <li className="game">
      <span className="knob">{knob}</span>
      <a className="info" onClick={=> this.toggleDetailed()}>[i]</a>
      <span className="savesCounter">({this.props.game.Saves.length})</span>
      <a className="gameName" onClick={=> this.toggleExpanded()}>{this.props.game.Name}</a>
      <span className="gameStamp">{this.props.game.Stamp}</span>
      {if this.state.isDetailed
        <div className="gameInfo">
          Path: {this.props.game.Path}
          Size: {this.props.game.Size}
          </div>
      }
      {if this.state.isExpanded
        <ol className="saves">
        {this.props.game.Saves.map((save) => <Save key={save.Stamp} save={save} />)}
        </ol>
      }
    </li>


class App extends React.Component
  constructor: (props) ->
    super props
    this.state = {
      isLoaded: false,
      error: null,
      cfg: {},
      currentGame: null,
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
      <div className="error">Error: {this.state.error.message}</div>
    else if !this.state.isLoaded
      <div>Loading...</div>
    else
      <div>
        Root: {this.state.cfg.Root}
        {if this.state.cfg.Games.length < 1
          <div className="Error">No games defined. Please use CLI.</div>
        else
          <ul className="games">
          {this.state.cfg.Games.map((game) => <Game key={game.Name} game={game} />)}
          </ul>
        }
      </div>

render <App />, document.getElementById('app')
