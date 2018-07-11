package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/bialas1993/go-notifier/db"
	"github.com/bialas1993/go-notifier/event"
	"github.com/bialas1993/go-notifier/schema"
	"github.com/bialas1993/go-notifier/util"
)

func createNotifyHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		ID string `json:"id"`
	}

	ctx := r.Context()

	// Read parameters
	title := template.HTMLEscapeString(r.FormValue("title"))
	if len(title) < 1 {
		util.ResponseError(w, http.StatusBadRequest, "Invalid title")
		return
	}

	body := template.HTMLEscapeString(r.FormValue("body"))
	if len(body) < 1 {
		util.ResponseError(w, http.StatusBadRequest, "Invalid body")
		return
	}

	service := template.HTMLEscapeString(r.FormValue("service"))
	if len(service) < 1 {
		util.ResponseError(w, http.StatusBadRequest, "Invalid service name")
		return
	}

	// Create notify
	createdAt := time.Now().UTC()
	id, err := ksuid.NewRandomWithTime(createdAt)
	if err != nil {
		util.ResponseError(w, http.StatusInternalServerError, "Failed to create notify")
		return
	}
	notify := schema.Notify{
		ID:        id.String(),
		Title:     title,
		Body:      body,
		Service:   service,
		CreatedAt: createdAt,
	}
	if err := db.Insert(ctx, notify); err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, "Failed to create notify")
		return
	}

	// Publish event
	if err := event.PublishNotifyCreated(notify); err != nil {
		log.Println(err)
	}

	// Return new notify
	util.ResponseOk(w, response{ID: notify.ID})
}
