package network_module

import(
		"math/rand"
		"fmt"
		"time"
		)

func Start_network(externalChan NetChannels){
	rand.Seed(time.Now().UTC().UnixNano())
	internalChan.network_internal_chan_init()
	//externalChan = NetChan
	go Send_im_alive()
	go Recieve_im_alive()
	go Watch_im_alive()
	go manageTCPconnections()

	for{
		select{
		case <- internalChan.setupFail:
			fmt.Println("net.Netwok --> setupFail \n ---------- Please wait while we try again ----------")
			internalChan.quitSendIMA <- true
			internalChan.quitWatchIMA <- true
			internalChan.quitRecieveIMA <- true
			internalChan.quitListenTCP <- true
			internalChan.quitTCPmap <- true
			go Send_im_alive()
			go Recieve_im_alive()
			go Watch_im_alive()
			go manageTCPconnections()

		case <- time.After(NETSETUP * time.Millisecond):
			//fmt.Println("Network setup success")
		}
	}
}