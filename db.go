package ot

import (
	"appengine"
	"appengine/datastore"
	"time"
)

type Char struct {
	ID        string
	LastLogin time.Time
	Server    string
	Faction   string
	Outfit    string
}

func (c Char) key(ctx appengine.Context) *datastore.Key {
	return datastore.NewKey(ctx, "char", c.ID, 0, nil)
}

func GetChar(ctx appengine.Context, id string) (c Char, err error) {
	c.ID = id
	err = datastore.Get(ctx, c.key(ctx), &c)
	if err == datastore.ErrNoSuchEntity {
		err = nil
		c.Outfit = "none"
	}

	return
}

func PutChar(ctx appengine.Context, c Char) error {
	_, err := datastore.Put(ctx, c.key(ctx), &c)
	return err
}

type Outfit struct {
	ID      string
	Name    string
	Tag     string
	Members int
	Server  string
	Faction string
}

func (o Outfit) key(ctx appengine.Context) *datastore.Key {
	return datastore.NewKey(ctx, "outfit", o.ID, 0, nil)
}

func GetOutfit(ctx appengine.Context, id string) (o Outfit, err error) {
	o.ID = id
	err = datastore.Get(ctx, o.key(ctx), &o)
	if err == datastore.ErrNoSuchEntity {
		err = nil
	}

	return
}

func PutOutfit(ctx appengine.Context, o Outfit) error {
	_, err := datastore.Put(ctx, o.key(ctx), &o)
	return err
}

type Movement struct {
	From, To string
	Amount   int
}

func (m Movement) key(ctx appengine.Context) *datastore.Key {
	return datastore.NewKey(ctx, "movement", m.From+"-"+m.To, 0, nil)
}

func GetMovement(ctx appengine.Context, from, to string) (m Movement, err error) {
	m.From = from
	m.To = to
	err = datastore.Get(ctx, m.key(ctx), &m)
	if err == datastore.ErrNoSuchEntity {
		err = nil
	}

	return
}

func PutMovement(ctx appengine.Context, m Movement) error {
	_, err := datastore.Put(ctx, m.key(ctx), &m)
	return err
}
