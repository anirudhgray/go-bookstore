package components

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Navbar struct {
	app.Compo

	auth bool
	dark string
}

func (n *Navbar) OnMount(ctx app.Context) {
	var token string
	ctx.LocalStorage().Get("token", &token)
	if token != "" {
		n.auth = true
	} else {
		n.auth = false
	}
	var theme string
	ctx.LocalStorage().Get("dark", &theme)
	n.dark = theme
	fmt.Println(n.dark, theme)
}

func (n *Navbar) logout(ctx app.Context, e app.Event) {
	ctx.LocalStorage().Del("token")
	ctx.LocalStorage().Del("user")
	ctx.Navigate("/")
}

func (n *Navbar) Render() app.UI {
	return app.Div().Body(
		app.Div().Body(
			app.A().Body(
				app.H1().Text("Secure Bookstore").Class("text-2xl font-bold dark:text-purple-200"),
			).Href("/"),
			app.P().Text("Built using Golang and WASM").Class("dark:text-purple-300"),
		),
		app.Div().Body(
			app.Div().ID("toggleLottieContainer").Class("sm:absolute sm:left-0 sm:-translate-x-full sm:-ml-3 grow basis-0").Body(
				app.If(n.dark == "true", app.Raw(`
					<lottie-player
						id="toggleLottie"
						src="/web/lottie/themeToggle.json"
						class="w-[4.5rem] cursor-pointer h-auto"
					>
					</lottie-player>
				`),
				).Else(app.Raw(`
				<lottie-player
					id="toggleLottie"
					src="/web/lottie/themeToggleInverse.json"
					class="w-[4.5rem] cursor-pointer h-auto"
				>
				</lottie-player>
			`)),
			),
			app.If(n.auth, app.Span().OnClick(n.logout).Class("bi dark:text-purple-200 bi-box-arrow-right cursor-pointer hover:font-bold text-4xl grow basis-0")),
			app.A().Class("bi bi-github text-4xl dark:text-purple-200 grow basis-0").Href("https://github.com/BalkanID-University/vit-2025-summer-engineering-internship-task-anirudhgray"),
		).Class("flex gap-5 items-center relative"),
		app.Script().Defer(true).Src("/web/js/themeToggle.js"),
	).Class("flex justify-between sm:flex-row items-center flex-col max-w-[80rem] w-full mx-auto glass py-3 px-4")
}
