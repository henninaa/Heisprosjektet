package main

import(
	"driver_module"
	"network_module"
	)

func main(){
	
	driver_elev.Elev_init()

	network_module.Network()

	deadChan := make(chan int)
	<-deadChan
}

