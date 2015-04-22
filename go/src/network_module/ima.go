package network_module

import (
        "net"
        "time"
        "printc"
)

func imaListen() {
        service := bcast + ":" + UDPport
        addr, err := net.ResolveUDPAddr("udp4", service)
        if err != nil {
                printc.DataWithColor(printc.COLOR_RED,"network.IMAListen()--> ResolveUDP error")
                internalChan.setupfail <- true
        }
        sock, err := net.ListenUDP("udp4", addr)
        if err != nil {
                printc.DataWithColor(printc.COLOR_RED,"network.IMAListen()--> ListenUDP error")
                internalChan.setupfail <- true
        }
        var data [512]byte
        for {
                select {
                case <-internalChan.quitImaListen:
                        return
                default:
                        _, remoteAddr, err := sock.ReadFromUDP(data[0:])
                        if err != nil {
                                printc.DataWithColor(printc.COLOR_RED,"network.IMAListen()--> ReadFromUDP error")
                                break
                        }
                        if LOCALIP != remoteAddr.IP.String() {
                                if err == nil {
                                        elevIP := remoteAddr.IP.String()
                                        internalChan.ima <- elevIP
                                } else {
                                        printc.DataWithColor(printc.COLOR_RED,"network.IMAListen()--> UDP read error")
                                } 
                        } 
                } 
        } 
} 

func imaWatcher() {
        peers := make(map[string]time.Time)
        deadline := IMALOSS * IMAPERIOD * time.Millisecond
        for {
                //printc.DataWithColor(printc.COLOR_CYAN,"ROUND ROUND GET AROUND I GET AROUND")
                select {
                case ip := <-internalChan.ima:
                        _, inMap := peers[ip]
                        if inMap {
                                peers[ip] = time.Now()
                        } else {
                                peers[ip] = time.Now()
                                internalChan.newIP <- ip
                        }
                case <-time.After(ALIVEWATCH * time.Millisecond):
                        for ip, timestamp := range peers {
                                if time.Now().After(timestamp.Add(deadline)) {
                                        printc.DataWithColor(printc.COLOR_RED, "network.imaWatcher --> Timeout", ip)
                                        externalChan.GetDeadElevator <- ip
                                        internalChan.closeConn <- ip
                                        delete(peers, ip)
                                }
                        }
                case deadIP := <-externalChan.SendDeadElevator:
                        internalChan.closeConn <- deadIP
                        delete(peers, deadIP)
                case errorIP := <-internalChan.errorIP:
                        printc.DataWithColor(printc.COLOR_CYAN,"I DID PASS!!!")
                        _, inMap := peers[errorIP]
                        if inMap {
                                printc.DataWithColor(printc.COLOR_CYAN,"You are giving me a panicattack!! FUCK YOU")
                                externalChan.Panic <- true
                        }
                case <-internalChan.quitImaWatcher:
                        return
                }
        }
}

func imaSend() {
        service := bcast + ":" + UDPport
        addr, err := net.ResolveUDPAddr("udp4", service)
        if err != nil {
                printc.DataWithColor(printc.COLOR_RED,"network.IMASend()--> Resolve error")
                internalChan.setupfail <- true
        }
        imaSock, err := net.DialUDP("udp4", nil, addr)
        if err != nil {
                printc.DataWithColor(printc.COLOR_RED,"network.IMASend()--> Dial error")
                internalChan.setupfail <- true
        }
        ima := []byte("IMA")
        for {
                select {
                case <-internalChan.quitImaSend:
                        return
                default:
                        _, err := imaSock.Write(ima)
                        if err != nil {
                                printc.DataWithColor(printc.COLOR_RED,"network.IMASend()--> UDP send error")
                        }
                        time.Sleep(IMAPERIOD * time.Millisecond)
                }
        }
}