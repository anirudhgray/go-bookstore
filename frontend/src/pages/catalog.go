package pages

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/anirudhgray/balkan-assignment/frontend/src/components"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type BookType struct {
	Name      string
	Author    string
	ISBN      string
	Price     int64
	Publisher string
}

type CatalogItem struct {
	Book      BookType
	AvgRating float64
}

type Response struct {
	Books []CatalogItem
}

type Catalog struct {
	app.Compo

	books []CatalogItem
	err   string

	search string
}

func (c *Catalog) OnMount(ctx app.Context) {
	var token string
	ctx.LocalStorage().Get("token", &token)
	if token == "" {
		ctx.Navigate("/")
	}
	c.getBooks(ctx)

	searchElem := app.Window().GetElementByID("catalogSearch")
	app.Window().Get("Tagify").New(searchElem, map[string]interface{}{
		"mode":    "mix",
		"pattern": "#",
	})
}

func (c *Catalog) getBooks(ctx app.Context) {
	req, err := http.NewRequest("GET", "/api/v1/books/catalog", nil)
	if err != nil {
		c.err = err.Error()
		return
	}
	var token string
	err = ctx.LocalStorage().Get("token", &token)

	req.Header.Set("Authorization", "Bearer "+token)

	if err != nil {
		c.err = err.Error()
		return
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		c.err = err.Error()
		return
	}
	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		c.err = err.Error()
		return
	}
	if res.StatusCode >= 400 {
		c.err = string(responseBody)
		return
	}

	var response = new(Response)
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		c.err = err.Error()
		return
	}

	c.err = ""

	c.books = response.Books
}

func (c *Catalog) Render() app.UI {
	return app.Div().Class("background md:p-10 p-5 min-h-screen flex flex-col").Body(
		&components.Navbar{},
		&components.Title{TitleString: "Catalog"},
		app.Div().Body(
			app.Form().Class("w-full relative mb-6 flex flex-col gap-4").Body(
				app.Textarea().ID("catalogSearch").Class("resize-none w-full py-1 px-2 rounded-md bg-white sm:pr-20").OnChange(c.ValueTo(&c.search)).Placeholder("Search for books!").Text(c.search),
				app.Button().Type("submit").Class("absolute sm:block hidden right-0 top-0 px-4 bi bi-search py-[0.65rem] bg-purple-500 hover:bg-purple-800 text-white rounded-md dark:bg-purple-600 dark:hover:bg-purple-400"),
				app.Button().Text("Search").Type("submit").Class("block sm:hidden px-4 bi bi-search py-3 bg-purple-500 hover:bg-purple-800 text-white rounded-md dark:bg-purple-600 dark:hover:bg-purple-400 ml-auto"),
			),
			app.P().Text(c.err).Class("text-red-900"),
			app.Div().Body(
				app.Range(c.books).Slice(func(i int) app.UI {
					return app.Div().Body(
						app.P().Text(c.books[i].Book.Author),
						app.H3().Text(c.books[i].Book.Name).Class("text-xl font-bold dark:text-purple-200 text-purple-800"),
						app.P().Text("ISBN: "+c.books[i].Book.ISBN).Class("mt-2 text-sm"),
						app.P().Text(c.books[i].Book.Publisher).Class("mt-2 text-sm"),
						app.P().Text(c.books[i].Book.Price).Class("text-right mt-2 text-lg dark:text-purple-200 text-purple-800"),
					).Class("col-span-1 glass-catalog p-3")
				}),
			).Class("grid gap-6 lg:grid-cols-4 md:grid-cols-3 sm:grid-cols-2 grid-cols-1 w-full"),
		).Class("max-w-[80rem] xl:mx-auto"),
		&components.Footer{},
	)
}
