package pages

import (
	"io"
	"net/http"

	"github.com/anirudhgray/balkan-assignment/frontend/src/components"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Verify struct {
	app.Compo

	email string
	otp   string

	err  string
	succ string

	loading bool
}

func (v *Verify) OnMount(ctx app.Context) {
	var token string
	ctx.LocalStorage().Get("token", &token)
	if token != "" {
		ctx.Navigate("/catalog")
	}
	v.email = app.Window().URL().Query()["email"][0]
	v.otp = app.Window().URL().Query()["otp"][0]
}

func (v *Verify) submit(ctx app.Context, e app.Event) {
	v.err = ""
	v.succ = ""
	v.loading = true

	e.PreventDefault()

	req, err := http.NewRequest("POST", "/api/v1/auth/verify", nil)
	if err != nil {
		v.err = err.Error()
		v.loading = false
		return
	}

	q := req.URL.Query()
	q.Add("email", v.email)
	q.Add("otp", v.otp)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		v.err = err.Error()
		v.loading = false
		return
	}
	defer res.Body.Close()
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		v.err = err.Error()
		v.loading = false
		return
	}
	if res.StatusCode >= 400 {
		v.err = string(responseBody)
		v.loading = false
		return
	}
	v.err = ""
	v.succ = string(responseBody)
	v.loading = false
}

func (v *Verify) Render() app.UI {
	return app.Div().Class("background md:p-10 p-5 min-h-screen flex flex-col").Body(
		&components.Navbar{},
		&components.Title{TitleString: "Account Verification"},
		app.Div().Body(
			app.Div().Class("grid grid-cols-2").Body(
				app.Form().Class("xl:col-span-2 md:col-span-1 col-span-2 max-w-[30rem] xl:mx-auto").Body(
					app.P().Text(v.email),
					// app.P().Text(v.otp),
					app.Button().Disabled(v.loading).Text("Verify Me").Class("px-3 py-2 bg-purple-500 hover:bg-purple-800 text-white rounded-md mt-6").OnClick(v.submit),
					app.P().Text(v.err).Class("text-red-900"),
					app.P().Text(v.succ).Class("text-green-900"),
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
