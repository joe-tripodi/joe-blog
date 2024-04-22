package blogserver

import (
	"database/sql"
	"html/template"
	"joe-blog/database"
	"log"
	"net/http"
	"regexp"
)

type ServerConfig struct {
	TemplatePrefix string
	Templates      *template.Template
	Database       *sql.DB
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
  allPosts, err := database.GetAllPosts(config.Database)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
	renderTemplate(w, "index", allPosts)
}

func readHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[len("/read/"):]
	log.Printf("User requested: %s", url)

	post, err := database.GetPost(config.Database, url)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	renderTemplate(w, "read", *post)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p any) {
	err := config.Templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var config = ServerConfig{
	TemplatePrefix: "/templates/",
	Templates:      nil,
	Database:       nil,
}

var templates *template.Template
var validPath = regexp.MustCompile("^/(read)/([a-zA-Z0-9]+)$")

func Server(tp string, db *sql.DB) {
	log.SetPrefix("blog-server: ")
	log.SetFlags(0)

	config.TemplatePrefix = tp
	config.Templates = template.Must(template.ParseFiles(
		config.TemplatePrefix+"read.html",
		config.TemplatePrefix+"index.html"),
	)

	config.Database = db

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/read/", readHandler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
