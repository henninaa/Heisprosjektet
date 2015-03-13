package network_module

import (
	"fmt"
	"net"
	"os"
	"time"
	"encoding/json"
	"strings"
	)

var ImAliveThread = make(chan ImAliveMessage,1)

const(
	port = ":20003"
	IP = "129.241.187.255"
)

type ImAliveMessage struct{
	Msg string
}


func SendImAlive(){
	var message ImAliveMessage
	message.Msg = "I'm alive"

	for{
		addr, err := net.ResolveUDPAddr("udp", IP+port)
		handleError(err)

		sock, err := net.DialUDP("udp",nil,addr)
		handleError(err)

		bmsg, err := json.Marshal(message)

		sock.Write(bmsg)

		handleError(err)

		time.Sleep(100*time.Millisecond)
	}
}



func RecieveImAlive() {

	buffer := make([]byte, 1024)
	var recievedMessage ImAliveMessage

	addr, _ := net.ResolveUDPAddr("udp",port)
	sock, _ := net.ListenUDP("udp", addr)

	for{
		n, _,error := sock.ReadFromUDP(buffer)
		if error != nil{
			handleError(error)
		}

		json.Unmarshal(buffer[:n], &recievedMessage)

		ImAliveThread <- recievedMessage
	}
}


func PrintRecievedMessages() {
	for{
		message := <- ImAliveThread
		fmt.Println("His message: ", message)
	}	
}




