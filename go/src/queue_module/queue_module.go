package queue_module

import (
	"driver_module"
	"network_module"
	"encoding/json"
	"math"
	. "debug_module"
)


const QUEUE_SIZE = 12

func (queue * Queue_t) Queue_insert(insert_floor int, insert_type driver_module.Elev_button_type_t, current_floor int){

	var input_post queue_post
	input_post.floor = insert_floor
	input_post.button_type = insert_type

	Debug_message("got queue insert " + string(insert_floor) + " " + string(current_floor), "Queue_insert")

	prev := current_floor
	var direction driver_module.Elev_button_type_t

	turn_on_light(insert_floor, insert_type)

	for post := range queue.queue; {

		if(post == input_post{
			break

		}else if(post.floor == -1){
			post = input_post
			break

		}else if(prev < post.floor){
			direction = driver_module.BUTTON_CALL_UP

			if(insert_floor < post.floor && insert_floor > prev){
				if(insert_type == direction || insert_type == driver_module.BUTTON_COMMAND){
				queue.queue_insert_to_pos(input_post, i)
				break
				}	
			}


		} else{
			direction = driver_module.BUTTON_CALL_DOWN
		
			if(insert_floor > post.floor && insert_floor < prev){
				if(insert_type == direction || insert_type == driver_module.BUTTON_COMMAND){
				queue.queue_insert_to_pos(input_post, i)
				break
				}	
			}
		}

		prev = post.floor

	}
	Debug_message("Ferdig!", "Queue_insert")
}

func Get_insertion_cost(insert_floor int, insert_type int, current_floor int, queue [QUEUE_SIZE]queue_post)(int){

	prev := current_floor
	direction := 0
	cost:=0

	for i := 0; i < QUEUE_SIZE; i++ {

		if(queue[i].floor == insert_floor){
			break

		}else if(queue[i].floor == -1){
			break

		}else if(prev < queue[i].floor){
			direction = driver_module.BUTTON_CALL_UP

			if(insert_floor < queue[i].floor && insert_floor > prev){
				if(insert_type == direction || insert_type == driver_module.BUTTON_COMMAND){
				break
				}	
			}

		} else{
			direction = driver_module.BUTTON_CALL_DOWN
		
			if(insert_floor > queue[i].floor && insert_floor < prev){
				if(insert_type == direction || insert_type == driver_module.BUTTON_COMMAND){
				break
				}	
			}
		}
		
		cost += int(math.Abs(float64(prev - queue[i].floor)))
		prev = queue[i].floor

	}

	cost += int(math.Abs(float64(prev - insert_floor)))

	return cost

}

func Init_queue()(queue [12]){

	for i := 0; i < QUEUE_SIZE; i++ {
			queue[i] = -1
	}

	return queue

}

func (queue * Queue_type) Get_new_direction(current_floor int) int{
	if(queue.queue[0]==-1 ){
		return -1
	}else if(current_floor < queue.queue[0]){
		return driver_module.UP
	} else{
		return driver_module.DOWN
	}

}

func (queue * queue_type) Should_elevator_stop(current_floor int, driver_module.Elev_button_type_t) bool{

	if(current_floor == queue.queue[0].floor; queue.queue[0].button_type == driver_module.BUTTON_COMMAND || queue.queue[0].button_type == direction){
		queue.queue_remove_multiple_floors(queue.queue[0])
		//turn_off_lights(current_floor)
		return true
	}
	return false
}

func (queue * queue_type) Get_lowest_cost_ip(insert_post queue_post)(lowest_cost_ip string){

	lowest_cost := 999999
	cost := -1

	for post := range queue.backup{

		cost = Get_insertion_cost(insert_post.floor, insert_post.button_type, current_floor, post.queue)

		if(cost < lowest_cost){
			lowest_cost_ip = post.IP
			lowest_cost = cost
		}
	}

	cost = Get_insertion_cost(insert_post.floor, insert_post.button_type, current_floor, queue.queue)

	if(cost < lowest_cost){
		lowest_cost_ip = "self"
	}

	return lowest_cost_ip
}


func (queue * queue_type) queue_insert_to_pos(insert_post queue_post, position int){

	var swap queue_post
	swap = insert_post
	var swap_tmp queue_post
	

	for i := position; i < QUEUE_SIZE; i++ {
		
		if(queue.queue[i].floor == -1){
			queue.queue[i] = swap
			break
		}

		swap_tmp = queue.queue[i]
		queue.queue[i] = swap
		swap = swap_tmp

	}

}

func (queue * queue_type) queue_remove_multiple_floors(post queue){

	previndex :=0

	for i := 0; i < QUEUE_SIZE; i++ {

		queue[previndex] = queue[i]

		if(queue[i]==post){continue}

		previndex++
	}

	for i := previndex; i< QUEUE_SIZE; i++ {queue[i].floor = -1}

}

func turn_on_light(insert_floor int, insert_type driver_module.Elev_button_type_t){

	
	if (order_lights[insert_floor][insert_type] == false){
		driver_module.Elev_set_button_lamp(insert_type, insert_floor, 1)
		order_lights[insert_floor][insert_type] = true
	}
}

func turn_off_lights(floor int){

	var order_type driver_module.Elev_button_type_t

	for order_type = 0; order_type<3; order_type++{

		if(order_lights[floor][order_type]){
			order_lights[floor][order_type] = false
			driver_module.Elev_set_button_lamp(order_type, floor, 0)
		}
	}

}

func convert_mail_to_queue_post()(){

	var mail network_module.Mail
}


func convert_mail_to_backup_post()(){


}