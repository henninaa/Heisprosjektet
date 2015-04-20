package main

import (
		"network_module"
		"time"
		"fmt"
		)

func main(){
	var NetChan network_module.NetChannels

	go NetChan.Network_external_chan_init()
	go network_module.Start_network(NetChan)
	go repeater(NetChan)
	for {
		fmt.Println("Moren din ")
		newMail := <- NetChan.Inbox
		fmt.Println("Moren din 2 ",newMail)
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