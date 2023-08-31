package main

import (
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// hello is a component that displays a simple "Hello World!". A component is a
// customizable, independent, and reusable UI element. It is created by
// embedding app.Compo into a struct.
type hello struct {
	app.Compo
}

// The Render method is where the component appearance is defined. Here, a
// "Hello World!" is displayed as a heading.
func (h *hello) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("Hello World! Mm. Coolio. NEEEEATO.").Class("text"),
		&hello2{},
	)
}

type hello2 struct {
	app.Compo
}

// The Render method is where the component appearance is defined. Here, a
// "Hello World!" is displayed as a heading.
func (h *hello2) Render() app.UI {
	return app.H1().Text("Hello again. Hahaha. Noish. Update? Woohoo. Dayumn.")
}

// The main function is the entry point where the app is configured and started.
// It is executed in 2 different environments: A client (the web browser) and a
// server.
func main() {
	// The first thing to do is to associate the hello component with a path.
	//
	// This is done by calling the Route() function,  which tells go-app what
	// component to display for a given path, on both client and server-side.
	app.Route("/", &hello{})
	app.Route("/nice", &hello2{})

	// Once the routes set up, the next thing to do is to either launch the app
	// or the server that serves the app.
	//
	// When executed on the client-side, the RunWhenOnBrowser() function
	// launches the app,  starting a loop that listens for app events and
	// executes client instructions. Since it is a blocking call, the code below
	// it will never be executed.
	//
	// When executed on the server-side, RunWhenOnBrowser() does nothing, which
	// lets room for server implementation without the need for precompiling
	// instructions.
	app.RunWhenOnBrowser()

	h := app.Handler{
		Name:        "Hello",
		Description: "An Hello World! example",
		Resources:   app.LocalDir("/"),
		Styles: []string{
			"/web/styles.css",
			"https://cdn.tailwindcss.com",
		},
	}

	// Finally, launching the server that serves the app is done by using the Go
	// standard HTTP package.
	//
	// The Handler is an HTTP handler that serves the client and all its
	// required resources to make it work into a web browser. Here it is
	// configured to handle requests with a path that starts with "/".
	http.Handle("/", &h)

	http.Handle("/nice", &h)

	err := app.GenerateStaticWebsite(".", &h)

	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":8001", nil); err != nil {
		log.Fatal(err)
	}
}
