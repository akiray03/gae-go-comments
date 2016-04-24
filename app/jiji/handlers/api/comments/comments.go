package comments

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/mjibson/goon"
	"jiji/models"
	"encoding/json"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/appengine/log"
	"net/url"
	"strconv"
	"io/ioutil"
	"io"
	"time"
	"fmt"
	"os"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	g := goon.NewGoon(r)

	q := datastore.NewQuery(models.KIND_COMMENT).Order("-created_at").Limit(10)

	pages := make([]models.Comment, 0, 10)
	if _, err := g.GetAll(q, &pages); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(pages)
}

func Show(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	page, err := getComment(r, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(page)
}

func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	var comment models.Comment

	if err := json.Unmarshal(body, &comment); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	g := goon.NewGoon(r)
	if _, err := g.Put(&comment); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	type SlackMessage struct {
		Text string `json:"text"`
	}
	slack_message := SlackMessage{}
	slack_message.Text = fmt.Sprintf("New comment arrived (by *%s* )\n>%s", comment.Author, comment.Body)

	webhook_url := os.Getenv("SLACK_WEBHOOK")
	if webhook_url != "" {
		ctx := appengine.NewContext(r)
		client := urlfetch.Client(ctx)
		values := url.Values{}
		payload, err := json.Marshal(slack_message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		values.Add("payload", string(payload))

		resp, err := client.PostForm(webhook_url, values)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Infof(ctx, "%p %p", resp.Status, resp.Body)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(comment); err != nil {
		panic(err)
	}
}

func getComment(r *http.Request, params httprouter.Params) (models.Comment, error) {
	g := goon.NewGoon(r)
	comment := models.Comment{}

	id, err := strconv.ParseInt(params.ByName("id"), 10, 0)
	if err != nil {
		return comment, err
	}

	comment = models.Comment{Id: id}
	if err := g.Get(&comment); err != nil {
		return comment, err
	}

	return comment, nil
}
