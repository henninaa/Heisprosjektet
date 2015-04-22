package bank_module

import("queue_module")

type internal_channels struct{

	insert_to_queue chan queue_module.Queue_post
	take_backup_order chan queue_backup_post
	auction_order chan queue_module.Queue_post
	new_direction chan int
}

func (internChan * internalChannels) internal_channels_init(){

	internChan.insert_to_queue = make(chan queue_module.Queue_post)
	internChan.take_backup_order = make(chan queue_module.Backup_post)
	internChan.auction_order = make(chan queue_module.Queue_post)
	internChan.new_direction = make(chan int)

}

type External_channels struct{

	new_order chan queue_module.Queue_post
	request_new_direction chan int
	new_floor chan int
}

func (externalChan * ExternalChannels) external_channels_init(){

	externalChan.new_order = make(chan queue_module.Queue_post)
	externalChan.request_new_direction = make (chan int)
	externalChan.new_floor = make(chan int)

}