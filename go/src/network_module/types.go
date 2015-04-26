package network_module

import (
	"net"
	"driver_module"
	)

const (
	IMA_PERIOD              = 25
	IMA_LOSS                = 4
	ALIVE_WATCH             = 10
	NET_SETUP               = 200
	DIAL_INT                = 50
	CONN_ATMPT              = 5
	WRITE_DL                = 10
	READ_DL                 = 10
	CONN_FAIL_TIMEOUT       = 2 * NET_SETUP
)

const (
	ORDER_TAKEN             = "OTK"
	ORDER_EXECUTED          = "OEX"
	TAKE_BACKUP_ORDER       = "TBO"
	BACKUP_ORDER_COMPLETE   = "BOC"
	TAKE_NEW_ORDER          = "TNO"
	TAKE_BACKUP_FLOOR       = "TBF"
	ENGINE_FAILURE			= "ENF"
	ENGINE_RECOVERY			= "ENR"
)


var (
	broad_cast      = "129.241.187.255"
	LOCAL_IP        = Get_my_IP()
	UDP_port        = "9001"
	TCP_port        = "9191"
)

type conn_map struct {
	tcp_map map[string]conn_chans
}

type conn_chans struct {
	send chan Mail
	quit chan bool
}

type tcp_connection struct {
	ip       string
	socket   net.Conn
	sendChan chan Mail
	quit     chan bool
}

type Message struct{
	Msg_type string
	Floor int
	Dir driver_module.Elev_button_type_t
	Cost_answer int
}

type Mail struct {
	IP  string
	Msg Message
}


