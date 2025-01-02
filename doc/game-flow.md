# Game flow

## Start options

| Mode | Description | New/Join | Public/Private | 
|------|-------------|----------|----------------|
| 1. Create new private game | Creates a new game that will not be public, so you must share the key. | New      | Private        |
| 2. Create new public game | Creates a new game that's going to be listed and open for anyone to join. | New      | Public         |
| 3. Join a private game | Join a game that's not listed, but that you received a key for. | Join     | Private        |
| 4. Join a public game | Browse the list of games and join one (this will fetch the list of games). | Join     | Public        |

> NOTE: probably want to also allow joining a _public_ game using a known key. 

> NOTE2: what about if I started a game and rejoin it later?


## Steps

| Step            | Description |
|-----------------|-------------|
| AskTheName      | get the name | 
| StartOrJoin     | start or join | 
| SelectGame      | choose an existing game key | 
| AskTheGameKey   | enter an existing game key | 
| ShowGameKey     | show the current game key | 
| StartGame       | start the game | 


## Pre-game flow

```plantuml
[*] --> AskTheName 
AskTheName : Ask the player's (unique) name

AskTheName --> StartOrJoin : confirm\nname
StartOrJoin : Show the user the 4 options listed 
StartOrJoin : under 'Start options', and make them 
StartOrJoin : choose one.

StartOrJoin --> AskTheName : back

StartOrJoin --> ShowGameKey : 1. new private game or\n 2. new public game
StartOrJoin --> AskTheGameKey : 3. join private game
StartOrJoin --> SelectGame : 4. join public game

state ShowGameKey {
    [*]-->GameCreated : create
    GameCreated: A game key is shown. 
    GameCreated: Waiting for another player
    GameCreated-->[*] : Another\nplayer joins

}

state AskTheGameKey {
    [*]-->WaitingForInput
    WaitingForInput-->ValidatingGameKey : Key entered
    ValidatingGameKey-->ShowError : Invalid game key
    ShowError-->WaitingForInput : Retry
    ValidatingGameKey-->[*] : Valid
}

state SelectGame {
    [*]-->ShowingGames
    ShowingGames-->GameSelected : select a game
    GameSelected-->[*]
}

SelectGame --> StartOrJoin : back
AskTheGameKey --> StartOrJoin : back
ShowGameKey --> StartOrJoin : back

state StartGame <<join>>

ShowGameKey --> StartGame
AskTheGameKey --> StartGame
SelectGame --> StartGame

StartGame --> StartTheGame

StartTheGame --> [*]



```
