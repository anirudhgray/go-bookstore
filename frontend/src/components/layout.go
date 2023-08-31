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
	l.children = v
	return l
}

func (l *Layout) Title(v string) *Layout {
	l.title = v
	return l
}

func (l *Layout) Render() app.UI {
	return app.Div().Class("bg-gray-400 p-10 min-h-screen flex flex-col").Body(
		&Navbar{},
		app.H2().Text(l.title).Class("text-4xl font-bold text-purple-900 mt-6 mb-6 xl:text-center"),
		app.Div().Body(
			app.Range(l.children).Slice(func(i int) app.UI {
				return l.children[i]
			}),
		).Class("max-w-[80rem] xl:mx-auto"),
		&Footer{},
	)
}
