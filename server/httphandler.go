// Package server handles web requests for webcig.
package server

import (
	"github.com/graphql-go/handler"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

// app is the heart of this package. It routes requests to handlers and does logging.
type app struct {
	router    *httprouter.Router // HTTP request routing
	staticDir string             // Location of general file resources (js, templates, images)
	templates *template.Template // all our templates
	log       *log.Logger        // General logging
}

// User holds the web app-relevant info for a user.
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// ServeHTTP fulfill's app's obligation to the Handler interface.
func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

// PlainPage contains barebones data for an HTML template.
type Page struct {
	Title string
}

// renderTemplate executes the specified template, code 500 if error
func (a *app) renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := a.templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		a.log.Printf("%s\n", err.Error())
		http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
	}
}

// webcigHomeRoute serves up the webcig home page.
func (a *app) webcigHomeRoute(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	page := Page{Title: "webcig Home"}
	a.renderTemplate(w, "base", &page)
}

// graphiqlRoute serves up the graphiql interface.
func (a *app) graphiqlRoute(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	page := Page{Title: "webcig graphql query interface"}
	a.renderTemplate(w, "graphiql", &page)
}

// graphqlRoute is a closure that passes a GraphQL http.Handler into an httprouter handler function.
func (a *app) graphqlRoute(graphqlHandler http.Handler) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		graphqlHandler.ServeHTTP(w, r)
	}
}

// New creates a new webcig app and returns it as an http.Handler.
func New(templateDir string, staticDir string) http.Handler {
	logger := log.New(os.Stdout, "web ", log.LstdFlags)

	router := httprouter.New()
	log := logger
	templates := template.Must(template.ParseFiles(templateDir+"base.html", templateDir+"graphiql.html"))
	app := &app{router, staticDir, templates, log}

	graphqlHandler := handler.New(&handler.Config{
		Schema: &Schema,
		Pretty: true,
	})

	router.GET("/", app.webcigHomeRoute)
	router.GET("/graphiql/", app.graphiqlRoute)
	router.POST("/graphql/", app.graphqlRoute(graphqlHandler))
	router.ServeFiles("/js/*filepath", http.Dir(staticDir+"/js/"))
	router.ServeFiles("/css/*filepath", http.Dir(staticDir+"/css/"))

	return app
}
