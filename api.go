package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tarantool/go-tarantool"
	"net/http"
	"strconv"
	"time"
)

// Name of the space in Tarantool
var space = "mail"

func api(w http.ResponseWriter, req *http.Request) {
	// Parsing body
	req.ParseForm()

	switch req.Method {
	case http.MethodGet:
		// GET method

		var logString string

		logString += "GET request FROM" + req.RemoteAddr + " ... AT " + time.Now().String() + "\n"

		// Parsing query parameters
		id := req.Form["id"]

		logString += "Requested ID: " + req.Form["id"][0] + "\n"

		// String to uint conv
		idU, err := strconv.ParseUint(id[0], 10, 64)
		// If error => 404, bc key is type of unsigned int
		if err != nil {
			logString += "Response: 404\n\n"
			logger(logString)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// SELECT query to Tarantool
		resp, _ := Conn.Select(space, "primary", 0, 1, tarantool.IterEq, []interface{}{idU})

		// If no such key => 404 Status
		if len(resp.Data) == 0 {
			logString += "Response: 404\n\n"
			logger(logString)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Parsing response from SELECT Tarantool query
		myStr := resp.Data[0].([]interface{})
		mymap := myStr[1].(map[interface{}]interface{})

		myGoodMap := make(map[string]string)

		for k, v := range mymap {
			strK := k.(string)
			strV := v.(string)
			myGoodMap[strK] = strV
		}

		bio := Bio{
			Name:      	myGoodMap["name"],
			SecondName: myGoodMap["secondName"],
		}

		jsonBio, _ := json.Marshal(bio)
		logString += "Response: " + string(jsonBio) + "\n\n"
		logger(logString)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(jsonBio))

	case http.MethodPost:
		// POST method

		var logString string
		logString += "POST request FROM" + req.RemoteAddr + " ... AT " + time.Now().String() + "\n"

		// Reading JSON data
		body := make([]byte, 1024)
		req.Body.Read(body)
		body = bytes.Trim(body, "\x00")

		logString += "Add: " + string(body) + "\n"

		// Parsing JSON body to custom struct
		data := new(Request)
		err := json.Unmarshal(body, &data)

		// Help vars for Tarantool query
		bio := make(map[string]string)
		bio["name"] = data.Val.Name
		bio["secondName"] = data.Val.SecondName

		// If JSON body is incorrect => 400 status code
		if data.Val.Name == "" || data.Val.SecondName == "" || err != nil {
			logString += "Response: 400 Bad Body\n\n"
			logger(logString)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Tarantool db query
		resp, _ := Conn.Insert(space, []interface{}{data.Key, bio})

		if resp.Code == tarantool.ErrTupleFound {
			logString += "Response: Key already exist 409\n\n"
			logger(logString)
			w.WriteHeader(http.StatusConflict)
			return
		}

		logString += "Success\n\n"
		logger(logString)

	case http.MethodPut:
		// PUT method

		var logString string

		logString += "PUT request FROM" + req.RemoteAddr + " ... AT " + time.Now().String() + "\n"

		// Reading JSON data
		body := make([]byte, 1024)
		req.Body.Read(body)
		body = bytes.Trim(body, "\x00")

		logString += "Add: " + string(body) + "\n"

		// Parsing JSON body to custom struct
		data := new(Bio)
		err := json.Unmarshal(body, &data)

		// Help vars for Tarantool query
		bio := make(map[string]string)
		bio["name"] = data.Name
		bio["secondName"] = data.SecondName

		// If JSON body is incorrect => 400 status code
		if data.Name == "" || data.SecondName == "" || err != nil {
			logString += "Response: 400 Bad Body\n\n"
			logger(logString)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		id := req.Form["id"]

		logString += "Requested ID: " + req.Form["id"][0] + "\n"

		idU, err := strconv.ParseUint(id[0], 10, 64)
		if err != nil {
			logString += "Response: 404\n\n"
			logger(logString)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		resp, err := Conn.Delete(space, "primary", []interface{}{idU})

		// If len(resp.Data) => there is no element with this key
		if len(resp.Data) == 0 {
			logString += "Response: No such key 404\n\n"
			logger(logString)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Tarantool db query
		resp, _ = Conn.Insert(space, []interface{}{idU, bio})

		logString += "Response: Successfully updated\n\n"
		logger(logString)

	case http.MethodDelete:
		// DELETE method

		var logString string

		logString += "Delete request FROM" + req.RemoteAddr + " ... AT " + time.Now().String() + "\n"

		id := req.Form["id"]

		logString += "Requested ID: " + req.Form["id"][0] + "\n"

		idU, err := strconv.ParseUint(id[0], 10, 64)
		if err != nil {
			logString += "Response: No such key 404\n\n"
			logger(logString)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		resp, err := Conn.Delete(space, "primary", []interface{}{idU})

		// If len(resp.Data) => there is no element with this key
		if len(resp.Data) == 0 {
			logString += "Response: No such key 404\n\n"
			logger(logString)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		logString += "Response: Successfully Deleted\n\n"
		logger(logString)

	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
