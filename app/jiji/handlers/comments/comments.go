package comments

import (
	"html/template"
	"net/http"
	"github.com/julienschmidt/httprouter"

	"google.golang.org/appengine/datastore"
	"github.com/mjibson/goon"

	"jiji/models"
	"strconv"
	"time"
)

var (
	templateIndex = template.Must(template.ParseFiles("html/comments/index.html"))
	templateNew = template.Must(template.ParseFiles("html/comments/new.html"))
	templateEdit = template.Must(template.ParseFiles("html/comments/edit.html"))
	templateShow = template.Must(template.ParseFiles("html/comments/show.html"))
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	g := goon.NewGoon(r)

	q := datastore.NewQuery(models.KIND_COMMENT).Order("-created_at").Limit(10)

	pages := make([]models.Comment, 0, 10)
	if _, err := g.GetAll(q, &pages); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err := templateIndex.Execute(w, pages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ShowOrNew(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if params.ByName("id") == "new" {
		New(w, r, params)
	} else {
		Show(w, r, params)
	}
}

func Show(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	p, err := getComment(r, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := templateShow.Execute(w, p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Edit(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	p, err := getComment(r, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := templateEdit.Execute(w, p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func New(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := templateNew.Execute(w, r.FormValue("content"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	g := goon.NewGoon(r)
	comment := models.NewComment()
	comment.Author = r.FormValue("author")
	comment.Body = r.FormValue("body")

	_, err := g.Put(&comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/comments", http.StatusFound)
}

func Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	g := goon.NewGoon(r)
	comment, err := getComment(r, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	comment.Author = r.FormValue("author")
	comment.Body = r.FormValue("body")
	comment.UpdatedAt = time.Now()

	if _, err := g.Put(&comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/comments", http.StatusFound)
}

func getComment(r *http.Request, params httprouter.Params) (models.Comment, error) {
	g := goon.NewGoon(r)
	p := models.Comment{}

	id, err := strconv.ParseInt(params.ByName("id"), 10, 0)
	if err != nil {
		return p, err
	}

	p = models.Comment{Id: id}
	if err := g.Get(&p); err != nil {
		return p, err
	}

	return p, nil
}
