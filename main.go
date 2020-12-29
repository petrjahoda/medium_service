package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/julienschmidt/sse"
	"github.com/kardianos/service"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	serviceIsRunning bool
	programIsRunning bool
	writingSync      sync.Mutex
)

const serviceName = "Medium service"
const serviceDescription = "Simple service, just for fun"

type program struct{}

func (p program) Start(s service.Service) error {
	fmt.Println(s.String() + " started")
	writingSync.Lock()
	serviceIsRunning = true
	writingSync.Unlock()
	go p.run()
	return nil
}

func (p program) Stop(s service.Service) error {
	writingSync.Lock()
	serviceIsRunning = false
	writingSync.Unlock()
	for programIsRunning {
		fmt.Println(s.String() + " stopping...")
		time.Sleep(1 * time.Second)
	}
	fmt.Println(s.String() + " stopped")
	return nil
}

func (p program) run() {
	logDirectoryCheck()
	go deleteOldLogFiles(48 * time.Hour)
	logInfo("RUN", "Program is running")
	router := httprouter.New()
	timer := sse.New()
	router.ServeFiles("/js/*filepath", http.Dir("js"))
	router.ServeFiles("/css/*filepath", http.Dir("css"))
	router.GET("/", serveHomepage)
	router.POST("/get_time", getTime)
	router.Handler("GET", "/time", timer)
	go streamTime(timer)
	err := http.ListenAndServe(":82", router)
	if err != nil {
		fmt.Println("Problem starting web server: " + err.Error())
		os.Exit(-1)
	}
}

func main() {
	serviceConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: serviceDescription,
	}
	prg := &program{}
	s, err := service.New(prg, serviceConfig)
	if err != nil {
		fmt.Println("Cannot create the service: " + err.Error())
	}
	err = s.Run()
	if err != nil {
		fmt.Println("Cannot start the service: " + err.Error())
	}
}
