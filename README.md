# Paprika - Router

Paprika is a simple router implementation for [raylib]. It's designed to 
abstract the navigation stack for a UI project. It follows a simple `push` and
`pop` structure. The following is an example of a page:

```go
type HomePage struct {
	temp int
}

// This example doesn't need Mount and Unmount
func (h *HomePage) Mount(r *router.Router)   {}
func (h *HomePage) Unmount(r *router.Router) {}

// On each frame update check if the Left or Right keys are pressed and go to
// the next page or previous page accordingly
func (h *HomePage) Update(r *router.Router) {
	// Go to Next Page
	if rl.IsKeyPressed(rl.KeyRight) {
		r.Push(&HomePage{
			temp: h.temp + 1, // Set the parameters of the next page
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
```

Using that example HomePage, we can easily run it by creating a raylib window 
and running the router. This example also sets the history to 15 pages instead
of the default 10.

```go
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

```

Things to keep in mind is that the your Page will need to implment the Mount, 
Unmount, Update and Draw functions. Each have the purposes stated in the 
following Page interface:

```go
type Page interface {
	// Mount is called when the page is first pushed and set as the active page
	// to display. You can set up any resources you need here.
	Mount(r *Router)

	// Unmount is called when the page is popped and is no longer the active
	// page to display. You should clean up any resources you set up in Mount
	// here.
	Unmount(r *Router)

	// Update is called every frame before the page is Drawn. You should run
	// any logic here that changes the state of the page.
	Update(r *Router)

	// Draw is called every frame after the page is Updated. You should run any
	// logic here that draws the current state of the page.
	Draw(r *Router)
}
```


[raylib]: https://github.com/gen2brain/raylib-go