package main

import (
	"github.com/kardianos/service"
	"time"
)

const serviceName = "Medium service"
const serviceDescription = "Simple service, just for fun"

type program struct{}

func (p program) Start(s service.Service) error {
	println(s.String() + " started")
	go p.run()
	return nil
}

func (p program) Stop(s service.Service) error {
	println(s.String() + " stopped")
	return nil
}

func (p program) run() {
	for {
		println("Service is running")
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
		println("Cannot create the service: " + err.Error())
	}
	err = s.Run()
	if err != nil {
		println("Cannot start the service: " + err.Error())
	}
}
