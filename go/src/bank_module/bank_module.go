package bank_module

import(
	"time"
	"queue_module"
	"network_module"
	"driver_module"
	"sensor_module"
	"FSM_module"
	)


func Elevator_main_control(){

	var elevator elevator_type
	elevator.elevator_type_init()

	var queue queue_module.Queue_type

	var(
		sensor_channels sensor_module.External_channels
		event_channels FSM_module.External_channels
		network_channels network_module.NetChannels
		external_chan External_channels
		internal_chan internal_channels
		)

	start_up(sensor_channels, event_channels, network_channels)
	
	var mail_buffer []network_module.Mail

	//----------Main control loop

	for{

		time.Sleep(ELEVATOR_MAIN_CONTROL_INTERVAL)
		select{

			case msg := <- network_channels.Inbox:

				handle_network_messgage(msg)

			case post := <- external_chan.new_order:

				handle_new_order(post)

			case post := <- internal_chan.insert_to_queue:

				insert_post_to_queue(post, elevator.floor)
				
			case post := <- internal_chan.auction_order:

				deliver_order(post, queue.Get_lowest_cost_ip(post, elevator.floor))

			case floor := <- external_chan.new_floor:

				handle_new_floor(floor, &(elevator.floor), event_channels.Reached_floor)

			case <- external_chan.request_new_direction:

				handle_new_direction(queue.Get_new_direction(elevator.floor))

			case direction := <- internal_chan.new_direction:

				handle_new_direction(direction)

			case ip := <- network_channels.New_connection:

				handle_new_connection(ip, queue.queue, elevator.floor)

		}

	}

}

func handle_network_messgage(mail network_module.Mail){

	switch (mail.Msg.msg_type){

	case network_module.ORDER_TAKEN:

	case network_module.ORDER_EXECUTED:

	case network_module.DELIVER_ORDER:

	case network_module.TAKE_NEW_ORDER:

		insert_post_to_queue(queue_module.Convert_mail_to_post(mail))

	case network_module.TAKE_BACKUP_ORDER:

	case network_module.BACKUP_ORDER_COMPLETE:

		internChan.take_backup_order <- queue_module.Convert_mail_to_backup_post(mail)

	case network_module.ERROR_MSG:

		idk.com

	case network_module.TAKE_BACKUP_FLOOR:

	}
}

func handle_new_order(post queue_module.Queue_post){

	if (post.button_type == driver_module.BUTTON_COMMAND){
		internChan.insert_to_queue <- post
	} else{

		internChan.auction_order <- post

	}

}

func deliver_order(post queue_module.Queue_post, IP string){

	var mail network_module.Mail

	if(IP == "self"){
		internChan.insert_to_queue <- post
	}else{

		mail.make_mail(IP, network_module.DELIVER_ORDER, JSON.Marshal(post))
		network_module.externalChan.send_to_one <- mail
	}

}

func insert_post_to_queue(post queue_module.Queue_post, floor int){

	var mail network_module.Mail

	mail.Make_mail("", network_module.ORDER_TAKEN, post.floor, driver_module.Button_type_to_int(post.button_type))
	queue.insert_queue(post.queue, post.button_type, elevator.floor)
}

func handle_new_floor(floor_input int, current_floor * int, new_floor_event_channel chan int){

	current_floor = floor_input
	new_floor_event_channel <- floor_iput
}

func handle_new_direction(direction int){


}

func handle_new_connection(ip string, ){


}

func start_up(sensor_channels sensor_module.External_channels, event_channels FSM_module.External_channels, network_channels network_module.NetChannels ){

	sensor_channels.sensor_external_channels_init()
	event_channels.event_external_channels_init()
	network_channels.netChanInit()

	go Sensors(sensor_channels)
	go Event_generator(event_channels)
	go network_module(network_channels)
}