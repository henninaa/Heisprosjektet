package network_module

import (
	"fmt"
	"net"
	"os"
	"time"
	"encoding/json"
	"math/rand"
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

func Get_my_IP()string {
	allIPs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("network.GetMyIP()--> Error receiving IPs. IP set to localhost. Consider setting IP manually")
		manualIP := ""
		fmt.Scanf("%s",&manualIP)
		return manualIP
	}
	
	IPString := make([]string, len(allIPs))
	for i := range allIPs {
			temp := allIPs[i].String()
			ip := strings.Split(temp, "/")
			IPString[i] = ip[0]
	}
	var myIP string
	for i := range IPString {
		if IPString[i][0:2] == "129" {
			myIP = IPString[i]
		}
	}
	return myIP
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
				//Legg til socket i oversikt over oppkoblinger
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

func PrintRecievedMessages() {
	for{
		message := <- ImAliveThread
		fmt.Println("His message: ", message)
	}	
}




