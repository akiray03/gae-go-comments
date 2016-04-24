package slack_notify

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	"os"
	"net/url"
	"encoding/json"
	"fmt"
	"jiji/models"
	"github.com/mjibson/goon"
	"strconv"
)

func Process(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(r)
	webhook_url := os.Getenv("SLACK_WEBHOOK")

	if webhook_url == "" {
		log.Infof(ctx, "SLACK_WEBHOOK not configured.")
		return
	}

	comment, err := getComment(r)
	if err != nil {
		log.Errorf(ctx, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Infof(ctx, "comment = %s", comment)

	type SlackMessage struct {
		Text string `json:"text"`
	}
	slack_message := SlackMessage{}
	slack_message.Text = fmt.Sprintf("New comment arrived (by *%s* )\n>%s", comment.Author, comment.Body)

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

func getComment(r *http.Request) (models.Comment, error) {
	g := goon.NewGoon(r)
	comment := models.Comment{}

	id, err := strconv.ParseInt(r.FormValue("comment_id"), 10, 0)
	if err != nil {
		return comment, err
	}

	comment = models.Comment{Id: id}
	if err := g.Get(&comment); err != nil {
		return comment, err
	}

	return comment, nil
}
