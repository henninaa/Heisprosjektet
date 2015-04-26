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

func manage_TCP_connections() {
	connections := conn_map{make(map[string]conn_chans)}
	go listen_for_TCP_connections()
	for {
		select {
		case new_IP := <-internal_chan.new_IP:
			connections.handle_new_IP(new_IP)

		case new_TCP_connection := <-internal_chan.update_TCP_map:
			connections.handle_new_connection(new_TCP_connection)

		case error_IP := <-internal_chan.connect_fail:
			connections.handle_failed_to_connect(error_IP)

		case error_IP := <-internal_chan.connection_error:
			connections.handle_connection_error(error_IP)

		case close_IP := <-internal_chan.close_conn:
			connections.handle_close_connection(close_IP)

		case mail := <-external_chan.Send_to_all:
			connections.handle_send_to_all(mail)

		case mail := <-external_chan.Send_to_one:
			connections.handle_send_to_one(mail)

		case <-internal_chan.quit_TCP_map:
			return
		} 
	} 
} 

func (connections *conn_map) handle_new_IP(new_IP string) {
	_, in_Map := connections.tcp_map[new_IP]
	if !in_Map {
		go connect_TCP(new_IP)
	} else {
		printc.Data_with_color(printc.COLOR_YELLOW,"network.monitorTCPConnections-->", new_IP, "already in connections")
	}
}

func (connections *conn_map) handle_new_connection(conn tcp_connection) {
	_, in_Map := connections.tcp_map[conn.ip]
	if !in_Map {
		connections.tcp_map[conn.ip] = conn_chans{send: make(chan Mail), quit: make(chan bool)} 
		printc.Data_with_color(printc.COLOR_GREEN,"network.monitorTCPConnections---> Connection made to ", conn.ip)
		conn.sendChan = connections.tcp_map[conn.ip].send
		conn.quit = connections.tcp_map[conn.ip].quit
		external_chan.New_connection <- conn.ip
		go conn.handle_connection()
		go peer_update(len(connections.tcp_map))
	} else {
		printc.Data_with_color(printc.COLOR_YELLOW,"network.monitorTCPConnections--> A connection already exist to", conn.ip)
		conn.socket.Close()
	}
}

func (connections *conn_map) handle_failed_to_connect(error_IP string) {
	_, in_Map := connections.tcp_map[error_IP]
	if in_Map {
		printc.Data_with_color(printc.COLOR_YELLOW,"network.monitorTCPConnections--> Could not dial up ", error_IP, "but a connection already exist")
	} else {
		printc.Data_with_color(printc.COLOR_RED,"network.monitorTCPConnections--> Could not connect to ", error_IP)
		internal_chan.error_IP <- error_IP //Notify imaWatcher of erroneous ip. Maybe it has timed out?
		printc.Data_with_color(printc.COLOR_CYAN,"YOU SHALL NOT PASS!!!!")
	}
}

func (connections *conn_map) handle_connection_error(error_IP string) {
	_, in_Map := connections.tcp_map[error_IP]
	if in_Map {
		delete(connections.tcp_map, error_IP)
	}
	go connect_TCP(error_IP)
}

func (connections *conn_map) handle_close_connection(close_IP string) {
	conn_chans, in_Map := connections.tcp_map[close_IP]
	printc.Data_with_color(printc.COLOR_CYAN,"conn_chans ", conn_chans, "in_Map ", in_Map)
	if in_Map {
		select {
		case conn_chans.quit <- true:
		case <-time.After(10 * time.Millisecond):
		}
		delete(connections.tcp_map, close_IP)
		num_of_conns := len(connections.tcp_map)
		if num_of_conns == 0 {
			go peer_update(num_of_conns)
		}
	} else {
		printc.Data_with_color(printc.COLOR_YELLOW,"network.monitorTCPConnections--> No connection to close ", close_IP)

	}
}

func (connections *conn_map) handle_send_to_one(mail Mail) {
	printc.Data_with_color(printc.COLOR_BLUE,"IP: ", mail.IP)
	switch mail.IP {
	case "":
		size := len(connections.tcp_map)
		if size != 0 {
			for _, conn_chans := range connections.tcp_map {
				conn_chans.send <- mail 
				break
			}
		}
	default:
		conn_chans, in_Map := connections.tcp_map[mail.IP]
		if in_Map {
			conn_chans.send <- mail
		} else {
			internal_chan.error_IP <- mail.IP
		}
	}
}

func (connections *conn_map) handle_send_to_all(mail Mail) {
	if len(connections.tcp_map) != 0 {
		for _, conn_chans := range connections.tcp_map {
			conn_chans.send <- mail
		}
	}
}

func (conn *tcp_connection) handle_connection() {
	quit_inbox := make(chan bool)
	go conn.inbox(quit_inbox)
	printc.Data_with_color(printc.COLOR_BLUE,"Network.handle_connection--> handle_connection for", conn.ip, "is running")
	connection_encoder := gob.NewEncoder(conn.socket)
	for {
		select {
		case mail := <-conn.sendChan:
			encoded_msg := mail.Msg
			err := connection_encoder.Encode(&encoded_msg)
			if err == nil {
				printc.Data_with_color(printc.COLOR_CYAN, "Message sendt without problem :)")
			} else {
				printc.Data_with_color(printc.COLOR_RED,"Network.handle_connection--> Error sending message to ", conn.ip, err)
				internal_chan.connection_error <- conn.ip 
			}
		case <-conn.quit:
			conn.socket.Close()
			printc.Data_with_color(printc.COLOR_YELLOW,"Network.handle_connections--> Connection to ", conn.ip, " has been terminated.")
			return
		case <-quit_inbox:
			conn.socket.Close()
			printc.Data_with_color(printc.COLOR_YELLOW,"Network.handle_connections--> Connection to ", conn.ip, " has been terminated.")
			internal_chan.connection_error <- conn.ip
			return
		}
	}
}

func (conn *tcp_connection) inbox(quit_inbox chan bool) {
	connection_decoder := gob.NewDecoder(conn.socket)
	for {
		decoded_msg := new(Message)
		err := connection_decoder.Decode(decoded_msg)
		switch err {
		case nil:
			new_mail := Mail{IP: conn.ip, Msg: *decoded_msg}
			external_chan.Inbox <- new_mail
		default:
			printc.Data_with_color(printc.COLOR_RED,"Network.inbox--> Error:", err)
			time.Sleep(IMA_PERIOD * IMA_LOSS * 2 * time.Millisecond)
			select {
			case quit_inbox <- true: 
			case <-time.After(WRITE_DL * time.Millisecond):
			}
			return
		}
	}
}

func connect_TCP(ip string) {
	attempts := 0
	for attempts < CONN_ATMPT {
		printc.Data_with_color(printc.COLOR_YELLOW,"Network.connect_TCP--> attempting to connect to ", ip)
		service := ip + ":" + TCP_port
		_, err := net.ResolveTCPAddr("tcp4", service)
		if err != nil {
			printc.Data_with_color(printc.COLOR_RED,"Network.connect_TCP--> ResolveTCPAddr failed")
			attempts++
			time.Sleep(DIAL_INT * time.Millisecond)
		} else {
			rand_sleep := time.Duration(rand.Intn(500)+500) * time.Microsecond
			printc.Data_with_color(printc.COLOR_MAGENTA,"Network.connect_TCP--> rand_sleep:", rand_sleep)
			time.Sleep(rand_sleep)
			socket, err := net.Dial("tcp4", service)
			if err != nil {
				printc.Data_with_color(printc.COLOR_RED,"Network.connect_TCP--> DialTCP error when connecting to", ip, " error: ", err)
				attempts++
				time.Sleep(DIAL_INT * time.Millisecond)
			} else {
				new_TCP_connection := tcp_connection{ip: ip, socket: socket}
				internal_chan.update_TCP_map <- new_TCP_connection 
				break
			}
		}
	}
	if attempts == CONN_ATMPT {
		select {
		case internal_chan.connect_fail <- ip: 
		case <-time.After(CONN_FAIL_TIMEOUT * time.Millisecond): 
			return
		}
	}
}

func listen_for_TCP_connections() {
	service := ":" + TCP_port
	tcp_addr, err := net.ResolveTCPAddr("tcp4", service)
	if err != nil {
		printc.Data_with_color(printc.COLOR_RED,"Network.listen_for_TCP_connections--> TCP resolve error")
		internal_chan.setup_fail <- true
	} else {
		listen_sock, err := net.ListenTCP("tcp4", tcp_addr)
		if err != nil {
			printc.Data_with_color(printc.COLOR_RED,"Network.connect_TCP--> ListenTCP error")
			internal_chan.setup_fail <- true
		} else {
			printc.Data_with_color(printc.COLOR_YELLOW,"Network.connect_TCP--> listening for new connections")
			for {
				select {
				case <-internal_chan.quit_listen_TCP:
					return
				default:
					socket, err := listen_sock.Accept()
					if err == nil {
						ip := clean_up_IP(socket.RemoteAddr().String())
						new_TCP_connection := tcp_connection{ip: ip, socket: socket}
						internal_chan.update_TCP_map <- new_TCP_connection
					}
				}
			}
		}
	}
}

func peer_update(num_of_peers int) {
	select {
	case external_chan.Num_of_peers <- num_of_peers:
	case <-time.After(500 * time.Millisecond):
	}
}

func clean_up_IP(garbage string) (cleanIP string) {
	split := strings.Split(garbage, ":") 
	cleanIP = split[0]
	return
}