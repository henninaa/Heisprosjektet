package bank_module

import(
	"time"
	"queue_module"
	"netwok_module"
	"driver_module"
	)


func Elevator_main_control(){

	elevator elevator_type
	elevator.elevator_type_init()

	queue queue_module
	queue.queue_type_init()
	
	mail_buffer []netwok_module.Mail

	for{

		select{

			case msg := <- networkChan.inbox:

				handle_network_messgage(msg)

			case post := externalChan.new_order:

				handle_new_order(post)

			case post := <- internChan.insert_to_queue:

				insert_post_to_queue(post, elevator.floor)
				
			case post := <- internChan.auction_order:

				deliver_order(post, queue.Get_lowest_cost_ip(post))

			case floor := <- externalChan.new_floor:

				handle_new_floor(floor, &(elevator.floor))

			case <- externalChan.request_new_direction:

				handle_new_direction(queue.request_new_direction())

			case direction := internalChan.new_direction:

				handle_new_direction(direction)

			case ip := <- externalChan.new_connection:

				handle_new_connection(queue.Get_queue_json, ip)

		}

	}

}

func handle_network_messgage(mail netwok_module.Mail){

	switch (mail.Msg.msg_type){

	case netwok_module.ORDER_TAKEN:

	case netwok_module.ORDER_EXECUTED:

	case netwok_module.DELIVER_ORDER:

	case netwok_module.TAKE_NEW_ORDER:

		insert_post_to_queue(queue_module.convert_mail_to_post(mail))

	case netwok_module.TAKE_BACKUP_ORDER:

	case netwok_module.BACKUP_ORDER_COMPLETE:

		internChan.take_backup_order <- queue_module.convert_mail_to_backup_post(mail)

	case netwok_module.ERROR_MSG:

		idk.com

	case netwok_module.TAKE_BACKUP_FLOOR:

	}
}

func handle_new_order(post queue_module.queue_post){

	if (post.button_type == driver_module.BUTTON_COMMAND){
		internChan.insert_to_queue <- post
	} else{

		internChan.auction_order <- post

	}

}

func deliver_order(post queue_module.backup_post, IP string){

	var mail netwok_module.Mail

	if(IP == "self"){
		internChan.insert_to_queue <- post
	}else{

		mail.make_mail(IP, netwok_module.DELIVER_ORDER, JSON.Marshal(post))
		netwok_module.externalChan.send_to_one <- mail
	}

}

func insert_post_to_queue(post queue_module.queue_post, floor int){

	var mail netwok_module.Mail

	mail.Make_mail("", netwok_module.ORDER_TAKEN, post.floor, driver_module.button_type_to_int(post.button_type))
	queue.insert_queue(post.queue, post.button_type, elevator.floor)
}

func handle_new_floor(floor_input int, current_floor * int){

	current_floor = floor_input
}

func handle_new_direction(){


}

func handle_new_connection(queue_backup []byte, ip string){


}