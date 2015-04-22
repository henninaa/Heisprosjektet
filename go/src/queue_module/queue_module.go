package queue_module

import (
	"driver_module"
	"math"
	. "debug_module"
)


const QUEUE_SIZE = 12

func (queue_class * Queue_type) Insert_to_own_queue(insert_floor int, insert_type driver_module.Elev_button_type_t, current_floor int){

	queue_class.queue.insert(insert_floor, insert_type, current_floor)

}

func (queue_class * Queue_type) Insert_to_backup_queue(insert_floor int, insert_type driver_module.Elev_button_type_t, current_floor int, ip string){

	for i:= 0; i<len(queue_class.backup); i++{
		if(queue_class.backup[i].IP == ip){
			queue_class.backup[i].queue.insert(insert_floor, insert_type, current_floor)
			return
		}
	}
	

}

func (queue_class * Queue_type) Send_backup_to_auction(IP string, order_chan chan Queue_post){

for i:= 0; i<len(queue_class.backup); i++{
		if(queue_class.backup[i].IP == IP){
			
			for j := range queue_class.backup[i].queue.list{

				if(queue_class.backup[i].queue.list[j].Button_type != driver_module.BUTTON_COMMAND){
					order_chan <- queue_class.backup[i].queue.list[j]
				}
			}
		}
	}

}

func (queue * queue_list) insert(insert_floor int, insert_type driver_module.Elev_button_type_t, current_floor int){

	var input_post Queue_post
	input_post.Floor = insert_floor
	input_post.Button_type = insert_type

	Debug_message("got queue insert " + string(insert_floor) + " " + string(current_floor), "Queue_insert")

	prev := current_floor
	var direction driver_module.Elev_button_type_t

	for i := range queue.list {

		if(queue.list[i] == input_post){
			break

		}else if(queue.list[i].Floor == -1){
			queue.list[i] = input_post
			break

		}else if(prev < queue.list[i].Floor){
			direction = driver_module.BUTTON_CALL_UP

			if(insert_floor < queue.list[i].Floor && insert_floor > prev){
				if(insert_type == direction || insert_type == driver_module.BUTTON_COMMAND){
				queue.insert_to_pos(input_post, i)
				break
				}	
			}


		} else{
			direction = driver_module.BUTTON_CALL_DOWN
		
			if(insert_floor > queue.list[i].Floor && insert_floor < prev){
				if(insert_type == direction || insert_type == driver_module.BUTTON_COMMAND){
				queue.insert_to_pos(input_post, i)
				break
				}	
			}
		}

		prev = queue.list[i].Floor

	}
	Debug_message("Ferdig!", "Queue_insert")
}

func (queue * queue_list) get_insertion_cost(insert_floor int, insert_type driver_module.Elev_button_type_t, current_floor int)(int){


	var input_post Queue_post
	input_post.Floor = insert_floor
	input_post.Button_type = insert_type
	cost:=0

	prev := current_floor
	var direction driver_module.Elev_button_type_t

	for i := range queue.list {

		if(queue.list[i] == input_post){
			break

		}else if(queue.list[i].Floor == -1){
			break

		}else if(prev < queue.list[i].Floor){
			direction = driver_module.BUTTON_CALL_UP

			if(insert_floor < queue.list[i].Floor && insert_floor > prev){
				if(insert_type == direction || insert_type == driver_module.BUTTON_COMMAND){
				break
				}	
			}


		} else{
			direction = driver_module.BUTTON_CALL_DOWN
		
			if(insert_floor > queue.list[i].Floor && insert_floor < prev){
				if(insert_type == direction || insert_type == driver_module.BUTTON_COMMAND){
				break
				}	
			}
		}

		cost += int(math.Abs(float64(prev - queue.list[i].Floor)))
		prev = queue.list[i].Floor

	}

	cost += int(math.Abs(float64(prev - insert_floor)))

	return cost
}


func Init_queue()(queue [12]Queue_post){

	for i := 0; i < QUEUE_SIZE; i++ {
			queue[i].Floor = -1
			queue[i].Button_type = driver_module.BUTTON_COMMAND
	}

	return queue

}

func (queue_class * Queue_type) Get_new_direction(current_floor int) int{
	if(queue_class.queue.list[0].Floor==-1 ){
		return -1
	}else if(current_floor < queue_class.queue.list[0].Floor){
		return driver_module.UP
	} else{
		return driver_module.DOWN
	}

}

func (queue_class * Queue_type) Should_elevator_stop(current_floor int, direction driver_module.Elev_button_type_t) bool{

	if(current_floor == queue_class.queue.list[0].Floor && (queue_class.queue.list[0].Button_type == driver_module.BUTTON_COMMAND || queue_class.queue.list[0].Button_type == direction)){
		queue_class.queue.queue_remove_multiple_floors(queue_class.queue.list[0])
		//turn_off_lights(current_floor)
		return true
	}
	return false
}

func (queue * Queue_type) Remove_post_from_backup_queue(remove_floor int, remove_type driver_module.Elev_button_type_t, ip string){

	var post Queue_post
	post.Floor = remove_floor
	post.Button_type = remove_type

	for i:= 0; i<len(queue.backup); i++{
		if(queue.backup[i].IP == ip){

			queue.backup[i].queue.queue_remove_multiple_floors(post)

		}
	}
}

func (queue * Queue_type) Get_lowest_cost_ip(insert_post Queue_post, current_floor int) (lowest_cost_ip string){
	lowest_cost := 999999
	cost := -1

	for i := range queue.backup{

		cost = queue.backup[i].queue.get_insertion_cost(insert_post.Floor, insert_post.Button_type, queue.backup[i].floor)

		if(cost < lowest_cost){
			lowest_cost_ip = queue.backup[i].IP
			lowest_cost = cost
		}
	}

	cost = queue.queue.get_insertion_cost(insert_post.Floor, insert_post.Button_type, current_floor)

	if(cost <= lowest_cost){
		lowest_cost_ip = "self"
	}

	return lowest_cost_ip
}


func (queue * queue_list) insert_to_pos(insert_post Queue_post, position int){

	var swap Queue_post
	swap = insert_post
	var swap_tmp Queue_post
	

	for i := position; i < QUEUE_SIZE; i++ {
		
		if(queue.list[i].Floor == -1){
			queue.list[i] = swap
			break
		}

		swap_tmp = queue.list[i]
		queue.list[i] = swap
		swap = swap_tmp

	}

}

func (queue * queue_list) queue_remove_multiple_floors(post Queue_post){

	previndex :=0

	for i := 0; i < QUEUE_SIZE; i++ {

		queue.list[previndex] = queue.list[i]

		if(queue.list[i]==post){continue}

		previndex++
	}

	for i := previndex; i< QUEUE_SIZE; i++ {queue.list[i].Floor = -1}

}

func (queue * Queue_type) turn_on_light(insert_floor int, insert_type driver_module.Elev_button_type_t){

	
	if (queue.order_lights[insert_floor][insert_type] == false){
		driver_module.Elev_set_button_lamp(insert_type, insert_floor, 1)
		queue.order_lights[insert_floor][insert_type] = true
	}
}

func (queue * Queue_type) turn_off_lights(floor int){

	var order_type driver_module.Elev_button_type_t

	for order_type = 0; order_type<3; order_type++{

		if(queue.order_lights[floor][order_type]){
			queue.order_lights[floor][order_type] = false
			driver_module.Elev_set_button_lamp(order_type, floor, 0)
		}
	}

}

func convert_mail_to_Queue_post()(){

}


func convert_mail_to_backup_post()(){


}