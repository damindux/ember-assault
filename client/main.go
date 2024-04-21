package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	WINDOW_WIDTH  = 1280
	WINDOW_HEIGHT = 720
	MAX_BULLETS   = 30
)

type Bullet struct {
	Pos    rl.Vector2
	Vel    rl.Vector2
	Active bool
}

var gameOver bool = false
var activeBullets int = 0
var bullets [MAX_BULLETS]Bullet

func drawGame(mousePosition rl.Vector2, player *Player, bullets *[MAX_BULLETS]Bullet) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	topLeftPos := rl.Vector2{
		X: player.Pos.X - float32(player.Width)/2,
		Y: player.Pos.Y - float32(player.Height)/2,
	}

	for i := 0; i < MAX_BULLETS; i++ {
		if bullets[i].Active {
			rl.DrawCircleV(bullets[i].Pos, 5, rl.White)
		}
	}

	rl.DrawLine(int32(mousePosition.X-10), int32(mousePosition.Y), int32(mousePosition.X+10), int32(mousePosition.Y), rl.Green)
	rl.DrawLine(int32(mousePosition.X), int32(mousePosition.Y-10), int32(mousePosition.X), int32(mousePosition.Y+10), rl.Green)

	rl.DrawRectangleV(topLeftPos, rl.Vector2{X: 50, Y: 50}, rl.Red)

	rl.EndDrawing()
}

func updateGame(player *Player, bullets *[MAX_BULLETS]Bullet) rl.Vector2 {
	mousePosition := rl.GetMousePosition()
	direction := rl.Vector2Subtract(mousePosition, player.Pos)
	distance := rl.Vector2Length(direction)

	if distance > 0 {
		direction = rl.Vector2Scale(rl.Vector2Normalize(direction), 10)
	}

	if rl.IsKeyDown(rl.KeyW) {
		player.Pos.Y -= 5
	}
	if rl.IsKeyDown(rl.KeyS) {
		player.Pos.Y += 5
	}
	if rl.IsKeyDown(rl.KeyA) {
		player.Pos.X -= 5
	}
	if rl.IsKeyDown(rl.KeyD) {
		player.Pos.X += 5
	}

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) && activeBullets < MAX_BULLETS {
		for i := 0; i < MAX_BULLETS; i++ {
			if !bullets[i].Active {
				bullets[i].Active = true
				bullets[i].Pos = player.Pos
				bullets[i].Vel = direction
				activeBullets++
				break
			}
		}
	}

	for i := 0; i < MAX_BULLETS; i++ {
		if bullets[i].Active {
			bullets[i].Pos = rl.Vector2Add(bullets[i].Pos, bullets[i].Vel)
			if bullets[i].Pos.X < 0 || bullets[i].Pos.X > WINDOW_WIDTH || bullets[i].Pos.Y < 0 || bullets[i].Pos.Y > WINDOW_HEIGHT {
				bullets[i].Active = false
				activeBullets--
			}
		}
	}

	return mousePosition
}

func main() {
	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "Client")
	rl.SetTargetFPS(60)

	defer rl.CloseWindow()


    game := Game{
        player: &Player{
            Pos: rl.Vector2{X: 100, Y: 100},
            Width: 50,
            Height: 50,
            Health: 100,
        },
    }
    game.Init()

    defer game.Close()
    
	for !rl.WindowShouldClose() {
        if game.stateMachine.currentState != nil {
            game.stateMachine.currentState.OnExecute()
        }
	}
}
