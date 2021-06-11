package main

// Value JSON in request | response
type Bio struct {
	Name string `json:"name"`
	SecondName string `json:"secondName"`
}

// JSON for POST method request
type Request struct {
	Key uint32 `json:"key"`
	Val Bio `json:"value"`
}
