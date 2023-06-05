package main

import (
	"fmt"

	"github.com/lcox74/paprika/router"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type HomePage struct {
	temp int
}

func (h *HomePage) Mount(r *router.Router)   {}
func (h *HomePage) Unmount(r *router.Router) {}
func (h *HomePage) Update(r *router.Router) {
	// Go to Next Page
	if rl.IsKeyPressed(rl.KeyRight) {
		r.Push(&HomePage{
			temp: h.temp + 1,
		})
	}

	// Go to Previous Page
	if rl.IsKeyPressed(rl.KeyLeft) {
		r.Pop()
	}
}
func (h *HomePage) Draw(r *router.Router) {
	// Clear the screen
	rl.ClearBackground(rl.RayWhite)
	rl.DrawText(fmt.Sprintf("This is page %d", h.temp), 190, 200, 20, rl.LightGray)
}

func main() {
	rl.InitWindow(800, 450, "raylib [core] example - basic window")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	r := router.NewRouter(
		router.WithHistory(15),
		router.WithDefaultPage(&HomePage{
			temp: 1,
		}),
	)

	r.Run()
}
