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

const (
	ORDER_TAKEN 		= "OTK"
	ORDER_EXECUTED		= "OEX"
	DELIVER_ORDER 		= "DLO"
	ERROR_MSG 			= "ERM"
	TAKE_BACKUP_ORDER	= "TBO"
	BACKUP_ORDER_COMPLETE= "BOC"
	TAKE_NEW_ORDER		= "TNO"
	TAKE_BACKUP_FLOOR	= "TBF"
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
	Msg message
}

type Message struct {
	IP string
	msg_type string
	floor int
	button_type int
		
}

func (mail * Mail) Make_mail(ip string, msg_type string, floor int, button_type int){


}
