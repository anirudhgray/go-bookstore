package pages

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/anirudhgray/balkan-assignment/frontend/src/components"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type SetAfterForgot struct {
	app.Compo

	email       string
	otp         string
	newPassword string

	err  string
	succ string

	loading bool
}

func (s *SetAfterForgot) OnMount(ctx app.Context) {
	var token string
	ctx.LocalStorage().Get("token", &token)
	if token != "" {
		ctx.Navigate("/catalog")
	}
	s.email = app.Window().URL().Query()["email"][0]
	s.otp = app.Window().URL().Query()["otp"][0]
}

func (s *SetAfterForgot) submit(ctx app.Context, e app.Event) {
	s.err = ""
	s.succ = ""
	s.loading = true

	e.PreventDefault()

	values := map[string]string{"new_password": s.newPassword}
	jsonData, err := json.Marshal(values)
	if err != nil {
		s.err = err.Error()
		return
	}

	req, err := http.NewRequest("POST", "/api/v1/auth/set-forgotten-password", bytes.NewBuffer(jsonData))
	if err != nil {
		s.err = err.Error()
		s.loading = false
		return
	}

	q := req.URL.Query()
	q.Add("email", s.email)
	q.Add("otp", s.otp)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		s.err = err.Error()
		s.loading = false
		return
	}
	defer res.Body.Close()
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		s.err = err.Error()
		s.loading = false
		return
	}
	if res.StatusCode >= 400 {
		s.err = string(responseBody)
		s.loading = false
		return
	}
	s.err = ""
	s.succ = string(responseBody)
	s.loading = false
}

func (s *SetAfterForgot) Render() app.UI {
	return app.Div().Class("background p-10 min-h-screen flex flex-col").Body(
		&components.Navbar{},
		&components.Title{TitleString: "Set New Password"},
		app.Div().Body(
			app.Div().Class("grid grid-cols-2").Body(
				app.Form().Class("xl:col-span-2 md:col-span-1 col-span-2 max-w-[30rem] xl:mx-auto").Body(
					app.P().Text(s.email),
					// app.P().Text(s.otp),
					app.Label().For("password").Text("New Password"),
					app.Input().ID("password").Class("w-full mb-3 py-1 px-2 mt-2 rounded-md").Value(s.newPassword).Type("password").Placeholder("securePwd!0").OnChange(s.ValueTo(&s.newPassword)),
					app.Button().Disabled(s.loading).Text("Confirm").Class("px-3 py-2 bg-purple-500 hover:bg-purple-800 text-white rounded-md mt-6").OnClick(s.submit),
					app.P().Text(s.err).Class("text-red-900"),
					app.P().Text(s.succ).Class("text-green-900"),
					app.Span().Body(
						app.P().Text("Go to"),
						app.A().Text("Login.").Href("/login").Class("font-bold text-purple-600 hover:text-purple-800"),
					).Class("flex gap-1 mt-4"),
				),
			),
		).Class("max-w-[80rem] xl:mx-auto"),
		&components.Footer{},
	)
}
