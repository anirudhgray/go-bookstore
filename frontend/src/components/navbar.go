package components

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type Navbar struct {
	app.Compo

	auth bool
}

func (n *Navbar) OnMount(ctx app.Context) {
	var token string
	ctx.LocalStorage().Get("token", &token)
	if token != "" {
		n.auth = true
	} else {
		n.auth = false
	}
}

func (n *Navbar) logout(ctx app.Context, e app.Event) {
	ctx.LocalStorage().Del("token")
	ctx.LocalStorage().Del("user")
	ctx.Navigate("/")
}

func (n *Navbar) Render() app.UI {
	return app.Div().Body(
		app.Div().Body(
			app.A().Body(
				app.H1().Text("Secure Bookstore").Class("text-2xl font-bold"),
			).Href("/"),
			app.P().Text("Built using Golang and WASM"),
		),
		app.Div().Body(
			app.If(n.auth, app.P().Text("LOGOUT").OnClick(n.logout).Class("cursor-pointer hover:font-bold")),
			app.A().Class("bi bi-github text-4xl").Href("https://github.com/BalkanID-University/vit-2025-summer-engineering-internship-task-anirudhgray"),
		).Class("flex gap-5 items-center"),
	).Class("flex justify-between max-w-[80rem] w-full mx-auto")
}
