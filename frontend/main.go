package main

import (
	"log"
	"net/http"

	"github.com/anirudhgray/balkan-assignment/frontend/src/pages"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// The main function is the entry point where the app is configured and started.
// It is executed in 2 different environments: A client (the web browser) and a
// server.
func main() {
	// The first thing to do is to associate the hello component with a path.
	//
	// This is done by calling the Route() function,  which tells go-app what
	// component to display for a given path, on both client and server-side.
	app.Route("/", &pages.Landing{})
	app.Route("/login", &pages.Login{})
	app.Route("/catalog", &pages.Catalog{})

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
			"https://cdn.jsdelivr.net/npm/bootstrap-icons@1.10.5/font/bootstrap-icons.css",
		},
		Scripts: []string{
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
