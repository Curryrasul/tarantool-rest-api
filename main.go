package main

import (
	"github.com/tarantool/go-tarantool"
	"net/http"
	"os"
)

// Global var for connection
var Conn *tarantool.Connection
var LogFile *os.File
var LogFileName = "./log.txt"

func main() {
	Conn, _ = tarantool.Connect("127.0.0.1:3311", tarantool.Opts{
		User: "admin",
		Pass: "pass",
	})

	LogFile, _ = os.OpenFile(LogFileName, os.O_RDWR, 0644)

	defer Conn.Close()
	defer LogFile.Close()

	// Starting the server on 8080 Port
	http.HandleFunc("/kv", api)
	http.ListenAndServe("0.0.0.0:8080", nil)
}
