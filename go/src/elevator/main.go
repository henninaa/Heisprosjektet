package main

import(
	"driver_module"
	"bank_module"
	)

func main(){
	
	var deadChan = make(chan int, 1)

	driver_module.Elev_init()
	
	go bank_module.Elevator_main_control()

	<-deadChan
}

