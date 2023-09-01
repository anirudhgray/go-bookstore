package pages

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/anirudhgray/balkan-assignment/frontend/src/components"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type BookType struct {
	Name   string
	Author string
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
}

func (c *Catalog) OnMount(ctx app.Context) {
	var token string
	ctx.LocalStorage().Get("token", &token)
	if token == "" {
		ctx.Navigate("/")
	}
	c.getBooks(ctx)
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
			app.P().Text(c.err).Class("text-red-900"),
			app.Range(c.books).Slice(func(i int) app.UI {
				return app.Div().Body(
					app.H3().Text(c.books[i].Book.Name),
				)
			}),
		).Class("max-w-[80rem] xl:mx-auto"),
		&components.Footer{},
	)
}
