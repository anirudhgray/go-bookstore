package components

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type Title struct {
	app.Compo

	TitleString string
}

func (t *Title) Render() app.UI {
	return app.H2().Text(t.TitleString).Class("text-4xl font-bold text-purple-900 mt-6 mb-6 xl:text-center dark:text-purple-100")
}
