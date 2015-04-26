package network_module

import (
	"net"
	"time"
	"printc"
)

func ima_listen() {
	service := broad_cast + ":" + UDP_port
	addr, err := net.ResolveUDPAddr("udp4", service)
	if err != nil {
		printc.Data_with_color(printc.COLOR_RED,"network.IMAListen()--> ResolveUDP error")
		internal_chan.setup_fail <- true
	}
	sock, err := net.ListenUDP("udp4", addr)
	if err != nil {
		printc.Data_with_color(printc.COLOR_RED,"network.IMAListen()--> ListenUDP error")
		internal_chan.setup_fail <- true
	}
	var data [512]byte
	for {
		select {
		case <-internal_chan.quit_ima_listen:
			return
		default:
			_, remoteAddr, err := sock.ReadFromUDP(data[0:])
			if err != nil {
				printc.Data_with_color(printc.COLOR_RED,"network.IMAListen()--> ReadFromUDP error")
				break
			}
			if LOCAL_IP != remoteAddr.IP.String() {
				if err == nil {
					elevIP := remoteAddr.IP.String()
					internal_chan.ima <- elevIP
				} else {
					printc.Data_with_color(printc.COLOR_RED,"network.IMAListen()--> UDP read error")
				} 
			} 
		} 
	} 
} 

func ima_watcher() {
	peers := make(map[string]time.Time)
	deadline := IMA_LOSS * IMA_PERIOD * time.Millisecond
	for {
		select {
		case ip := <-internal_chan.ima:
			_, inMap := peers[ip]
			if inMap {
				peers[ip] = time.Now()
			} else {
				peers[ip] = time.Now()
				internal_chan.new_IP <- ip
			}
		case <-time.After(ALIVE_WATCH * time.Millisecond):
			for ip, timestamp := range peers {
				if time.Now().After(timestamp.Add(deadline)) {
					printc.Data_with_color(printc.COLOR_RED, "network.imaWatcher --> Timeout", ip)
					external_chan.Get_dead_elevator <- ip
					internal_chan.close_conn <- ip
					delete(peers, ip)
				}
			}
		case <-internal_chan.quit_ima_watcher:
			return
		}
	}
}

func ima_send() {
	service := broad_cast + ":" + UDP_port
	addr, err := net.ResolveUDPAddr("udp4", service)
	if err != nil {
		printc.Data_with_color(printc.COLOR_RED,"network.IMASend()--> Resolve error")
		internal_chan.setup_fail <- true
	}
	imaSock, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		printc.Data_with_color(printc.COLOR_RED,"network.IMASend()--> Dial error")
		internal_chan.setup_fail <- true
	}
	ima := []byte("IMA")
	for {
		select {
		case <-internal_chan.quit_ima_send:
			return
		default:
			_, err := imaSock.Write(ima)
			if err != nil {
				printc.Data_with_color(printc.COLOR_RED,"network.IMASend()--> UDP send error")
			}
			time.Sleep(IMA_PERIOD * time.Millisecond)
		}
	}
}