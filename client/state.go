package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
    "fmt"
)

type State interface {
	OnEnter()
	OnExecute()
	OnExit()
}

type Player struct {
    Pos                   rl.Vector2
    Width, Height, Health int
}

type Game struct {
    player       *Player
    stateMachine *StateMachine
    network      *NetworkClient
}

// ---- State Machine ----

type StateMachine struct {
	currentState State
    previousState State
	nextState    State
}

func (sm *StateMachine) shouldBeSendingData() bool {
    _, isPlaying := sm.currentState.(*PlayingState)
    return isPlaying
}

// ---- Change State ----

func (sm *StateMachine) ChangeState(newState State) {
    if sm.currentState != nil {
        sm.currentState.OnExit()
    }

    sm.previousState = sm.currentState
    sm.currentState = newState
    sm.currentState.OnEnter()
}

// ---- Playing State ----

type PlayingState struct {
	game *Game
}

func (p *PlayingState) OnEnter() {
    rl.HideCursor()
    if _, isGameOver := p.game.stateMachine.previousState.(*GameOverState); isGameOver {
        p.game.Init()
    }
}

func (p *PlayingState) OnExecute() {
    handleGamePlay(p.game.player, &bullets)

	if rl.IsKeyPressed(rl.KeyP) { // Press P to pause
		p.game.stateMachine.ChangeState(&PausedState{game: p.game})
	}
	if p.game.player.Health <= 0 {
		p.game.stateMachine.ChangeState(&GameOverState{game: p.game})
	}
}

func (p *PlayingState) OnExit() {
	// Clean up the playing state
}

// ---- Paused State ----

type PausedState struct {
	game *Game
}

func (p *PausedState) OnEnter() {
	// Paused state enter logic
}

func (p *PausedState) OnExecute() {
	if rl.IsKeyPressed(rl.KeyP) { // Press P to resume
		p.game.stateMachine.ChangeState(&PlayingState{game: p.game})
	}
}

func (p *PausedState) OnExit() {
	// Clean up the paused state
}

// ---- Game Over State ----

type GameOverState struct {
	game *Game
}

func (g *GameOverState) OnEnter() {
    rl.ShowCursor()
    rl.DrawText("Game Over", 100, 100, 20, rl.Red)
}

func (g *GameOverState) OnExecute() {
	if rl.IsKeyPressed(rl.KeyEnter) { 
		g.game.stateMachine.ChangeState(&PlayingState{game: g.game})
	}
}

func (g *GameOverState) OnExit() {
	// Clean up game over state
}

func (g *Game) Init() {
    g.player.Pos = rl.Vector2{X: 100, Y: 100}
    g.player.Height = 100

    var err error
    g.network, err = InitNetworkClient("127.0.0.1", 10000, g.player)
    if err != nil {
        fmt.Println("Error initializing network client:", err)
        return
    }

    go g.network.SendPlayerData()
    go g.network.ReceiveUpdates()

    for i := range bullets {
        bullets[i].Active = false
    }
    activeBullets = 0
    
    g.stateMachine.ChangeState(&PlayingState{game: g})
}

// TODO
func (g *Game) updateGameStateBasedOnNetworkData() {
    _, isPlaying := g.stateMachine.(*PlayingState) {
    if !isPlaying {
        return
    }
}

func (g *Game) Close() {
    g.network.conn.Close()
}

func handleGamePlay(player *Player, bullets *[MAX_BULLETS]Bullet) {
    mousePosition := updateGame(player, bullets)
    drawGame(mousePosition, player, bullets)
}
