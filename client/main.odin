package main

import rb "vendor:raylib"
import "core:net"

Player :: struct {
    pos: rb.Vector2,
    width: i32,
    height: i32,
    health: i32,
}

Bullet :: struct {
    pos: rb.Vector2,
    vel: rb.Vector2,
    active: bool,
}

MAX_BULLETS :: 30
SCREEN_WIDTH :: 1280
SCREEN_HEIGHT :: 720

game_over := false
bullets := [MAX_BULLETS]Bullet{}
active_bullets := 0
p := Player{rb.Vector2{100, 100},50, 50, 100}

draw_game :: proc (mouse_position: rb.Vector2) {
    rb.BeginDrawing()
    rb.ClearBackground(rb.BLACK)
    top_left_pos: rb.Vector2 = { p.pos.x - f32(p.width / 2), p.pos.y - f32(p.height / 2) }

    for i := 0; i < MAX_BULLETS; i+=1 {
        if bullets[i].active {
            rb.DrawCircle(i32(bullets[i].pos.x), i32(bullets[i].pos.y), 5, rb.WHITE)
        }
    }

    rb.DrawLine(i32(mouse_position.x - 10), i32(mouse_position.y), i32(mouse_position.x + 10), i32(mouse_position.y), rb.GREEN);
    rb.DrawLine(i32(mouse_position.x), i32(mouse_position.y - 10), i32(mouse_position.x), i32(mouse_position.y + 10), rb.GREEN);
    rb.DrawRectangleV(top_left_pos, rb.Vector2{50, 50}, rb.RED)
    rb.EndDrawing()
}

update_game :: proc () -> rb.Vector2 {
    mouse_position := rb.GetMousePosition()
    direction := mouse_position - transmute(rb.Vector2)p.pos
    distance := rb.Vector2Length(direction)

    if distance > 0 {
        direction = rb.Vector2Normalize(direction)
    }

    if rb.IsKeyDown(rb.KeyboardKey.W) {
        p.pos.y -= 5
    }
    if rb.IsKeyDown(rb.KeyboardKey.S) {
        p.pos.y += 5
    }
    if rb.IsKeyDown(rb.KeyboardKey.A) {
        p.pos.x -= 5
    }
    if rb.IsKeyDown(rb.KeyboardKey.D) {
        p.pos.x += 5
    }

    if (rb.IsMouseButtonPressed(rb.MouseButton.LEFT) && (active_bullets < MAX_BULLETS)) {
        for i := 0; i < MAX_BULLETS; i+=1 {
            if !bullets[i].active {
                bullets[i].active = true
                bullets[i].pos = p.pos 
                bullets[i].vel = rb.Vector2Scale(direction, 10)
                break
            } 

        }

    }

    for i := 0; i < MAX_BULLETS; i+=1 {
        if bullets[i].active {
            bullets[i].pos.x += bullets[i].vel.x
            bullets[i].pos.y += bullets[i].vel.y

            if bullets[i].pos.y < 0 || bullets[i].pos.y > SCREEN_HEIGHT || bullets[i].pos.x < 0 || bullets[i].pos.x > SCREEN_WIDTH {
                bullets[i].active = false
            }
        }
    }
    return mouse_position
}

update_draw_game :: proc () {
    draw_game(update_game())
}

init_game :: proc () {
    for i := 0; i < MAX_BULLETS; i+=1 {
        bullets[i].active = false
    }
}

main :: proc () {
    rb.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "Client")
    rb.HideCursor()
    init_game()
    rb.SetTargetFPS(60)

    for !rb.WindowShouldClose() {
        update_draw_game()
    }

    rb.CloseWindow()
}
