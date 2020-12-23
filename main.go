package main

import (
	"fmt"
	"github.com/kardianos/service"
	"sync"
	"time"
)

const serviceName = "Medium service"
const serviceDescription = "Simple service, just for fun"

type program struct{}

func (p program) Start(s service.Service) error {
	fmt.Println(s.String() + " started")
	go p.run()
	return nil
}

func (p program) Stop(s service.Service) error {
	fmt.Println(s.String() + " stopped")
	return nil
}

func (p program) run() {

	for serviceIsRunning {
		writingSync.Lock()
		programIsRunning = true
		writingSync.Unlock()
		fmt.Println("Service is running")
		time.Sleep(2 * time.Second)
		writingSync.Lock()
		programIsRunning = false
		writingSync.Unlock()

	for {
		fmt.Println("Service is running")
		time.Sleep(1 * time.Second)
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
