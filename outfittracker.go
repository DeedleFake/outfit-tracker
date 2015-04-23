package ot

import (
	"appengine"
	"appengine/urlfetch"
	"fmt"
	"net/http"
)

func handleMain(rw http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	client := urlfetch.Client(ctx)

	cq := &CensusQuery{
		Client: client,
		Object: "character",
		Query:  "c:resolve=outfit&c:show=character_id,times.last_login",
	}
	fmt.Fprintf(rw, "%v characters.", cq.Len())

	r, err := cq.Next()
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(rw, "%v", r)
}

func init() {
	http.HandleFunc("/", handleMain)
}
