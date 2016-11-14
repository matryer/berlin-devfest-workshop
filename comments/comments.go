package comments

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

const entityKindComment = "Comment"

type commentEntity struct {
	Key     *datastore.Key `json:"id" datastore:"-"`
	Body    string         `json:"body" datastore:",noindex"`
	Created time.Time      `json:"ctime"`
	Author  string         `json:"author" datastore:",noindex"`
	Tags    []string       `json:"tags"`
}

func commentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		commentCreateHandler(w, r)
	case "GET":
		log.Println("GET BABY")
		commentsGetHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}

func commentsGetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	var comments []*commentEntity
	keys, err := datastore.
		NewQuery(entityKindComment).
		Order("-Created").
		Limit(25).
		GetAll(ctx, &comments)
	if err != nil {
		http.Error(w, fmt.Sprintf("GetAll: %s", err), 500)
		return
	}
	for i, key := range keys {
		comments[i].Key = key
	}
	err = json.NewEncoder(w).Encode(comments)
	if err != nil {
		http.Error(w, fmt.Sprintf("Encode: %s", err), 500)
		return
	}
}

func commentCreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	var comment commentEntity
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if comment.Author == "/dev/null" {
		http.Error(w, "PLEASE STOP", 500)
		return
	}
	comment.Created = time.Now()
	err = saveComment(ctx, &comment)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = json.NewEncoder(w).Encode(comment)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func saveComment(ctx context.Context, c *commentEntity) error {
	key := datastore.NewIncompleteKey(ctx, entityKindComment, nil)
	var err error
	key, err = datastore.Put(ctx, key, c)
	if err != nil {
		return err
	}
	c.Key = key
	return nil
}
