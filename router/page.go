package router

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

func mountPage(r *Router, page Page) {
	if page == nil {
		return
	}
	page.Mount(r)
}
func unmountPage(r *Router, page Page) {
	if page == nil {
		return
	}
	page.Unmount(r)
}
func updatePage(r *Router, page Page) {
	if page == nil {
		return
	}
	page.Update(r)
}
func drawPage(r *Router, page Page) {
	if page == nil {
		return
	}
	page.Draw(r)
}
