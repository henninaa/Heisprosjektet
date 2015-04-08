package network_module

import (
	"fmt"
	"net"
	"os"
	"time"
	//"encoding/json"
	"math/rand"
	"strings"
	)




func ConnectTCP(ipAdr string, port string){
	attempts := 0

	for attempts < 5 {
		fmt.Println("Network.connectTCP--> attempting to connect to ", ipAdr)
		_, err := net.ResolveTCPAddr("tcp",ipAdr+port)
		if checkError(err){
			fmt.Println("***Network.connectTCP--> ResolveTCPAddr failed")
			attempts ++
			time.Sleep(100 * time.Millisecond)
		}else{
			service := ipAdr+":9191"
			randSleep := time.Duration(rand.Intn(500)+500) * time.Microsecond
			fmt.Println("Network.connectTCP--> randSleep:", randSleep)
			time.Sleep(randSleep)
			socket, err := net.Dial("tcp", service);	
			if checkError(err){
				fmt.Println("***Network.connectTCP--> DialTCP error when connecting to", ipAdr)
				attempts++
				time.Sleep(500 * time.Millisecond)
			}else{
				newTCPConnection := tcpConnection{ip: ip, socket: socket}
				internalChan.updateTCPMap <- newTCPConnection
				break
			}
		}
	}
}

func Listen_for_TCP_connection(){
	service := ":9191"

	tcpAddr, err := net.ResolveTCPAddr("tcp", service)

	listenSocket, err := net.ListenTCP("tcp", tcpAddr)

	fmt.Println("Network.connectTCP--> listening for new connections")


}

func cleanUpIP(garbage string) (cleanIP string) {
        split := strings.Split(garbage, ":") //Hackjob to separate ip from local socket. (Seems like a "fault" in the net package)
        cleanIP = split[0]
        return
}	

func checkError(err error) bool{
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
		return true
	}
	return false
}

func handleError(err error){
	if err != nil{
		fmt.Println(err)
	}

}