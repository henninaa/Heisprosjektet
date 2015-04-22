package network_module

import(
		"net"
		"fmt"
		"strings"
		)


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
		if IPString[i][0:3] == "129" {
			myIP = IPString[i]
		}
	}
	return myIP
}