package ot

import (
	"appengine"
	"appengine/datastore"
	"appengine/urlfetch"
	"encoding/json"
	"fmt"
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
		Query:  "times.last_login=]1429747200&c:resolve=outfit,world&c:show=character_id,faction_id,times.last_login",
	}

	switch initial {
	case true:
		log.Printf("Populating with %v characters...", cq.Len())
	case false:
		log.Printf("Updating %v characters...", cq.Len())
	}

	for cq.HasNext() {
		if p := cq.At() * 100 / cq.Len(); p%10 == 0 {
			log.Printf("Finished %v%% of update.", p)
		}

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

			var moved bool
			if !initial {
				yc, err := GetChar(ctx, c.ID)
				if err != nil {
					panic(err)
				}

				if yc.Outfit != c.Outfit {
					moved = true

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

					log.Printf("%v moved from %q to %q", c.ID, oo.Name, ro["name"])
				}
			}

			if initial || moved {
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
			}

			err = PutChar(ctx, c)
			if err != nil {
				panic(err)
			}
		}
	}

	switch initial {
	case true:
		log.Printf("Population complete.")
	case false:
		log.Printf("Update complete.")
	}
}

func handleView(rw http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	server := q.Get("server")
	faction := q.Get("faction")

	err := tmpl.ExecuteTemplate(rw, "view", map[string]interface{}{
		"server":    Servers[server],
		"faction":   Factions[faction],
		"serverID":  server,
		"factionID": faction,
	})
	if err != nil {
		panic(err)
	}
}

func handleData(rw http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	server := q.Get("server")
	faction := q.Get("faction")

	ctx := appengine.NewContext(req)

	oq := datastore.NewQuery("outfit")
	oq = oq.Filter("Server =", server)
	oq = oq.Filter("Faction =", faction)

	mq := datastore.NewQuery("movement")

	num, err := oq.Count(ctx)
	if err != nil {
		panic(err)
	}

	nodes := make([]map[string]interface{}, 0, num)
	edges := make([]map[string]interface{}, 0, num)

	var ocur Outfit
	oiter := oq.Run(ctx)
	for {
		_, err := oiter.Next(&ocur)
		if err != nil {
			if err == datastore.Done {
				break
			}

			panic(err)
		}

		label := fmt.Sprintf("%v (%v)", ocur.Name, ocur.Members)
		if ocur.Tag != "" {
			label = fmt.Sprintf("[%v] %v", ocur.Tag, label)
		}

		nodes = append(nodes, map[string]interface{}{
			"data": map[string]interface{}{
				"id":    ocur.ID,
				"label": label,
			},
		})

		var mcur Movement
		miter := mq.Filter("From =", ocur.ID).Run(ctx)
		for {
			mkey, err := miter.Next(&mcur)
			if err != nil {
				if err == datastore.Done {
					break
				}

				panic(err)
			}

			edges = append(edges, map[string]interface{}{
				"data": map[string]interface{}{
					"id":     mkey.StringID(),
					"weight": mcur.Amount,
					"source": mcur.From,
					"target": mcur.To,
				},
			})
		}
	}

	e := json.NewEncoder(rw)
	err = e.Encode(map[string]interface{}{
		"nodes": nodes,
		"edges": edges,
	})
	if err != nil {
		panic(err)
	}
}

func handleMain(rw http.ResponseWriter, req *http.Request) {
	err := tmpl.ExecuteTemplate(rw, "main", nil)
	if err != nil {
		panic(err)
	}
}

func init() {
	http.HandleFunc("/update", handleUpdate)
	http.HandleFunc("/view", handleView)
	http.HandleFunc("/data", handleData)
	http.HandleFunc("/", handleMain)
}
