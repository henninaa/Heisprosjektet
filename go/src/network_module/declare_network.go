package network_module

import "net"

const (
        IMAPERIOD       = 100
        IMALOSS         = 4
        ALIVEWATCH      = 10
        NETSETUP        = 200
        WRITEDL         = 10
        READDL          = 10
        CONNATMPT		= 5
)

var(
	UDPport = "9001"
	TCPport = "9191"
	bcast = "129.241.187.255"
	MyIP = Get_my_IP()
)

type connectionMap struct {
	tcpMap map[string]connectionChans
}

type connectionChans struct{
	send chan Mail
	quit chan bool
}

type tcpConnection struct{
	ip string
	socket net.Conn
	sendChan chan Mail
	quit chan bool
}

type Mail struct {
	IP  string
	Msg []byte
}