package pages

import (
	"github.com/anirudhgray/balkan-assignment/frontend/src/components"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Catalog struct {
	app.Compo
}

func (c *Catalog) OnMount(ctx app.Context) {
	var token string
	ctx.LocalStorage().Get("token", &token)
	if token == "" {
		ctx.Navigate("/")
	}
}

func (c *Catalog) Render() app.UI {
	return app.Div().Class("background p-10 min-h-screen flex flex-col").Body(
		&components.Navbar{},
		&components.Title{TitleString: "Catalog"},
		app.Div().Body(
		//
		).Class("max-w-[80rem] xl:mx-auto"),
		&components.Footer{},
	)
}
