package queue_module

import (
	"driver_module"
	"encoding/json"
	"math"
	. "debug_module"
)

var queue [QUEUE_SIZE]int
var order_lights [4][3]bool

const QUEUE_SIZE = 12

func Queue_insert(insert_floor int, insert_type driver_module.Elev_button_type_t, current_floor int){

	Debug_message("got queue insert " + string(insert_floor) + " " + string(current_floor), "Queue_insert")

	prev := current_floor
	var direction driver_module.Elev_button_type_t

	turn_on_light(insert_floor, insert_type)

	for i := 0; i < QUEUE_SIZE; i++ {

		if(queue[i] == insert_floor){
			break

		}else if(queue[i] == -1){
			queue[i] = insert_floor
			break

		}else if(prev < queue[i]){
			direction = driver_module.BUTTON_CALL_UP

			if(insert_floor < queue[i] && insert_floor > prev){
				if(insert_type == direction || insert_type == driver_module.BUTTON_COMMAND){
				queue_insert_to_pos(insert_floor, insert_type, i)
				break
				}	
			}


		} else{
			direction = driver_module.BUTTON_CALL_DOWN
		
			if(insert_floor > queue[i] && insert_floor < prev){
				if(insert_type == direction || insert_type == driver_module.BUTTON_COMMAND){
				queue_insert_to_pos(insert_floor, insert_type, i)
				break
				}	
			}
		}

		prev = queue[i]

	}
	Debug_message("Ferdig!", "Queue_insert")
}

func Get_insertion_cost(insert_floor int, insert_type int, current_floor int)(int){

	prev := current_floor
	direction := 0
	cost:=0

	for i := 0; i < QUEUE_SIZE; i++ {

		if(queue[i] == insert_floor){
			break

		}else if(queue[i] == -1){
			break

		}else if(prev < queue[i]){
			direction = driver_module.BUTTON_CALL_UP

			if(insert_floor < queue[i] && insert_floor > prev){
				if(insert_type == direction || insert_type == driver_module.BUTTON_COMMAND){
				break
				}	
			}

		} else{
			direction = driver_module.BUTTON_CALL_DOWN
		
			if(insert_floor > queue[i] && insert_floor < prev){
				if(insert_type == direction || insert_type == driver_module.BUTTON_COMMAND){
				break
				}	
			}
		}
		
		cost += int(math.Abs(float64(prev - queue[i])))
		prev = queue[i]

	}

	cost += int(math.Abs(float64(prev - insert_floor)))

	return cost

}

func Init_queue()(){

	var j driver_module.Elev_button_type_t

	for i := 0; i < QUEUE_SIZE; i++ {
			queue[i] = -1
	}

	for i :=0; i<driver_module.N_FLOORS; i++{
		for j = 0; j<driver_module.N_BUTTONS; j++{
			order_lights[i][j] = false
			driver_module.Elev_set_button_lamp(j,i,0)
		}
	}

}

func One_direction(current_floor int) int{
	if(queue[0]==-1 ){
		return -1
	}else if(current_floor < queue[0]){
		return driver_module.UP
	} else{
		return driver_module.DOWN
	}

}

func Should_elevator_stop(current_floor int) bool{
	if(current_floor == queue[0]){
		queue_remove_multiple_floors(current_floor)
		turn_off_lights(current_floor)
		return true
	}
	return false
}

func Pop_queue() (result int){

	result = queue[0]

	for i := 1; i < QUEUE_SIZE; i++ {

		queue[i-1] = queue[i]
		if(queue[i] == -1){break}
	}

	queue[QUEUE_SIZE-1] = -1

	queue_remove_multiple_floors(result)

	return result
}

func Get_queue_json(queue [QUEUE_SIZE]int)(queue_encoded []byte){

	queue_encoded, _ = json.Marshal(queue)
	return queue_encoded
}


func queue_insert_to_pos(insert_floor int, insert_type driver_module.Elev_button_type_t, position int){

	var swap int
	swap = insert_floor
	var swap_tmp int
	

	for i := position; i < QUEUE_SIZE; i++ {
		
		if(queue[i] == -1){
			queue[i] = swap
			break
		}

		swap_tmp = queue[i]
		queue[i] = swap
		swap = swap_tmp

	}

}

func queue_remove_multiple_floors(floor int){

	previndex :=0

	for i := 0; i < QUEUE_SIZE; i++ {

		queue[previndex] = queue[i]

		if(queue[i]==floor){continue}

		previndex++
	}

	for i := previndex; i< QUEUE_SIZE; i++ {queue[i] = -1}

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