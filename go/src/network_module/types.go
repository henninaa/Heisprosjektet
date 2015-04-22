package network_module

import "net"

const (
        IMAPERIOD       = 25
        IMALOSS         = 4
        ALIVEWATCH      = 10
        NETSETUP        = 200
        DIALINT         = 50
        CONNATMPT       = 5
        WRITEDL         = 10
        READDL          = 10
        CONNFAILTIMEOUT = 2 * NETSETUP
)

const (
        ORDER_TAKEN             = "OTK"
        ORDER_EXECUTED          = "OEX"
        DELIVER_ORDER           = "DLO"
        ERROR_MSG               = "ERM"
        TAKE_BACKUP_ORDER       = "TBO"
        BACKUP_ORDER_COMPLETE   = "BOC"
        TAKE_NEW_ORDER          = "TNO"
        TAKE_BACKUP_FLOOR       = "TBF"
)


var (
        bcast   = "129.241.187.255"
        LOCALIP = GetMyIP()
        UDPport = "9001"
        TCPport = "9191"
)

type connMap struct {
        tcpMap map[string]connChans
}

type connChans struct {
        send chan Mail
        quit chan bool
}

type tcpConnection struct {
        ip       string
        socket   net.Conn
        sendChan chan Mail
        quit     chan bool
}

type Message struct{
        Msg_type string
        Floor int
        Dir string
        CostAnsw int
        TakeOrdre string
}

type Mail struct {
        IP  string
        Msg Message
}

func (mail Mail) Make_mail(){
        
}