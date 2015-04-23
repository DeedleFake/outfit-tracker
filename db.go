package ot

import (
	"appengine"
	"appengine/datastore"
)

type Char struct {
	ID        string
	LastLogin time.Time
	OutfitID  string
}

func (c Char) key(ctx appengine.Context) *datastore.Key {
	return datastore.NewKey(ctx, "char", c.ID, nil)
}

func GetChar(ctx appengine.Context, id string) (c Char, err error) {
	c.ID = id
	err = datastore.Get(ctx, c.key(ctx), &c)
	return
}

func PutChar(ctx appengine.Context, c Char) error {
	return datastore.Put(ctx, c.key(ctx), c)
}

type Outfit struct {
	ID      int
	Name    string
	Tag     string
	Members int
}

func (o Outfit) key(ctx appengine.Context) *datastore.Key {
	return datastore.NewKey(ctx, "outfit", id, nil)
}

func GetOutfit(ctx appengine.Context, id string) (o Outfit, err error) {
	o.ID = id
	err = datastore.Get(ctx, o.key(ctx), &o)
	return
}

func PutOutfit(ctx appengine.Context, o Outfit) error {
	return datastore.Put(ctx, o.key(ctx), o)
}
