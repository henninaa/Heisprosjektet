package network_module

type internalChannels struct {
        setupfail       chan bool
        ima             chan string
        newIP           chan string
        deadElevator    chan string
        updateTCPMap    chan tcpConnection
        connectFail     chan string
        closeConn       chan string
        errorIP         chan string
        connectionError chan string
        deleteConn      chan string
        quitImaSend     chan bool
        quitImaListen   chan bool
        quitImaWatcher  chan bool
        quitListenTCP   chan bool
        quitTCPMap      chan bool
}

type NetChannels struct {
        GetDeadElevator  chan string
        SendDeadElevator chan string
        SendToAll        chan Mail
        SendToOne        chan Mail
        Inbox            chan Mail
        NumOfPeers       chan int
        Panic            chan bool
}

var internalChan internalChannels
var externalChan NetChannels

func (internalChan *internalChannels) init() {
        internalChan.setupfail = make(chan bool)
        internalChan.ima = make(chan string)
        internalChan.newIP = make(chan string)
        internalChan.deadElevator = make(chan string)
        internalChan.updateTCPMap = make(chan tcpConnection)
        internalChan.connectFail = make(chan string)
        internalChan.connectionError = make(chan string)
        internalChan.errorIP = make(chan string)
        internalChan.closeConn = make(chan string)
        internalChan.deleteConn = make(chan string)
        internalChan.quitImaSend = make(chan bool)
        internalChan.quitImaListen = make(chan bool)
        internalChan.quitImaWatcher = make(chan bool)
        internalChan.quitListenTCP = make(chan bool)
        internalChan.quitTCPMap = make(chan bool)
}

func (externalChan *NetChannels) NetChanInit() {
        externalChan.GetDeadElevator = make(chan string)
        externalChan.SendDeadElevator = make(chan string)
        externalChan.SendToAll = make(chan Mail)
        externalChan.SendToOne = make(chan Mail)
        externalChan.Inbox = make(chan Mail)
        externalChan.NumOfPeers = make(chan int)
        externalChan.Panic = make(chan bool)
}