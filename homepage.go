package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/julienschmidt/sse"
	"html/template"
	"net/http"
	"time"
)

type HomePage struct {
	Time string
}

type TimeDataInput struct {
	Name string
	Time string
}

type TimeDataOutput struct {
	Result   string
	Text     string
	Time     string
	Duration string
}

func getTime(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("MAIN", "Get time function called from "+request.RemoteAddr)
	timer := time.Now()
	var data TimeDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", "Problem decoding data: "+err.Error())
		var responseData TimeDataOutput
		responseData.Result = "nok"
		responseData.Text = "problem with user json data"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Get time function ended in "+time.Since(timer).String())
		return
	}
	var responseData TimeDataOutput
	responseData.Result = "ok"
	responseData.Text = "everything went smooth"
	responseData.Time = time.Now().Format("02/01/2006, 15:04:05")
	responseData.Duration = time.Since(timer).String()
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo("MAIN", "Get time function ended in "+time.Since(timer).String())
	return
}

func serveHomepage(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("MAIN", "Serving homepage to the user")
	timer := time.Now()
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
	logInfo("MAIN", "Homepage served in "+time.Since(timer).String())
}

func streamTime(timer *sse.Streamer) {
	logInfo("MAIN", "Streaming time started")
	for serviceIsRunning {
		timer.SendString("", "time", time.Now().Format("02/01/2006, 15:04:05"))
		time.Sleep(1 * time.Second)
	}
}
