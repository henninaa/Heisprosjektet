package network_module

import (

        "math/rand"
        "time"
        "printc"
        "driver_module"
)


func Network_setup(Net_chan Net_channels) {
        rand.Seed(time.Now().UTC().UnixNano())
        
        internal_chan.init()
        
        external_chan = Net_chan
        
        ima_start()
        network_start()
}

func ima_start() {
        go ima_watcher()
        go ima_listen()
        go ima_send()
}

func network_start() {
        go manage_TCP_connections()
        
        for {
                select {
                
                case <-internal_chan.setup_fail:
                        printc.DataWithColor(printc.COLOR_RED,"net.Startup--> Setupfail. Retrying...")
                        
                        internal_chan.quit_ima_send <- true
                        internal_chan.quit_ima_listen <- true
                        internal_chan.quit_ima_watcher <- true
                        internal_chan.quit_listen_TCP <- true
                        internal_chan.quit_TCP_map <- true
                        time.Sleep(time.Millisecond)
                        
                        ima_start()
                        
                        go manage_TCP_connections()
                
                case <-time.After(NET_SETUP * time.Millisecond):
                        printc.DataWithColor(printc.COLOR_GREEN,"net.Startup --> Network setup complete")
                        return
                }
        }
}

func (mail * Mail) Make_mail(IP string, msg_type string, floor int, dir driver_module.Elev_button_type_t, cost int){
        
        mail.IP = IP
        mail.Msg.Msg_type = msg_type
        mail.Msg.Floor = floor
        mail.Msg.Dir = dir
        mail.Msg.Cost_answer = cost
}