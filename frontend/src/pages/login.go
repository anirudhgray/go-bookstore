package pages

import (
	"github.com/anirudhgray/balkan-assignment/frontend/src/components"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Login struct {
	app.Compo
}

func (h *Login) Render() app.UI {
	return components.NewPage().Title("Login").Children()
}
