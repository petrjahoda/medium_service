package main

import (
	"github.com/julienschmidt/httprouter"
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

	time.Sleep(10 * time.Second)

	var homepage HomePage
	homepage.Time = time.Now().String()
	tmpl := template.Must(template.ParseFiles("html/homepage.html"))
	_ = tmpl.Execute(writer, homepage)
	writingSync.Lock()
	programIsRunning = false
	writingSync.Unlock()
}
