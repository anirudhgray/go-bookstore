package pages

import (
	"github.com/anirudhgray/balkan-assignment/frontend/src/components"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// A component is a
// customizable, independent, and reusable UI element. It is created by
// embedding app.Compo into a struct.
type Landing struct {
	app.Compo
}

// The Render method is where the component appearance is defined. Here, a
// "Hello World!" is displayed as a heading.
func (h *Landing) Render() app.UI {
	return components.NewPage().Title("Welcome").Children(
		app.P().Text("The backend API is COMPLETE and you can see the docs in the github readme (base url is at /api). Backend is written in Golang (using Gin and Gorm, though more on why the latter was a bad idea later), with Postgres as a DB."),
		app.P().Text("The frontend is also pretty cool, since instead of a JS framework, I'm using Go compiled to WebAssembly. However, at the time of this writing, the frontend consists of basically what you see in front of you. Fun exercise nonetheless.").Class("mt-2 text-justify"),
		app.Img().Alt("20-minute-adventure").Src("/web/images/rickandmorty.jpeg").Class("mt-6 md:w-2/3 md:max-w-[40rem] mx-auto w-full"),
		app.Span().Body(
			app.P().Text("Let's go, in and out. 20 minute adventure.").Class("italic text-sm"),
			app.P().Text("- me, looking cheerfully at go and wasm earlier.").Class("text-sm"),
		).Class("text-center mb-6"),
		app.A().Text("Enter Bookstore").Class("px-3 py-2 bg-purple-500 hover:bg-purple-800 text-white rounded-md max-w-full block w-[20rem] text-center mt-6 mx-auto").Href("/login"),
	)
}
