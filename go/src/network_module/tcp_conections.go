package network_module

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
	"os"
	)


func manageTCPconnections(){
	connections := connectionMap{make(map[string]connectionChans)}
	go Listen_for_TCP_connection()
	for {
		select {
		case newIP := <-internalChan.newIP:
			connections.handleNewIP(newIP)

		case newTCPConnection := <-internalChan.updateTCPMap:
			connections.handleNewConnection(newTCPConnection)

		case errorIP := <-internalChan.connectFail:
			connections.handleFailedToConnect(errorIP)

		case errorIP := <-internalChan.connectionError:
			connections.handleConnectionError(errorIP)

		case closeIP := <-internalChan.closeConnection:
			connections.handleCloseConnection(closeIP)

		case mail := <-externalChan.SendToAll:
			connections.handleSendToAll(mail)

		case mail := <-externalChan.SendToOne:
			connections.handleSendToOne(mail)

		case <-internalChan.quitTCPmap:
			return
		}
	} 
}

func (connections *connectionMap) handleNewIP(newIP string){
	_, inMap := connections.tcpMap[newIP]
	if !inMap {
		go ConnectTCP(newIP)
	} else {
		fmt.Println("network.tcp_connections.handleNewIP -->", newIP, " already connected")
	}
}

func (connections *connectionMap) handleNewConnection(newConnection tcpConnection){
	_,inMap := connections.tcpMap[newConnection.ip]
	if !inMap {
		connections.tcpMap[newConnection.ip] = connectionChans{send: make(chan Mail), quit: make(chan bool)}
		fmt.Println("network.tcp_connections.handleNewConnection --> Connected to: ", newConnection.ip)
		newConnection.sendChan = connections.tcpMap[newConnection.ip].send
		newConnection.quit = connections.tcpMap[newConnection.ip].quit
		go newConnection.handleConnection()
	} else {
		fmt.Println("network.tcp_connections.handleNewConnection -->", newConnection.ip, " connection already in connectionMap")
		newConnection.socket.Close()
	}
}

func (connections *connectionMap) handleFailedToConnect(errorIP string){
	_,inMap := connections.tcpMap[errorIP]
	if inMap{
		fmt.Println("network.tcp_connections.handleFailedToConnect --> Already a connected to ", errorIP)
	} else {
		fmt.Println("network.tcp_connections.handleFailedToConnect --> Could not connect to ", errorIP)
		internalChan.errorIP <- errorIP
	}
}

func (connections *connectionMap) handleConnectionError(errorIP string){
	_,inMap := connections.tcpMap[errorIP]
	if inMap{
		delete(connections.tcpMap, errorIP)
	} 
	go ConnectTCP(errorIP)
}

func (connections *connectionMap) handleCloseConnection(closeIP string){
	connectionChans, inMap := connections.tcpMap[closeIP]
	if inMap{
		select{
		case connectionChans.quit <- true:
		case <-time.After(10 * time.Millisecond):
		}
		delete(connections.tcpMap, closeIP)
	} else {
		fmt.Println("network.tcp_connections.handleCloseConnection --> No connection to close ", closeIP)
	}
}

func (connections *connectionMap) handleSendToAll(mail Mail){
	if len(connections.tcpMap) != 0{
		for _,connectionChans := range connections.tcpMap{
			connectionChans.send <- mail
		}
	}
}

func (connections *connectionMap) handleSendToOne(mail Mail){
	switch mail.IP{
	case "":
		size := len(connections.tcpMap)
		if size != 0{
			for _,connectionChans := range connections.tcpMap{
				connectionChans.send <- mail
				break
			}
		}
	default:
		connectionChans, inMap := connections.tcpMap[mail.IP]
		if inMap{
			connectionChans.send <- mail
		} else {
			internalChan.errorIP <- mail.IP
		}
	}
}

func (connection *tcpConnection) inbox(quitInbox chan bool) {
	var msg [512]byte
	for{
		nBytes, err := connection.socket.Read(msg[0:])
		err = nil
		switch err{
		case nil:
			fmt.Println("I was never here biach")
			newMail := Mail{IP: connection.ip, Msg: msg[0:nBytes]}
			externalChan.Inbox <- newMail

		default:
			fmt.Println("network.tcp_connections.inbox --> Error:", err)
			time.Sleep(IMAPERIOD * IMALOSS * 2 * time.Millisecond)
			select {
			case quitInbox <- true: 
			case <-time.After(WRITEDL * time.Millisecond):
			}
		return
		}
	}
}
func (connection *tcpConnection) handleConnection(){
	quitInbox := make(chan bool)
	go connection.inbox(quitInbox)
	fmt.Println("network.tcp_connections.handleConnection --> handleConnection for", connection.ip, "is running")
	for {
		select{
		case mail := <-connection.sendChan:
			connection.socket.SetWriteDeadline(time.Now().Add(WRITEDL * time.Millisecond))
			//_, err := connection.socket.Write(mail.Msg)
			nBytes, err := connection.socket.Write(mail.Msg)
			if err == nil {
				fmt.Println("Network.handleConnection--> Msg of", nBytes, "bytes sent to ", connection.ip)
			} else {
				fmt.Println("***Network.handleConnection--> Error sending message to ", connection.ip, err)
				internalChan.connectionError <- connection.ip //Notify manager of fault
			}
		case <-connection.quit:
			connection.socket.Close()
			fmt.Println("Network.handleConnections--> Connection to ", connection.ip, " has been terminated.")
			return

		case <-quitInbox:
			connection.socket.Close()
			fmt.Println("Network.handleConnections--> Connection to ", connection.ip, " has been terminated.")
			internalChan.connectionError <- connection.ip
			return

		}
	}
}

func ConnectTCP(ipAdr string){
	attempts := 0

	for attempts < CONNATMPT {
		fmt.Println("network.tcp_connections.connectTCP --> attempting to connect to ", ipAdr)
		_, err := net.ResolveTCPAddr("tcp",ipAdr + ":" + TCPport)
		if checkError(err){
			fmt.Println("network.tcp_connections.connectTCP --> ResolveTCPAddr failed")
			attempts ++
			time.Sleep(100 * time.Millisecond)
		}else{
			service := ipAdr+":9191"
			randSleep := time.Duration(rand.Intn(500)+500) * time.Microsecond
			fmt.Println("network.tcp_connections.connectTCP --> randSleep:", randSleep)
			time.Sleep(randSleep)
			socket, err := net.Dial("tcp", service);	
			if checkError(err){
				fmt.Println("network.tcp_connections.connectTCP --> DialTCP error when connecting to", ipAdr)
				attempts++
				time.Sleep(500 * time.Millisecond)
			}else{
				newTCPConnection := tcpConnection{ip: ipAdr, socket: socket}
				internalChan.updateTCPMap <- newTCPConnection
				break
			}
		}
	}
}

func Listen_for_TCP_connection(){
	service := ":9191"

	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		fmt.Println("network.tcp_connections.listen_for_TCP_connection --> TCP resolve error")
		internalChan.setupFail <- true
	} else{
		listenSocket, err := net.ListenTCP("tcp", tcpAddr)
		if err != nil {
			fmt.Println("network.tcp_connections.listen_for_TCP_connection --> TCP listen error")
			internalChan.setupFail <- true
		} else {
			fmt.Println("network.tcp_connections.listen_for_TCP_connection --> listening for new connections")
			for{
				select{
				case <- internalChan.quitListenTCP:
					return

				default:
					socket, err := listenSocket.Accept()
					if err != nil {
						ip := clean_up_IP(socket.RemoteAddr().String())
						newTCPConnection := tcpConnection{ip: ip, socket: socket}
						internalChan.updateTCPMap <- newTCPConnection
					}
				}
			}
		}
	}
}

func clean_up_IP(garbage string) (cleanIP string) {
        split := strings.Split(garbage, ":")
        cleanIP = split[0]
        return
}	


func checkError(err error) bool{
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
		return true
	}
	return false
}

