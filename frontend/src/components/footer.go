package components

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type Footer struct {
	app.Compo
}

func (f *Footer) Render() app.UI {
	return app.Div().Class("mx-auto mt-auto pt-8 w-full").Body(
		app.Div().Class("glass py-3 px-4 max-w-[80rem] w-full").Body(
			app.P().Text("Made by Anirudh Mishra").Class("text-center dark:text-purple-300"),
			app.P().Text("API at /api/v1").Class("mt-auto text-center dark:text-purple-300"),
		),
	)
}
