package components

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type Footer struct {
	app.Compo
}

func (f *Footer) Render() app.UI {
	return app.P().Text("Made by Anirudh Mishra").Class("mt-auto pt-6 text-center")
}
