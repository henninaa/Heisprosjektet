package bank_module

import(
	"time"
	"queue_module"
	"network_module"
	"driver_module"
	"sensor_module"
	"FSM_module"
	"fmt"
	"printc"
	)


func Elevator_main_control(){

	var elevator elevator_type

	elevator.elevator_type_init()

	var queue queue_module.Queue_type
	queue.Init()

	var(
		sensor_channels sensor_module.External_channels
		event_channels FSM_module.External_channels
		network_channels network_module.NetChannels
		external_chan External_channels
		internal_chan internal_channels
		)

	internal_chan.init()
	external_chan.init()
	sensor_channels.Init()
	event_channels.Init()
	network_channels.NetChanInit()


	start_up(sensor_channels, event_channels, network_channels)
	fmt.Println("clickibng: ", sensor_channels.Order_chan)
	//var mail_buffer []network_module.Mail

	//----------Main control loop

	for{

		printc.DataWithColor(printc.COLOR_YELLOW, "ny runde")
		time.Sleep(ELEVATOR_MAIN_CONTROL_INTERVAL)
		select{

			case msg := <- network_channels.Inbox:
				printc.DataWithColor(printc.COLOR_GREEN, "inbox")
				handle_network_messgage(msg, internal_chan)

			case post := <- sensor_channels.Order_chan:

				printc.DataWithColor(printc.COLOR_GREEN, "order_chan")

				handle_new_order(post,internal_chan.insert_to_queue, internal_chan.auction_order)

			case post := <- internal_chan.insert_to_queue:

				printc.DataWithColor(printc.COLOR_GREEN, "insert_to_queue")

				insert_post_to_queue(post, &queue, elevator.floor, external_chan.request_new_direction)
				
			case post := <- internal_chan.auction_order:

				printc.DataWithColor(printc.COLOR_GREEN, "auction_order")

				deliver_order(post, queue.Get_lowest_cost_ip(post, elevator.floor), internal_chan.insert_to_queue, network_channels.SendToOne)

			case mail := <- internal_chan.order_executed:

				handle_order_executed(mail, &queue)

				printc.DataWithColor(printc.COLOR_GREEN, "order_executed")

			case floor := <- sensor_channels.Floor_chan:

				handle_new_floor(floor, &(elevator.floor), event_channels.Reached_floor, internal_chan.check_stop_conditions)

				printc.DataWithColor(printc.COLOR_GREEN, "new fllor")

			case <- external_chan.request_new_direction:
				
				printc.DataWithColor(printc.COLOR_GREEN, "req_new_direction")
				handle_new_direction(queue.Get_new_direction(elevator.floor), &elevator, event_channels.New_direction)

			case direction := <- internal_chan.new_direction:

				printc.DataWithColor(printc.COLOR_GREEN, "new_direction")
				handle_new_direction(direction, &elevator, event_channels.New_direction)

			case ip := <- network_channels.New_connection:
				printc.DataWithColor(printc.COLOR_GREEN, "New_connection")

				handle_new_connection(ip, queue, elevator.floor)

			case <- internal_chan.check_stop_conditions:
				printc.DataWithColor(printc.COLOR_GREEN, "check_stop")

				check_stop_cond(&queue, &elevator, event_channels.Stop)

			case <- event_channels.Req_direction:

				printc.DataWithColor(printc.COLOR_GREEN, "Req_new_direction")
				handle_new_direction(queue.Get_new_direction(elevator.floor), &elevator, event_channels.New_direction)



		}

	}

}

func handle_network_messgage(mail network_module.Mail, internal_chan internal_channels){

	switch (mail.Msg.Msg_type){

	case network_module.ORDER_TAKEN:

	case network_module.ORDER_EXECUTED:

		internal_chan.order_executed <- mail

	case network_module.DELIVER_ORDER:

	case network_module.TAKE_NEW_ORDER:

		internal_chan.insert_to_queue <- queue_module.Convert_mail_to_queue_post(mail)

	case network_module.TAKE_BACKUP_ORDER:

	case network_module.BACKUP_ORDER_COMPLETE:

		internal_chan.take_backup_order <- convert_mail_to_backup_post(mail, mail.IP)

	case network_module.ERROR_MSG:

	case network_module.TAKE_BACKUP_FLOOR:

	}
}

func handle_new_order(post queue_module.Queue_post, insert_to_queue chan queue_module.Queue_post, auction_order chan queue_module.Queue_post){

	if (post.Button_type == driver_module.BUTTON_COMMAND){
		insert_to_queue <- post
	} else{

		auction_order <- post

	}

	fmt.Println("yeye")

}

func deliver_order(post queue_module.Queue_post, IP string, insert_to_queue chan queue_module.Queue_post, send_to_one chan network_module.Mail){

	var mail network_module.Mail

	if(IP == "self"){
		insert_to_queue <- post
	}else{

		mail.Make_mail()
		send_to_one <- mail
	}

}

func insert_post_to_queue(post queue_module.Queue_post, queue * queue_module.Queue_type, current_floor int, req_dir_chan chan int){

	var mail network_module.Mail

	mail.Make_mail()
	queue.Insert_to_own_queue(post, current_floor)

	req_dir_chan <- 1

}

func handle_new_floor(floor_input int, current_floor * int, new_floor_event_channel chan int, stop_check_chan chan int){

	if(floor_input!=-1){

		*current_floor = floor_input
		new_floor_event_channel <- floor_input
		stop_check_chan <- 1

	}
}

func handle_new_direction(direction int, elevator * elevator_type, new_dir_chan chan int){

	if (direction == -1){
		elevator.moving = false
	} else if (direction == driver_module.UP){
		elevator.moving = true
		elevator.direction = driver_module.UP
		new_dir_chan <- direction
	} else{
		elevator.moving = true
		elevator.direction = driver_module.DOWN
		new_dir_chan <- direction
	}

	printc.DataWithColor(printc.COLOR_BLUE, "for")
	new_dir_chan <- direction
	printc.DataWithColor(printc.COLOR_BLUE, "etter")

}

func handle_new_connection(ip string, queue queue_module.Queue_type, current_floor int){


}

func start_up(sensor_channels sensor_module.External_channels, event_channels FSM_module.External_channels, network_channels network_module.NetChannels ){



	sensor_module.Sensors(sensor_channels)
	go FSM_module.Event_generator(event_channels)
    network_module.NetworkSetup(network_channels)
}

func convert_mail_to_backup_post(mail network_module.Mail, ip string)(backup_post queue_backup_post){

	backup_post.post = queue_module.Convert_mail_to_queue_post(mail)
	backup_post.IP = ip

	return backup_post
}

func handle_order_executed(mail network_module.Mail, queue * queue_module.Queue_type){

	post := queue_module.Convert_mail_to_queue_post(mail)

	queue.Remove_post_from_backup_queue(post, mail.IP)

}

func check_stop_cond(queue * queue_module.Queue_type, elevator * elevator_type, stop_event_chan chan int){
		fmt.Println(elevator.floor)
	if(queue.Should_elevator_stop(elevator.floor, driver_module.Convert_dir_to_button(elevator.direction))){

		stop_event_chan <- 1
		elevator.direction = -1
		elevator.moving = false
	}

}