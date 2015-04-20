package main

import (
		"network_module"
		"time"
		"fmt"
		)

func main(){
	var NetChan network_module.NetChannels

	//go network_module.Network_external_chan_init()
	go network_module.Start_network(NetChan)
	go repeater(NetChan)
	for {
		newMail := <- NetChan.Inbox
		fmt.Println(newMail)
	}

	deadChan := make(chan int)
	<-deadChan
}

func repeater(netChan network_module.NetChannels){
	for {
		mail := network_module.Mail{Msg: []byte("TST")}
		netChan.SendToAll <- mail
		time.Sleep(100*time.Millisecond)
	}
}