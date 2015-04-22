package bank_module

import(
	"queue_module"
	"network_module"
	)

type internal_channels struct{

	insert_to_queue chan queue_module.Queue_post
	take_backup_order chan queue_backup_post
	auction_order chan queue_module.Queue_post
	new_direction chan int
	order_executed chan network_module.Mail
	check_stop_conditions chan int
}

func (intern_chan * internal_channels) init(){

	intern_chan.insert_to_queue = make(chan queue_module.Queue_post,2)
	intern_chan.take_backup_order = make(chan queue_backup_post,2)
	intern_chan.auction_order = make(chan queue_module.Queue_post,2)
	intern_chan.new_direction = make(chan int,2)
	intern_chan.order_executed = make(chan network_module.Mail,2)
	intern_chan.check_stop_conditions = make(chan int,2)

}

type External_channels struct{

	new_order chan queue_module.Queue_post
	request_new_direction chan int
	new_floor chan int
}

func (external_chan * External_channels) init(){

	external_chan.new_order = make(chan queue_module.Queue_post, 2)
	external_chan.request_new_direction = make (chan int, 2)
	external_chan.new_floor = make(chan int,2 )

}