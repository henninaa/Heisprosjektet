package bank_module

import(
	"time"
	"queue_module"
	"network_module"
	"driver_module"
	"sensor_module"
	"FSM_module"
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
		network_channels network_module.Net_channels
		external_chan External_channels
		internal_chan internal_channels
		)

	internal_chan.init()
	external_chan.init()
	sensor_channels.Init()
	event_channels.Init()
	network_channels.Init()


	start_up(sensor_channels, event_channels, network_channels)
	//var mail_buffer []network_module.Mail

	//----------Main control loop

	go_to_defined_floor(&elevator, sensor_channels)

	for{
		time.Sleep(ELEVATOR_MAIN_CONTROL_INTERVAL)
		select{

			case msg := <- network_channels.Inbox:
				printc.Data_with_color(printc.COLOR_GREEN, "Getting a message from network_channels.Inbox")
				handle_network_messgage(msg, internal_chan)

			
			case post := <- sensor_channels.Order_chan:
				printc.Data_with_color(printc.COLOR_GREEN, "Getting a message from sensor_channels.Order_chan")
				handle_new_order(post,internal_chan.insert_to_queue, internal_chan.auction_order)

			
			case post := <- internal_chan.insert_to_queue:
				printc.Data_with_color(printc.COLOR_GREEN, "Getting a message from internal_chan.insert_to_queue")
				insert_post_to_queue(post, &queue, elevator.floor, network_channels.Send_to_all, event_channels)
				
			
			case post := <- internal_chan.auction_order:
				printc.Data_with_color(printc.COLOR_GREEN, "Getting a message from internal_chan.auction_order")
				deliver_order(post, queue.Get_lowest_cost_ip(post, elevator.floor), internal_chan.insert_to_queue, network_channels.Send_to_one)

			
			case mail := <- internal_chan.order_executed:
				handle_order_executed(mail, &queue)
				printc.Data_with_color(printc.COLOR_GREEN, "Getting a message from internal_chan.order_executed")

			
			case floor := <- sensor_channels.Floor_chan:
				printc.Data_with_color(printc.COLOR_GREEN, "Getting a message from sensor_channels.Floor_chan")
				handle_new_floor(floor, &(elevator.floor), internal_chan.check_stop_conditions)

			
			case direction := <- internal_chan.new_direction:
				printc.Data_with_color(printc.COLOR_GREEN, "Getting a message from internal_chan.new_direction")
				handle_new_direction(direction, &elevator, event_channels)

			
			case ip := <- network_channels.New_connection:
				printc.Data_with_color(printc.COLOR_GREEN, "Getting a message from network_channels.New_connection")
				handle_new_connection(ip, queue, elevator.floor)

			
			case <- internal_chan.check_stop_conditions:
				printc.Data_with_color(printc.COLOR_GREEN, "Getting a message from internal_chan.check_stop")
				check_stop_cond(&queue, &elevator, event_channels, network_channels)

			
			case <- event_channels.Get_new_direction:
				printc.Data_with_color(printc.COLOR_GREEN, "Getting a message from event_channels.Get_new_direction")
				handle_new_direction(queue.Get_new_direction(elevator.floor), &elevator, event_channels)

			
			case <- event_channels.Get_new_action:
				printc.Data_with_color(printc.COLOR_GREEN, "Getting a message from event_channels.Get_new_action")
				make_new_action(&queue, elevator, event_channels, network_channels)

			
			case dead_elevator := <- network_channels.Get_dead_elevator:
				printc.Data_with_color(printc.COLOR_RED, "Getting a message from network_channels.Get_dead_elevator: ", dead_elevator)

			
			case <- internal_chan.take_backup_order:

			
			case <- internal_chan.take_backup_order:
/*
			
			case <- event_channels.Get_new_direction:

				printc.Data_with_color(printc.COLOR_GREEN, "Req_new_direction")
				handle_new_direction(queue.Get_new_direction(elevator.floor), &elevator, event_channels)
*/




		}

	}

}

func handle_network_messgage(mail network_module.Mail, internal_chan internal_channels){

	switch (mail.Msg.Msg_type){

	case network_module.ORDER_EXECUTED:

		internal_chan.order_executed <- mail
		printc.Data_with_color(printc.COLOR_GREEN, "ORDER_EXECUTED.MAIL ", mail)

	case network_module.DELIVER_ORDER:

	case network_module.TAKE_NEW_ORDER:

		internal_chan.insert_to_queue <- queue_module.Convert_mail_to_queue_post(mail)
		printc.Data_with_color(printc.COLOR_GREEN, "TAKE_NEW_ORDER.MAIL ", mail)

	case network_module.TAKE_BACKUP_ORDER:

		internal_chan.take_backup_order <- convert_mail_to_backup_post(mail, mail.IP)
		printc.Data_with_color(printc.COLOR_GREEN, "TAKE_NEW_BACKUP_ORDER.MAIL ", mail)

	case network_module.TAKE_BACKUP_FLOOR:



	}
}

func handle_new_order(post queue_module.Queue_post, insert_to_queue chan queue_module.Queue_post, auction_order chan queue_module.Queue_post){

	if (post.Button_type == driver_module.BUTTON_COMMAND){
		insert_to_queue <- post
	} else{

		auction_order <- post

	}

}

func deliver_order(post queue_module.Queue_post, IP string, insert_to_queue chan queue_module.Queue_post, send_to_one chan network_module.Mail){

	var mail network_module.Mail

	if(IP == "self"){
		insert_to_queue <- post
	}else{

		mail.Make_mail(IP, network_module.TAKE_NEW_ORDER, post.Floor, post.Button_type, 0)
		send_to_one <- mail
	}

}

func insert_post_to_queue(post queue_module.Queue_post, queue * queue_module.Queue_type, current_floor int,send_to_all chan network_module.Mail, event_channels FSM_module.External_channels){

	var mail network_module.Mail

	mail.Make_mail("", network_module.TAKE_BACKUP_ORDER, post.Floor, post.Button_type, 0)
	send_to_all <- mail

	queue.Insert_to_own_queue(post, current_floor)

	event_channels.New_order <- 1

}

func handle_new_floor(floor_input int, current_floor * int, stop_check_chan chan int){

	if(floor_input!=-1){

		*current_floor = floor_input
		stop_check_chan <- 1

	}
}

func handle_new_direction(direction int, elevator * elevator_type, event_chan FSM_module.External_channels){

	printc.Data_with_color(printc.COLOR_YELLOW, "Retning: ", direction)

	if (direction == -1){
		elevator.moving = false
	} else if (direction == driver_module.UP){
		elevator.moving = true
		elevator.direction = driver_module.UP
		event_chan.New_direction_up <- 1
	} else if (direction == driver_module.DOWN){
		elevator.moving = true
		elevator.direction = driver_module.DOWN
		event_chan.New_direction_down <- 1
	}

}

func handle_new_connection(ip string, queue queue_module.Queue_type, current_floor int){


}

func start_up(sensor_channels sensor_module.External_channels, event_channels FSM_module.External_channels, network_channels network_module.Net_channels ){



	sensor_module.Sensors(sensor_channels)
	go FSM_module.Event_generator(event_channels)
    network_module.Network_setup(network_channels)
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

func check_stop_cond(queue * queue_module.Queue_type, elevator * elevator_type, event_channels FSM_module.External_channels, network_channels network_module.Net_channels){
		printc.Data_with_color(printc.COLOR_GREEN, "check_stop_cond elevator.floor: ", elevator.floor)


	var mail network_module.Mail

		var post queue_module.Queue_post

	if(queue.Should_elevator_stop(elevator.floor, driver_module.Convert_dir_to_button(elevator.direction), &post)){

		event_channels.Right_floor <- 1
		elevator.direction = -1
		elevator.moving = false

		mail.Make_mail("", network_module.ORDER_EXECUTED, post.Floor, post.Button_type, 0)
		network_channels.Send_to_all <- mail


	}

}

func make_new_action(queue * queue_module.Queue_type, elevator elevator_type, event_channels FSM_module.External_channels, network_channels network_module.Net_channels){

	var post queue_module.Queue_post
	var mail network_module.Mail

	if(queue.Should_elevator_stop(elevator.floor, driver_module.Convert_dir_to_button(elevator.direction), &post)){

		event_channels.Right_floor <- 1
		mail.Make_mail("", network_module.ORDER_EXECUTED, post.Floor, post.Button_type, 0)
		network_channels.Send_to_all <- mail

	}else if(queue.Get_new_direction(elevator.floor) == driver_module.UP){
		event_channels.New_direction_up <- 1
	}else if(queue.Get_new_direction(elevator.floor) == driver_module.DOWN){
		event_channels.New_direction_down <- 1
	}

	printc.Data_with_color(printc.COLOR_GREEN, "jeg skal: ", queue.Get_new_direction(elevator.floor))

}

func go_to_defined_floor(elevator * elevator_type, sensor_channels sensor_module.External_channels){

	var floor int
	

	driver_module.Elev_start_engine(driver_module.DOWN)

	for floor = -1 ; floor <= -1;{

		floor = <- sensor_channels.Floor_chan

	}
	driver_module.Elev_stop_engine()
	elevator.floor = floor
	elevator.direction = -1
	elevator.moving = false

}