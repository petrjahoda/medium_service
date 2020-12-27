package main

import (
	"encoding/json"
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
	var data TimeDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		fmt.Println(err.Error())
		var responseData TimeDataOutput
		responseData.Result = "nok"
		responseData.Text = "problem with user json data"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		return
	}
	fmt.Println(data.Name)
	fmt.Println(data.Time)
	timer := time.Now()
	time.Sleep(1 * time.Second)
	end := time.Since(timer)
	fmt.Println("processing takes: " + end.String())
	var responseData TimeDataOutput
	responseData.Result = "ok"
	responseData.Text = "everything went smooth"
	responseData.Time = time.Now().Format("02/01/2006, 15:04:05")
	responseData.Duration = end.String()
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	return
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
