package pages

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/anirudhgray/balkan-assignment/frontend/src/components"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Login struct {
	app.Compo

	email    string
	password string
	err      string
}

func (l *Login) OnMount(ctx app.Context) {
	var token string
	ctx.LocalStorage().Get("token", &token)
	if token != "" {
		ctx.Navigate("/catalog")
	}
	ctx.ObserveState("error").Value(&l.err)
}

func (l *Login) submit(ctx app.Context, e app.Event) {
	e.PreventDefault()
	values := map[string]string{"email": l.email, "password": l.password}
	jsonData, err := json.Marshal(values)
	if err != nil {
		ctx.SetState("error", err.Error())
		return
	}
	req, err := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
	if err != nil {
		ctx.SetState("error", err.Error())
		return
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		ctx.SetState("error", err.Error())
		return
	}
	defer res.Body.Close()
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		ctx.SetState("error", err.Error())
		return
	}
	if res.StatusCode >= 400 {
		ctx.SetState("error", string(responseBody))
		return
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(responseBody, &m)
	if err != nil {
		ctx.SetState("error", err.Error())
		return
	}
	ctx.SetState("error", "")

	ctx.LocalStorage().Set("token", m["token"])
	ctx.LocalStorage().Set("user", m["user"])

	ctx.Navigate("/catalog")
}

func (l *Login) Render() app.UI {
	return app.Div().Class("bg-gray-400 p-10 min-h-screen flex flex-col").Body(
		&components.Navbar{},
		&components.Title{TitleString: "Login"},
		app.Div().Body(
			app.Div().Class("grid grid-cols-2").Body(
				app.Form().Class("xl:col-span-2 md:col-span-1 col-span-2 max-w-[30rem] xl:mx-auto").Body(
					app.Label().For("email").Text("Email"),
					app.Input().ID("email").Class("w-full mb-3 py-1 px-2 rounded-md").Value(l.email).Placeholder("test@anrdhmshr.tech").OnChange(l.ValueTo(&l.email)),
					app.Label().For("password").Text("Password"),
					app.Input().ID("password").Class("w-full mb-3 py-1 px-2 rounded-md").Value(l.password).Placeholder("securePwd!0").OnChange(l.ValueTo(&l.password)),
					app.Button().Text("Login").Class("px-3 py-2 bg-purple-500 hover:bg-purple-800 text-white rounded-md mt-6").OnClick(l.submit),
					app.P().Text(l.err).Class("text-red-900 "+l.err),
				),
				// app.P().Text("Nice").Class("md:col-span-1 col-span-2"),
			),
		).Class("max-w-[80rem] xl:mx-auto"),
		&components.Footer{},
	)
}
