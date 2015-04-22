package network_module

import (
        //"fmt"
        "math/rand"
        "net"
        "strings"
        "time"
        "encoding/gob"
        "printc"
)

func manageTCPConnections() {
        connections := connMap{make(map[string]connChans)}
        go listenForTCPConnections()
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

                case closeIP := <-internalChan.closeConn:
                        connections.handleCloseConnection(closeIP)

                case mail := <-externalChan.SendToAll:
                        connections.handleSendToAll(mail)

                case mail := <-externalChan.SendToOne:
                        connections.handleSendToOne(mail)

                case <-internalChan.quitTCPMap:
                        return
                } 
        } 
} 

func (connections *connMap) handleNewIP(newIP string) {
        _, inMap := connections.tcpMap[newIP]
        if !inMap {
                go connectTCP(newIP)
        } else {
                printc.DataWithColor(printc.COLOR_YELLOW,"network.monitorTCPConnections-->", newIP, "already in connections")
        }
}

func (connections *connMap) handleNewConnection(conn tcpConnection) {
        _, inMap := connections.tcpMap[conn.ip]
        if !inMap {
                connections.tcpMap[conn.ip] = connChans{send: make(chan Mail), quit: make(chan bool)} 
                printc.DataWithColor(printc.COLOR_GREEN,"network.monitorTCPConnections---> Connection made to ", conn.ip)
                conn.sendChan = connections.tcpMap[conn.ip].send
                conn.quit = connections.tcpMap[conn.ip].quit
                go conn.handleConnection()
                go peerUpdate(len(connections.tcpMap))
        } else {
                printc.DataWithColor(printc.COLOR_YELLOW,"network.monitorTCPConnections--> A connection already exist to", conn.ip)
                conn.socket.Close()
        }
}

func (connections *connMap) handleFailedToConnect(errorIP string) {
        _, inMap := connections.tcpMap[errorIP]
        if inMap {
                printc.DataWithColor(printc.COLOR_YELLOW,"network.monitorTCPConnections--> Could not dial up ", errorIP, "but a connection already exist")
        } else {
                printc.DataWithColor(printc.COLOR_RED,"network.monitorTCPConnections--> Could not connect to ", errorIP)
                internalChan.errorIP <- errorIP //Notify imaWatcher of erroneous ip. Maybe it has timed out?
                printc.DataWithColor(printc.COLOR_CYAN,"YOU SHALL NOT PASS!!!!")
        }
}

func (connections *connMap) handleConnectionError(errorIP string) {
        _, inMap := connections.tcpMap[errorIP]
        if inMap {
                delete(connections.tcpMap, errorIP)
        }
        go connectTCP(errorIP)
}

func (connections *connMap) handleCloseConnection(closeIP string) {
        connChans, inMap := connections.tcpMap[closeIP]
        printc.DataWithColor(printc.COLOR_CYAN,"ConnChans ", connChans, "inMap ", inMap)
        if inMap {
                select {
                case connChans.quit <- true:
                case <-time.After(10 * time.Millisecond):
                }
                delete(connections.tcpMap, closeIP)
                numOfConns := len(connections.tcpMap)
                if numOfConns == 0 {
                        go peerUpdate(numOfConns)
                }
        } else {
                printc.DataWithColor(printc.COLOR_YELLOW,"network.monitorTCPConnections--> No connection to close ", closeIP)

        }
}

func (connections *connMap) handleSendToOne(mail Mail) {
        printc.DataWithColor(printc.COLOR_BLUE,"IP: ", mail.IP)
        switch mail.IP {
        case "":
                size := len(connections.tcpMap)
                if size != 0 {
                        for _, connChans := range connections.tcpMap {
                                connChans.send <- mail 
                                break
                        }
                }
        default:
                connChans, inMap := connections.tcpMap[mail.IP]
                if inMap {
                        connChans.send <- mail
                } else {
                        internalChan.errorIP <- mail.IP
                }
        }
}

func (connections *connMap) handleSendToAll(mail Mail) {
        if len(connections.tcpMap) != 0 {
                for _, connChans := range connections.tcpMap {
                        connChans.send <- mail
                }
        }
}

func (conn *tcpConnection) handleConnection() {
        quitInbox := make(chan bool)
        go conn.inbox(quitInbox)
        printc.DataWithColor(printc.COLOR_BLUE,"Network.handleConnection--> handleConnection for", conn.ip, "is running")
        conEnc := gob.NewEncoder(conn.socket)
        for {
                select {
                case mail := <-conn.sendChan:
                        encodedMsg := mail.Msg
                        err := conEnc.Encode(&encodedMsg)
                        if err == nil {
                                printc.DataWithColor(printc.COLOR_CYAN, "Message sendt without problem :)")
                        } else {
                                printc.DataWithColor(printc.COLOR_RED,"Network.handleConnection--> Error sending message to ", conn.ip, err)
                                internalChan.connectionError <- conn.ip 
                        }
                case <-conn.quit:
                        conn.socket.Close()
                        printc.DataWithColor(printc.COLOR_YELLOW,"Network.handleConnections--> Connection to ", conn.ip, " has been terminated.")
                        return
                case <-quitInbox:
                        conn.socket.Close()
                        printc.DataWithColor(printc.COLOR_YELLOW,"Network.handleConnections--> Connection to ", conn.ip, " has been terminated.")
                        internalChan.connectionError <- conn.ip
                        return
                }
        }
}

func (conn *tcpConnection) inbox(quitInbox chan bool) {
        conDec := gob.NewDecoder(conn.socket)
        for {
                decodedMsg := new(Message)
                err := conDec.Decode(decodedMsg)
                switch err {
                case nil:
                        newMail := Mail{IP: conn.ip, Msg: *decodedMsg}
                        externalChan.Inbox <- newMail
                default:
                        printc.DataWithColor(printc.COLOR_RED,"Network.inbox--> Error:", err)
                        time.Sleep(IMAPERIOD * IMALOSS * 2 * time.Millisecond)
                        select {
                        case quitInbox <- true: 
                        case <-time.After(WRITEDL * time.Millisecond):
                        }
                        return
                }
        }
}

func connectTCP(ip string) {
        attempts := 0
        for attempts < CONNATMPT {
                printc.DataWithColor(printc.COLOR_YELLOW,"Network.connectTCP--> attempting to connect to ", ip)
                service := ip + ":" + TCPport
                _, err := net.ResolveTCPAddr("tcp4", service)
                if err != nil {
                        printc.DataWithColor(printc.COLOR_RED,"Network.connectTCP--> ResolveTCPAddr failed")
                        attempts++
                        time.Sleep(DIALINT * time.Millisecond)
                } else {
                        randSleep := time.Duration(rand.Intn(500)+500) * time.Microsecond
                        printc.DataWithColor(printc.COLOR_MAGENTA,"Network.connectTCP--> randSleep:", randSleep)
                        time.Sleep(randSleep)
                        socket, err := net.Dial("tcp4", service)
                        if err != nil {
                                printc.DataWithColor(printc.COLOR_RED,"Network.connectTCP--> DialTCP error when connecting to", ip, " error: ", err)
                                attempts++
                                time.Sleep(DIALINT * time.Millisecond)
                        } else {
                                newTCPConnection := tcpConnection{ip: ip, socket: socket}
                                internalChan.updateTCPMap <- newTCPConnection 
                                break
                        }
                }
        }
        if attempts == CONNATMPT {
                select {
                case internalChan.connectFail <- ip: 
                case <-time.After(CONNFAILTIMEOUT * time.Millisecond): 
                        return
                }
        }
}

func listenForTCPConnections() {
        service := ":" + TCPport
        tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
        if err != nil {
                printc.DataWithColor(printc.COLOR_RED,"Network.listenForTCPConnections--> TCP resolve error")
                internalChan.setupfail <- true
        } else {
                listenSock, err := net.ListenTCP("tcp4", tcpAddr)
                if err != nil {
                        printc.DataWithColor(printc.COLOR_RED,"Network.connectTCP--> ListenTCP error")
                        internalChan.setupfail <- true
                } else {
                        printc.DataWithColor(printc.COLOR_YELLOW,"Network.connectTCP--> listening for new connections")
                        for {
                                select {
                                case <-internalChan.quitListenTCP:
                                        return
                                default:
                                        socket, err := listenSock.Accept()
                                        if err == nil {
                                                ip := cleanUpIP(socket.RemoteAddr().String())
                                                newTCPConnection := tcpConnection{ip: ip, socket: socket}
                                                internalChan.updateTCPMap <- newTCPConnection
                                        }
                                }
                        }
                }
        }
}

func peerUpdate(numOfPeers int) {
        select {
        case externalChan.NumOfPeers <- numOfPeers:
        case <-time.After(500 * time.Millisecond):
        }
}

func cleanUpIP(garbage string) (cleanIP string) {
        split := strings.Split(garbage, ":") 
        cleanIP = split[0]
        return
}