package network_module

import (

        "math/rand"
        "time"
        "printc"
)


func NetworkSetup(NetChan NetChannels) {
        rand.Seed(time.Now().UTC().UnixNano())
        
        internalChan.init()
        
        externalChan = NetChan
        
        imaStart()
        networkStart()
}

func imaStart() {
        go imaWatcher()
        go imaListen()
        go imaSend()
}

func networkStart() {
        go manageTCPConnections()
        
        for {
                select {
                
                case <-internalChan.setupfail:
                        printc.DataWithColor(printc.COLOR_RED,"net.Startup--> Setupfail. Retrying...")
                        
                        internalChan.quitImaSend <- true
                        internalChan.quitImaListen <- true
                        internalChan.quitImaWatcher <- true
                        internalChan.quitListenTCP <- true
                        internalChan.quitTCPMap <- true
                        time.Sleep(time.Millisecond)
                        
                        imaStart()
                        
                        go manageTCPConnections()
                
                case <-time.After(NETSETUP * time.Millisecond):
                        printc.DataWithColor(printc.COLOR_GREEN,"net.Startup --> Network setup complete")
                        return
                }
        }
}