package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/bialas1993/go-notifier/db"
	"github.com/bialas1993/go-notifier/event"
	"github.com/bialas1993/go-notifier/schema"
	"github.com/bialas1993/go-notifier/search"
	"github.com/bialas1993/go-notifier/util"
)

func onNotifyCreated(m event.NotifyCreatedMessage) {
	notify := schema.Notify{
		ID:        m.ID,
		Title:     m.Title,
		Body:      m.Body,
		Service:   m.Service,
		CreatedAt: m.CreatedAt,
	}
	if err := search.Insert(context.Background(), notify); err != nil {
		log.Println(err)
	}
}

func listNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error

	// Read parameters
	skip := uint64(0)
	skipStr := r.FormValue("skip")
	take := uint64(100)
	takeStr := r.FormValue("take")
	if len(skipStr) != 0 {
		skip, err = strconv.ParseUint(skipStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid skip parameter")
			return
		}
	}
	if len(takeStr) != 0 {
		take, err = strconv.ParseUint(takeStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid take parameter")
			return
		}
	}

	notifications, err := db.List(ctx, skip, take)
	if err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, "Could not fetch notifications")
		return
	}

	util.ResponseOk(w, notifications)
}

func searchNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()

	// Read parameters
	query := r.FormValue("query")
	if len(query) == 0 {
		util.ResponseError(w, http.StatusBadRequest, "Missing query parameter")
		return
	}
	skip := uint64(0)
	skipStr := r.FormValue("skip")
	take := uint64(100)
	takeStr := r.FormValue("take")
	if len(skipStr) != 0 {
		skip, err = strconv.ParseUint(skipStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid skip parameter")
			return
		}
	}
	if len(takeStr) != 0 {
		take, err = strconv.ParseUint(takeStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid take parameter")
			return
		}
	}

	notifications, err := search.Search(ctx, query, skip, take)
	if err != nil {
		log.Println(err)
		util.ResponseOk(w, []schema.Notify{})
		return
	}

	util.ResponseOk(w, notifications)
}
