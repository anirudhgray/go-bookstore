package components

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type Layout struct {
	app.Compo

	children []app.UI
	title    string
}

func NewPage() *Layout {
	return &Layout{}
}

func (l *Layout) Children(v ...app.UI) *Layout {
	l.children = app.FilterUIElems(v...)
	return l
}

func (l *Layout) Title(v string) *Layout {
	l.title = v
	return l
}

func (l *Layout) Render() app.UI {
	return app.Div().Class("bg-gray-400 p-10 min-h-screen flex flex-col").Body(
		app.Div().Body(
			app.Div().Body(
				app.H1().Text("Secure Bookstore").Class("text-2xl font-bold"),
				app.P().Text("Built using Golang and WASM"),
			),
			app.A().Class("bi bi-github text-4xl").Href("https://github.com/BalkanID-University/vit-2025-summer-engineering-internship-task-anirudhgray"),
		).Class("flex justify-between"),
		app.H2().Text(l.title).Class("text-4xl font-bold text-purple-900 mt-6 mb-2"),
		app.Range(l.children).Slice(func(i int) app.UI {
			return l.children[i]
		}),
		app.P().Text("Made by Anirudh Mishra").Class("mt-auto pt-6 text-center"),
	)
}
