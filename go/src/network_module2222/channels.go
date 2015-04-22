package network_module

//////////Internal channels//////////
type internalChannels struct{
	newIP 			chan string
	connectionFail 	chan string
	connectionError chan string
	closeConnection chan string
	ima 			chan string
	errorIP			chan string
	connectFail     chan string
	setupFail 		chan bool
	quitSendIMA 	chan bool
	quitWatchIMA 	chan bool
	quitRecieveIMA 	chan bool
	quitListenTCP 	chan bool
	quitTCPmap		chan bool
	quitListenIMA	chan bool
	updateTCPMap 	chan tcpConnection
}

var internalChan internalChannels

func (internalChan *internalChannels) network_internal_chan_init(){
	internalChan.newIP = make(chan string)
	internalChan.closeConnection = make(chan string)
	internalChan.ima = make(chan string)

	internalChan.quitSendIMA = make(chan bool)
	internalChan.quitWatchIMA = make(chan bool)
	internalChan.quitListenIMA = make(chan bool)
	internalChan.quitRecieveIMA = make(chan bool)
	internalChan.quitListenTCP = make(chan bool)
	internalChan.quitTCPmap = make(chan bool)
	internalChan.setupFail = make(chan bool)

	internalChan.connectionFail = make(chan string)
	internalChan.connectionError = make(chan string)
	internalChan.errorIP = make(chan string)
	internalChan.connectFail = make(chan string)

	internalChan.updateTCPMap = make(chan tcpConnection)
}


//////////External channels//////////
type NetChannels struct{
	GetDeadElevator 	chan string
	SendDeadElevator 	chan string
	Panic 				chan bool
	Inbox 				chan Mail
	SendToAll 			chan Mail
	SendToOne 			chan Mail
}

	
var externalChan NetChannels


func (externalChan *NetChannels) Network_external_chan_init(){
	externalChan.GetDeadElevator = make(chan string)
	externalChan.SendDeadElevator = make(chan string)
	externalChan.Panic = make(chan bool)
	externalChan.Inbox = make(chan Mail)
	externalChan.SendToOne = make (chan Mail)
	externalChan.SendToAll = make (chan Mail)

}