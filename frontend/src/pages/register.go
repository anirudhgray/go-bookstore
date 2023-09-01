package pages

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/anirudhgray/balkan-assignment/frontend/src/components"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Register struct {
	app.Compo

	email           string
	password        string
	name            string
	confirmPassword string

	err  string
	succ string

	loading bool
}

func (r *Register) OnMount(ctx app.Context) {
	var token string
	ctx.LocalStorage().Get("token", &token)
	if token != "" {
		ctx.Navigate("/catalog")
	}
}

func (r *Register) submit(ctx app.Context, e app.Event) {
	r.err = ""
	r.succ = ""
	r.loading = true

	e.PreventDefault()
	values := map[string]string{"email": r.email, "password": r.password, "name": r.name}
	jsonData, err := json.Marshal(values)
	if err != nil {
		r.err = err.Error()
		r.loading = false
		return
	}
	req, err := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	if err != nil {
		r.err = err.Error()
		r.loading = false
		return
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		r.err = err.Error()
		r.loading = false
		return
	}
	defer res.Body.Close()
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		r.err = err.Error()
		r.loading = false
		return
	}
	if res.StatusCode >= 400 {
		r.err = string(responseBody)
		r.loading = false
		return
	}
	r.err = ""
	r.succ = string(responseBody)
	r.loading = false
}

func (r *Register) resend(ctx app.Context, e app.Event) {
	r.err = ""
	r.succ = ""
	r.loading = true

	e.PreventDefault()

	email := r.email

	req, err := http.NewRequest("GET", "/api/v1/auth/request-verification", nil)
	if err != nil {
		r.err = err.Error()
		r.loading = false
		return
	}

	q := req.URL.Query()
	q.Add("email", email)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		r.err = err.Error()
		r.loading = false
		return
	}
	defer res.Body.Close()
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		r.err = err.Error()
		r.loading = false
		return
	}
	if res.StatusCode >= 400 {
		r.err = string(responseBody)
		r.loading = false
		return
	}
	r.err = ""
	r.succ = string(responseBody)
	r.loading = false
}

func (r *Register) Render() app.UI {
	return app.Div().Class("background md:p-10 p-5 min-h-screen flex flex-col").Body(
		&components.Navbar{},
		&components.Title{TitleString: "Register"},
		app.Div().Body(
			app.Div().Class("grid grid-cols-2").Body(
				app.Form().Class("xl:col-span-2 md:col-span-1 col-span-2 max-w-[30rem] xl:mx-auto").Body(
					app.Label().For("name").Text("Name"),
					app.Input().ID("name").Class("w-full mb-3 py-1 px-2 rounded-md").Value(r.name).Placeholder("John Doe").OnChange(r.ValueTo(&r.name)),
					app.Label().For("email").Text("Email"),
					app.Input().ID("email").Class("w-full mb-3 py-1 px-2 rounded-md").Type("email").Value(r.email).Placeholder("test@anrdhmshr.tech").OnChange(r.ValueTo(&r.email)),
					app.Label().For("password").Text("Password"),
					app.Input().ID("password").Class("w-full mb-3 py-1 px-2 rounded-md").Value(r.password).Type("password").Placeholder("securePwd!0").OnChange(r.ValueTo(&r.password)),
					app.Label().For("repeatpassword").Text("Repeat Password"),
					app.Input().ID("repeatpassword").Class("w-full mb-3 py-1 px-2 rounded-md").Value(r.confirmPassword).Type("password").Placeholder("securePwd!0").OnChange(r.ValueTo(&r.confirmPassword)),
					app.Button().Disabled(r.loading).Text("Register").Class("px-3 py-2 bg-purple-500 hover:bg-purple-800 text-white rounded-md mt-6 dark:bg-purple-600 dark:hover:bg-purple-400").OnClick(r.submit),
					app.P().Text(r.err).Class("text-red-900"),
					app.P().Text(r.succ).Class("text-green-900"),
					app.If(r.succ != "", app.P().Text("Resend verification mail.").Class("font-bold text-purple-600 hover:text-purple-800 dark:text-purple-500 dark:hover:text-purple-300").OnClick(r.resend)),

					app.Span().Body(
						app.P().Text("Have an existing account?"),
						app.A().Text("Log in.").Href("/login").Class("font-bold text-purple-600 hover:text-purple-800 dark:text-purple-500 dark:hover:text-purple-300"),
					).Class("flex gap-1 mt-4"),
				),
			),
		).Class("max-w-[80rem] xl:mx-auto"),
		&components.Footer{},
	)
}
