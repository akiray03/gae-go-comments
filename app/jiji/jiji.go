package jiji

import (
	"net/http"
	"github.com/julienschmidt/httprouter"

	"jiji/handlers/comments"
	api_comments "jiji/handlers/api/comments"
	"jiji/handlers/static"
)

func init() {
	router := httprouter.New()
	router.GET("/comments", comments.Index)
	router.GET("/comments/:id", comments.ShowOrNew)
	router.GET("/comments/:id/edit", comments.Edit)
	router.POST("/comments", comments.Create)
	router.POST("/comments/:id", comments.Update)
	router.PUT("/comments/:id", comments.Update)

	router.GET("/api/comments", api_comments.Index)
	router.GET("/api/comments/:id", api_comments.Show)
	router.POST("/api/comments", api_comments.Create)

	router.GET("/", static.Index)

	http.Handle("/", router)
}
