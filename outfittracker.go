package ot

import (
	"appengine"
	"appengine/urlfetch"
	"log"
	"net/http"
	"strconv"
	"time"
)

func handleUpdate(rw http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)

	q := req.URL.Query()
	initial := q.Get("initial") != ""

	client := urlfetch.Client(ctx)
	cq := &CensusQuery{
		Client: client,
		Object: "character",
		Query:  "c:resolve=outfit,world&c:show=character_id,faction_id,times.last_login",
	}

	log.Printf("Updating %v characters...", cq.Len())

	for cq.HasNext() {
		r, err := cq.Next()
		if err != nil {
			panic(err)
		}

		if rerr := r.Err(); rerr != nil {
			panic(rerr)
		}

		for _, rci := range r.List("character") {
			rc := rci.(map[string]interface{})

			rll := rc["times"].(map[string]interface{})["last_login"].(string)
			ill, _ := strconv.ParseInt(rll, 0, 10)
			ll := time.Unix(ill, 0)

			var ro map[string]interface{}
			if r, ok := rc["outfit"]; ok {
				ro = r.(map[string]interface{})
			}

			c := Char{
				ID:        rc["character_id"].(string),
				LastLogin: ll,
				Server:    rc["world_id"].(string),
				Faction:   rc["faction_id"].(string),
				Outfit:    "none",
			}
			if ro != nil {
				c.Outfit = ro["outfit_id"].(string)
			}
			log.Printf("Updating char: %v", c.ID)

			yc, err := GetChar(ctx, c.ID)
			if err != nil {
				panic(err)
			}

			if !initial && (yc.Outfit != c.Outfit) {
				m, err := GetMovement(ctx, yc.Outfit, c.Outfit)
				if err != nil {
					panic(err)
				}

				m.Amount++

				err = PutMovement(ctx, m)
				if err != nil {
					panic(err)
				}

				oo, err := GetOutfit(ctx, yc.Outfit)
				if err != nil {
					panic(err)
				}

				oo.Members--

				err = PutOutfit(ctx, oo)
				if err != nil {
					panic(err)
				}

				log.Printf("Movement from %q to %q", oo.Name, ro["name"])
			}

			no, err := GetOutfit(ctx, c.Outfit)
			if err != nil {
				panic(err)
			}

			if no.ID != "none" {
				no.Name = ro["name"].(string)
				no.Tag = ro["alias"].(string)
			}
			no.Server = c.Server
			no.Faction = c.Faction
			no.Members++

			err = PutOutfit(ctx, no)
			if err != nil {
				panic(err)
			}

			err = PutChar(ctx, c)
			if err != nil {
				panic(err)
			}
		}
	}

	log.Printf("Update complete.")
}

func init() {
	http.HandleFunc("/", handleUpdate)
}
