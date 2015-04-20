package main

import(
	"driver_module"
	"queue_module"
	"sensor_module"
	"FSM_module"
	)

func main(){
	
	var deadChan = make(chan int, 1)

	driver_module.Elev_init()
	queue_module.Init_queue()

	go sensor_module.Sensors()
	go FSM_module.FSM()


	<-deadChan
}

