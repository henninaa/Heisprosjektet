package network_module

type internal_channels struct {
	setup_fail       chan bool
	ima             chan string
	new_IP           chan string
	dead_elevator    chan string
	update_TCP_map    chan tcp_connection
	connect_fail     chan string
	close_conn       chan string
	error_IP         chan string
	connection_error chan string
	quit_ima_send     chan bool
	quit_ima_listen   chan bool
	quit_ima_watcher  chan bool
	quit_listen_TCP   chan bool
	quit_TCP_map      chan bool
}

type Net_channels struct {
	Get_dead_elevator  chan string
	Send_dead_elevator chan string
	Send_to_all        chan Mail
	Send_to_one        chan Mail
	Inbox              chan Mail
	Num_of_peers       chan int
	New_connection     chan string

}

var internal_chan internal_channels
var external_chan Net_channels

func (internal_chan *internal_channels) init() {
	internal_chan.setup_fail = make(chan bool)
	internal_chan.ima = make(chan string)
	internal_chan.new_IP = make(chan string)
	internal_chan.dead_elevator = make(chan string)
	internal_chan.update_TCP_map = make(chan tcp_connection)
	internal_chan.connect_fail = make(chan string)
	internal_chan.connection_error = make(chan string)
	internal_chan.error_IP = make(chan string)
	internal_chan.close_conn = make(chan string)
	internal_chan.quit_ima_send = make(chan bool)
	internal_chan.quit_ima_listen = make(chan bool)
	internal_chan.quit_ima_watcher = make(chan bool)
	internal_chan.quit_listen_TCP = make(chan bool)
	internal_chan.quit_TCP_map = make(chan bool)
}

func (external_chan *Net_channels) Init() {
	external_chan.Get_dead_elevator = make(chan string)
	external_chan.Send_to_all = make(chan Mail)
	external_chan.Send_to_one = make(chan Mail)
	external_chan.Inbox = make(chan Mail)
	external_chan.Num_of_peers = make(chan int)
	external_chan.New_connection = make(chan string)

}