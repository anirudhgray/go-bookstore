package pages

import (
	"io"
	"net/http"

	"github.com/anirudhgray/balkan-assignment/frontend/src/components"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Forgot struct {
	app.Compo

	email string

	err  string
	succ string

	loading bool
}

func (f *Forgot) OnMount(ctx app.Context) {
	var token string
	ctx.LocalStorage().Get("token", &token)
	if token != "" {
		ctx.Navigate("/catalog")
	}
}

func (f *Forgot) submit(ctx app.Context, e app.Event) {
	f.err = ""
	f.succ = ""
	f.loading = true

	e.PreventDefault()

	email := f.email

	req, err := http.NewRequest("GET", "/api/v1/auth/forgot-password", nil)
	if err != nil {
		f.err = err.Error()
		f.loading = false
		return
	}

	q := req.URL.Query()
	q.Add("email", email)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		f.err = err.Error()
		f.loading = false
		return
	}
	defer res.Body.Close()
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		f.err = err.Error()
		f.loading = false
		return
	}
	if res.StatusCode >= 400 {
		f.err = string(responseBody)
		f.loading = false
		return
	}
	f.err = ""
	f.succ = string(responseBody)
	f.loading = false
}

func (f *Forgot) Render() app.UI {
	return app.Div().Class("background md:p-10 p-5 min-h-screen flex flex-col").Body(
		&components.Navbar{},
		&components.Title{TitleString: "Forgot Password"},
		app.Div().Body(
			app.Div().Class("grid grid-cols-2").Body(
				app.Form().Class("xl:col-span-2 md:col-span-1 col-span-2 max-w-[30rem] xl:mx-auto").Body(
					app.Label().For("email").Text("Email"),
					app.Input().ID("email").Class("w-full mb-3 py-1 px-2 rounded-md").Type("email").Value(f.email).Placeholder("test@anrdhmshr.tech").OnChange(f.ValueTo(&f.email)),
					app.Button().Disabled(f.loading).Text("Send OTP").Class("px-3 py-2 bg-purple-500 hover:bg-purple-800 text-white rounded-md mt-6").OnClick(f.submit),
					app.P().Text(f.err).Class("text-red-900"),
					app.P().Text(f.succ).Class("text-green-900"),
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
