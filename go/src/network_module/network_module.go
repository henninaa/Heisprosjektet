/*
Mål med denne modulen:
1. Broadcast "I'm alive" - Go
2. Motta "I'm alive" - Go
3. Lese av IP til broadcasterene
4. Koble opp TCP mot broadcasterene
5. Få Json fra resten av programmet
6. Sende Json til alle broadcasterene - Go
7. Motta Json fra alle broadcasterene - Go
8. Levere Json til resten av programmet
9. Sjekke om alle er i live - Go


Fullførte mål:
 - Nesten 1: Kan sende til enkeltadresser, ikke broadcasting
 - Nesten 2: Kan motta fra enkeltadresser, ikke broadcasting

*/

package network_module

import(
		"net"
		"fmt"
		"strings"
		)

type Kanal struct{
	Send_ko chan int
}
var Kan Kanal

func Init_chan() {
	var unimportant_variable bool
	unimportant_variable = true
}

func Get_my_IP()string {
	allIPs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("network_module.Get_my_IP()--> Error receiving IP. Type your IP here: ")
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

func Network(){

	go SendImAlive()
	go RecieveImAlive()
	go PrintRecievedMessages()

}