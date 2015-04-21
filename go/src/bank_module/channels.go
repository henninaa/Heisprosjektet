package bank_module

type internalChannels struct{

	insert_to_queue chan queue_module.queue_post
	take_backup_order chan queue_module.backup_post
	auction_order chan queue_module.queue_post
}

func (internChan * internalChannels) internal_channels_init(){

	internChan.insert_to_queue = make(chan queue_module.queue_post)
	internChan.take_backup_order = make(chan queue_module.backup_post)
	internChan.auction_order = make(chan queue_module.queue_post)

}

type ExternalChannels struct{

	new_order chan queue_module.queue_post
	request_new_direction chan int
}

func (externalChan * ExternalChannels) external_channels_init(){

	externalChan.new_order = make(chan queue_module.queue_post)
	externalChan.request_new_direction = make (chan int)

}