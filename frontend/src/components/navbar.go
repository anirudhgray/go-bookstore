package components

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type Navbar struct {
	app.Compo
}

func (n *Navbar) Render() app.UI {
	return app.Div().Body(
		app.Div().Body(
			app.A().Body(
				app.H1().Text("Secure Bookstore").Class("text-2xl font-bold"),
			).Href("/"),
			app.P().Text("Built using Golang and WASM"),
		),
		app.A().Class("bi bi-github text-4xl").Href("https://github.com/BalkanID-University/vit-2025-summer-engineering-internship-task-anirudhgray"),
	).Class("flex justify-between max-w-[80rem] w-full mx-auto")
}
