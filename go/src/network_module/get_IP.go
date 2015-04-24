package network_module

import (
        "net"
        "strings"
)

func Get_my_IP() string {
        baddr,_ := net.ResolveUDPAddr("udp4", broad_cast + ":" + UDP_port)

        tempConn,_ := net.DialUDP("udp4", nil, baddr)
        
        defer tempConn.Close()
        
        myAddr := tempConn.LocalAddr()
        
        my_IP := strings.Split(myAddr.String(), ":")[0]
        
        return my_IP
}