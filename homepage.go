package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/julienschmidt/sse"
	"html/template"
	"net/http"
	"time"
)

type HomePage struct {
	Time string
}

func serveHomepage(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	writingSync.Lock()
	programIsRunning = true
	writingSync.Unlock()
	var homepage HomePage
	homepage.Time = time.Now().Format("02/01/2006, 15:04:05")
	tmpl := template.Must(template.ParseFiles("html/homepage.html"))
	_ = tmpl.Execute(writer, homepage)
	writingSync.Lock()
	programIsRunning = false
	writingSync.Unlock()
}

func streamTime(timer *sse.Streamer) {
	fmt.Println("Streaming time started")
	for serviceIsRunning {
		timer.SendString("", "time", time.Now().Format("02/01/2006, 15:04:05"))
		time.Sleep(1 * time.Second)
	}
}
