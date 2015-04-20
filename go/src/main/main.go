package main

import . "network_module"

func main(){
	var NetChan NetChannels

	go Start_network(NetChan)

	deadChan := make(chan int)
	<-deadChan
}

