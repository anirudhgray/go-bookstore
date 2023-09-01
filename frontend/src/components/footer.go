package components

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type Footer struct {
	app.Compo
}

func (f *Footer) Render() app.UI {
	return app.Div().Class("mt-auto glass py-3 px-4 max-w-[80rem] w-full mx-auto").Body(
		app.P().Text("Made by Anirudh Mishra").Class("text-center"),
		app.P().Text("API at /api/v1").Class("mt-auto text-center"),
	)
}
