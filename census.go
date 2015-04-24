package ot

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	CensusURL = "https://census.daybreakgames.com"
	CensusID  = "outfittracker"
)

var (
	Servers = map[string]string{
		"19": "Jaeger",
		"25": "Briggs",
		"13": "Cobalt",
		"10": "Miller",
		"1":  "Connery",
		"17": "Emerald",
	}

	Factions = map[string]string{
		"0": "Nanite Systems",
		"1": "Vanu Sovereignty",
		"2": "New Conglomerate",
		"3": "Terran Republic",
	}
)

type CensusCountResult struct {
	Count int `json:"count"`
}

type CensusGetResult map[string]interface{}

func (r CensusGetResult) Returned() int {
	return int(r["returned"].(float64))
}

func (r CensusGetResult) List(name string) []interface{} {
	return r[name+"_list"].([]interface{})
}

func (r CensusGetResult) Err() error {
	if err, ok := r["error"]; ok {
		return CensusError{
			Err: err.(string),
		}
	}

	return nil
}

type CensusError struct {
	Err string `json:"error"`
}

func (err CensusError) Error() string {
	return err.Err
}

type CensusQuery struct {
	Client *http.Client

	Object string
	Query  string
	Num    int

	at     int
	length int
}

func (cq *CensusQuery) getURL() string {
	return fmt.Sprintf("%v/s:%v/get/ps2/%v?%v&c:limitPerDB=%v&c:start=%v",
		CensusURL,
		CensusID,
		cq.Object,
		cq.Query,
		cq.Num,
		cq.at,
	)
}

func (cq *CensusQuery) countURL() string {
	return fmt.Sprintf("%v/s:%v/count/ps2/%v?%v",
		CensusURL,
		CensusID,
		cq.Object,
		cq.Query,
	)
}

func (cq *CensusQuery) init() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	if cq.length == 0 {
		rsp, err := cq.Client.Get(cq.countURL())
		if err != nil {
			panic(err)
		}
		defer rsp.Body.Close()

		dec := json.NewDecoder(rsp.Body)

		var r CensusCountResult
		err = dec.Decode(&r)
		if err != nil {
			panic(err)
		}

		cq.length = r.Count
	}

	if cq.Num == 0 {
		cq.Num = 20
	}

	return nil
}

func (cq *CensusQuery) At() int {
	return cq.at
}

func (cq *CensusQuery) Len() int {
	cq.init()

	return cq.length
}

func (cq *CensusQuery) HasNext() bool {
	return cq.at < cq.Len()
}

func (cq *CensusQuery) Next() (r CensusGetResult, err error) {
	cq.init()

	if !cq.HasNext() {
		err = errors.New("No more data to get")
		return
	}

	rsp, err := cq.Client.Get(cq.getURL())
	if err != nil {
		return
	}
	defer rsp.Body.Close()

	d := json.NewDecoder(rsp.Body)

	err = d.Decode(&r)
	if err != nil {
		return
	}

	cq.at += cq.Num

	return
}
