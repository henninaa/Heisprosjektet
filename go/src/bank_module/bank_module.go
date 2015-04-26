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

	go_to_defined_floor(&elevator, sensor_channels, event_channels)

	var queue queue_module.Queue_type
	queue.Init(elevator.floor)

	for{
		time.Sleep(ELEVATOR_MAIN_CONTROL_INTERVAL)
		select{

			case msg := <- network_channels.Inbox:
				printc.Data_with_color(printc.COLOR_GREEN, "Getting a message from network_channels.Inbox")
				handle_network_messgage(msg, internal_chan)

			
			case post := <- sensor_channels.Order_chan:
				printc.Data_with_color(printc.COLOR_GREEN, "Getting a message from sensor_channels.Order_chan")
				handle_new_order(post,internal_chan.insert_to_queue, internal_chan.auction_order)

			
			case order := <- internal_chan.insert_to_queue:
				printc.Data_with_color(printc.COLOR_GREEN, "Getting a message from internal_chan.insert_to_queue")
				insert_order_to_local_queue(order, &queue, elevator.floor, network_channels.Send_to_all, event_channels)

			
			case post := <- internal_chan.auction_order:
				printc.Data_with_color(printc.COLOR_GREEN, "Getting a message from internal_chan.auction_order")
				auction_order(post, &queue, elevator.floor, internal_chan.insert_to_queue, network_channels.Send_to_one)

			
			case mail := <- internal_chan.remote_order_executed:
				handle_remote_order_executed(mail, &queue)
				printc.Data_with_color(printc.COLOR_GREEN, "Getting a message from internal_chan.remote_order_executed")

			
			case floor := <- sensor_channels.Floor_chan:
				printc.Data_with_color(printc.COLOR_GREEN, "Getting a message from sensor_channels.Floor_chan")
				handle_new_floor(floor, &(elevator.floor), internal_chan.check_stop_conditions, network_channels.Send_to_all)

			
			case direction := <- internal_chan.new_direction:
				printc.Data_with_color(printc.COLOR_GREEN, "Getting a message from internal_chan.new_direction")
				handle_new_direction(direction, &elevator, event_channels)

			
			case ip := <- network_channels.New_connection:
				printc.Data_with_color(printc.COLOR_GREEN, "Getting a message from network_channels.New_connection: ", ip)
				handle_new_connection(ip, queue, elevator.floor, network_channels)
			
			case <- internal_chan.check_stop_conditions:
				printc.Data_with_color(printc.COLOR_GREEN, "Getting a message from internal_chan.check_stop")
				check_stop_cond(&queue, &elevator, event_channels, network_channels)
			
			case <- event_channels.Get_new_action:
				printc.Data_with_color(printc.COLOR_GREEN, "Getting a message from event_channels.Get_new_action")
				make_new_action(&queue, &elevator, event_channels, network_channels)

			
			case dead_elevator := <- network_channels.Get_dead_elevator:
				printc.Data_with_color(printc.COLOR_RED, "Getting a message from network_channels.Get_dead_elevator: ", dead_elevator)
				handle_dead_elevator(dead_elevator, &queue, elevator.floor, event_channels)
			
			case mail :=<- internal_chan.take_backup_order:

				handle_take_backup_order(mail, &queue)

			case mail := <- internal_chan.take_backup_floor:

				handle_take_backup_floor(mail, &queue)

			case <- event_channels.Engine_error:

				handle_engine_error(network_channels, sensor_channels, event_channels, internal_chan.abort_light_show, &queue, &elevator)

			case dead_elevator := <- internal_chan.lost_engine_on_network:

				handle_dead_engine_on_network(dead_elevator, &queue, elevator.floor, event_channels)

			case IP := <- internal_chan.engine_recovery_on_network:

				handle_engine_recovery_on_network(&queue, IP)

			case dir := <- event_channels.New_dir:

				elevator.direction = dir
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

		internal_chan.remote_order_executed <- mail
		printc.Data_with_color(printc.COLOR_GREEN, "ORDER_EXECUTED.MAIL ", mail)

	case network_module.TAKE_NEW_ORDER:

		internal_chan.insert_to_queue <- queue_module.Convert_mail_to_queue_post(mail)
		printc.Data_with_color(printc.COLOR_GREEN, "TAKE_NEW_ORDER.MAIL ", mail)
		printc.Data_with_color(printc.COLOR_BLACK, "				mail mottat: ", mail)

	case network_module.TAKE_BACKUP_ORDER:

		internal_chan.take_backup_order <- mail
		printc.Data_with_color(printc.COLOR_GREEN, "TAKE_NEW_BACKUP_ORDER.MAIL ", mail)

	case network_module.TAKE_BACKUP_FLOOR:

		internal_chan.take_backup_floor <- mail
		printc.Data_with_color(printc.COLOR_GREEN, "TAKE_NEW_BACKUP_FLOOR.MAIL ", mail)

	case network_module.ENGINE_FAILURE:

		internal_chan.lost_engine_on_network <- mail.IP

	case network_module.ENGINE_RECOVERY:

		internal_chan.engine_recovery_on_network <- mail.IP


	}
}

func handle_new_order(order queue_module.Queue_post, insert_to_queue chan queue_module.Queue_post, auction_order chan queue_module.Queue_post){

	if (order.Button_type == driver_module.BUTTON_COMMAND){
		insert_to_queue <- order
	} else{

		auction_order <- order

	}

}

func auction_order(post queue_module.Queue_post, queue * queue_module.Queue_type, current_floor int, insert_to_queue chan queue_module.Queue_post, send_to_one chan network_module.Mail){

	var mail network_module.Mail

	IP := queue.Get_lowest_cost_ip(post, current_floor)

	if(IP == "self"){
		insert_to_queue <- post
	}else{

		mail.Make_mail(IP, network_module.TAKE_NEW_ORDER, post.Floor, post.Button_type, 0)
		send_to_one <- mail
		printc.Data_with_color(printc.COLOR_BLACK, "				post.Floor: ", post.Floor)
		printc.Data_with_color(printc.COLOR_BLACK, "				mail sendt: ", mail)
	}

}

func insert_order_to_local_queue(post queue_module.Queue_post, queue * queue_module.Queue_type, current_floor int,send_to_all chan network_module.Mail, event_channels FSM_module.External_channels){

	var mail network_module.Mail

	mail.Make_mail("", network_module.TAKE_BACKUP_ORDER, post.Floor, post.Button_type, 0)
	send_to_all <- mail

	queue.Insert_to_own_queue(post, current_floor)

	event_channels.New_order <- 1

}

func handle_new_floor(floor_input int, current_floor * int, stop_check_chan chan int, send_to_all chan network_module.Mail){

	if(floor_input!=-1){

		*current_floor = floor_input
		stop_check_chan <- 1

		var mail network_module.Mail
		mail.Make_mail("", network_module.TAKE_BACKUP_FLOOR, floor_input, driver_module.BUTTON_COMMAND, 0)
		send_to_all <- mail
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

func handle_new_connection(ip string, queue queue_module.Queue_type, current_floor int, network_channels network_module.Net_channels){

	var mail network_module.Mail
	mail.Make_mail(ip, network_module.TAKE_BACKUP_FLOOR, current_floor, 0,0)

	network_channels.Send_to_one <- mail

}

func start_up(sensor_channels sensor_module.External_channels, event_channels FSM_module.External_channels, network_channels network_module.Net_channels ){

	sensor_module.Sensors(sensor_channels)
	go FSM_module.Event_generator(event_channels)
    go network_module.Network_setup(network_channels)
}

func convert_mail_to_backup_post(mail network_module.Mail, ip string)(backup_post queue_backup_post){

	backup_post.post = queue_module.Convert_mail_to_queue_post(mail)
	backup_post.IP = ip

	return backup_post
}

func handle_remote_order_executed(mail network_module.Mail, queue * queue_module.Queue_type){

	post := queue_module.Convert_mail_to_queue_post(mail)

	queue.Remove_post_from_backup_queue(post, mail.IP)

}

func handle_take_backup_floor(mail network_module.Mail, queue * queue_module.Queue_type){

	queue.New_backup_floor(mail.Msg.Floor, mail.IP)
}

func check_stop_cond(queue * queue_module.Queue_type, elevator * elevator_type, event_channels FSM_module.External_channels, network_channels network_module.Net_channels){
		printc.Data_with_color(printc.COLOR_GREEN, "check_stop_cond elevator.floor: ", elevator.floor)
		printc.Data_with_color(printc.COLOR_RED, "Tuleluuuuuu")


	var mail network_module.Mail

		var post queue_module.Queue_post

	if(queue.Should_elevator_stop(elevator.floor, driver_module.Convert_dir_to_button(elevator.direction), &post)){

		printc.Data_with_color(printc.COLOR_CYAN, "Tuleluuuuuu")
		event_channels.Right_floor <- 1
		elevator.direction = -1
		elevator.moving = false

		mail.Make_mail("", network_module.ORDER_EXECUTED, post.Floor, post.Button_type, 0)
		network_channels.Send_to_all <- mail


	}else{
		printc.Data_with_color(printc.COLOR_RED, "HIIIIIIIIT------------------------------------------------------")
		if(queue.Going_too_far_up(elevator.floor, elevator.direction)){
			event_channels.Too_far_up <- 1
			printc.Data_with_color(printc.COLOR_RED, "DOWNS")
		}else if(queue.Going_too_far_down(elevator.floor, elevator.direction)){
			event_channels.Too_far_down <- 1
			printc.Data_with_color(printc.COLOR_RED, "UP")
		}
	}

}

func make_new_action(queue * queue_module.Queue_type, elevator * elevator_type, event_channels FSM_module.External_channels, network_channels network_module.Net_channels){

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

func go_to_defined_floor(elevator * elevator_type, sensor_channels sensor_module.External_channels, event_channels FSM_module.External_channels){

	var floor int
	

	driver_module.Elev_start_engine(driver_module.DOWN)

	for floor = -1 ; floor <= -1;{

		floor = <- sensor_channels.Floor_chan

	}
	driver_module.Elev_stop_engine()
	elevator.floor = floor
	elevator.direction = -1
	elevator.moving = false

	event_channels.New_order <- 1

}


func handle_dead_elevator(dead_elevator string, queue * queue_module.Queue_type, current_floor int, event_channels FSM_module.External_channels){

	queue.Annex_backup(dead_elevator, current_floor)

	event_channels.New_order <- 1

}

func handle_take_backup_order(mail network_module.Mail, queue * queue_module.Queue_type){

	queue.Insert_to_backup_queue(mail.Msg.Floor, mail.Msg.Dir, mail.IP)
}

func handle_engine_error(network_channels network_module.Net_channels, sensor_channels sensor_module.External_channels, event_channels FSM_module.External_channels, abort_light_show chan bool, queue * queue_module.Queue_type, elevator * elevator_type){

	var mail network_module.Mail
	mail.Make_mail("", network_module.ENGINE_FAILURE, 0, driver_module.BUTTON_COMMAND,0)

	queue.Delete_local_orders()

	go light_show(abort_light_show)

	network_channels.Send_to_all <- mail

	sensor_channels.Deactivate_orders <- true

	floor := <- sensor_channels.Floor_chan

	if(floor != -1){elevator.floor = floor}

	abort_light_show <- true
	event_channels.Right_floor <- 1
	sensor_channels.Activate_orders <- true

	mail.Make_mail("", network_module.ENGINE_RECOVERY, 0, driver_module.BUTTON_COMMAND,0)
	network_channels.Send_to_all <- mail
	
	printc.Data_with_color(printc.COLOR_YELLOW, "Tuleluuuuuu")
}

func handle_engine_recovery_on_network(queue * queue_module.Queue_type, IP string){

	queue.Engine_recovered(IP)
}

func light_show(abort chan bool){

	var j driver_module.Elev_button_type_t
	for{

		select{
		case <- abort:
			return
		default:
			for j = 0; j<driver_module.N_BUTTONS; j++{
				for i :=0; i<driver_module.N_FLOORS; i++{
					driver_module.Elev_set_button_lamp(j,i,1)
					time.Sleep(100*time.Millisecond)
					driver_module.Elev_set_button_lamp(j,i,0)
				}
			}
		}
	}
}

func handle_dead_engine_on_network(dead_elevator string, queue * queue_module.Queue_type, floor int, event_channels FSM_module.External_channels){

	handle_dead_elevator(dead_elevator, queue, floor, event_channels)
	queue.Engine_failed(dead_elevator)

}