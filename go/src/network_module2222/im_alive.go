package network_module

import (
	"fmt"
	"net"
	"time"
	)


func Send_im_alive() {
	service := bcast + ":" + UDPport
	addr, err := net.ResolveUDPAddr("udp4", service)
	if err != nil {
		fmt.Println("network.im_alive.Send_im_alive() --> Resolve error", err)
		internalChan.setupFail <- true
	}
	imaSock, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		fmt.Println("network.im_alive.Send_im_alive() --> Dial error", err)
		internalChan.setupFail <- true
	}
	ima := []byte("I'm alive!!!")
	for {
		select {
		case <-internalChan.quitSendIMA:
			return
		default:
			_, err := imaSock.Write(ima)
			if err != nil {
				fmt.Println("network.im_alive.Send_im_alive() --> UDP send error")
			}
			time.Sleep(IMAPERIOD * time.Millisecond)
		}
	}
}


func Recieve_im_alive() {
		service := bcast + ":" + UDPport
		
		addr, err := net.ResolveUDPAddr("udp4", service)
		
		if err != nil {
			fmt.Println("network.im_alive.Recieve_im_alive() --> ResolveUDP error")
			internalChan.setupFail <- true
		}
		
		sock, err := net.ListenUDP("udp4", addr)
		
		if err != nil {
			fmt.Println("network.im_alive.Recieve_im_alive() --> ListenUDP error")
			internalChan.setupFail <- true
		}
		var data [512]byte
		
		for {
			select {
			case <-internalChan.quitListenIMA:
				return
			default:
				_, remoteAddr, err := sock.ReadFromUDP(data[0:])
				if err != nil {
					fmt.Println("network.im_alive.Recieve_im_alive()--> ReadFromUDP error")
					break
				}
				if MyIP != remoteAddr.IP.String() {
					if err == nil {
						elevIP := remoteAddr.IP.String()
						internalChan.ima <- elevIP
					} else {
						fmt.Println("network.im_alive.Recieve_im_alive()--> UDP read error")
				} 
			} 
		} 
	}
}

func Watch_im_alive(){
	peers := make(map[string]time.Time)
	deadline := IMALOSS * IMAPERIOD * time.Millisecond
	for{
		select{
			case ip := <- internalChan.ima:
				_, inMap := peers[ip]
				peers[ip] = time.Now()
				if !inMap {
					fmt.Println("Jeg fant en ny IP!!", ip)
					internalChan.newIP <- ip
				}

			case <- time.After(ALIVEWATCH * time.Millisecond):
				for ip, timestamp := range peers{
					if time.Now().After(timestamp.Add(deadline)) {
						fmt.Println("network.Wach_im_alive --> timeout", ip)
						externalChan.GetDeadElevator <- ip
						internalChan.closeConnection <- ip
						delete(peers, ip)
					}
				}
			case deadIP := <- externalChan.SendDeadElevator:
				internalChan.closeConnection <- deadIP
				delete(peers, deadIP)

			case errorIP := <- internalChan.errorIP:
				_, inMap := peers[errorIP]
				if inMap{
					externalChan.Panic <- true
				}
			case <- internalChan.quitWatchIMA:
				return
		}
	}
}




