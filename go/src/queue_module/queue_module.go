package queue_module

import (
	. "driver_elev"
	"encoding/json"
	"math"
	//"debug_elevator"
)

var Queue [QUEUE_SIZE]int



const(
	UP = 	0x01
	DOWN = 	0x02
	STILL = 0x03
	QUEUE_SIZE = 12
	)



func Queue_insert(insert_floor int, insert_type Elev_button_type_t, current_floor int, queue * [QUEUE_SIZE]int){

	prev := current_floor
	var direction Elev_button_type_t

	if(prev == insert_floor){return}

	for i := 0; i < QUEUE_SIZE; i++ {

		if(queue[i] == -1){
			queue[i] = insert_floor
		}else if(queue[i] == insert_floor){
			break;
		}

		if(prev < queue[i]){
			direction = BUTTON_CALL_UP
		} else{
			direction = BUTTON_CALL_DOWN
		}

		if(insert_type == direction || insert_type == BUTTON_COMMAND){
			queue_insert_to_pos(insert_floor, insert_type, i, queue)
			break
		}

	}
}

func Get_insertion_cost(insert_floor int, insert_type int, current_floor int, queue * [QUEUE_SIZE]int)(int){

	prev := current_floor
	direction := 0
	cost:=0

	for i := 0; i < QUEUE_SIZE; i++ {

		if(queue[i] == -1){
			break;
		}

		if(prev < queue[i]){
			direction = BUTTON_CALL_UP
		} else{
			direction = BUTTON_CALL_DOWN
		}
		

		if(insert_type == direction || insert_type == BUTTON_COMMAND){
			break;
		} else{
			cost += int(math.Abs(float64(prev - queue[i])))
		}

		prev = queue[i]
	
	}

	cost += int(math.Abs(float64(prev - insert_floor)))

	return cost

}

func Init_queue()(queue [QUEUE_SIZE]int){

	for i := 0; i < QUEUE_SIZE; i++ {
			queue[i] = -1
	}

	return queue
}

func One_direction(current_floor int, queue * [QUEUE_SIZE]int) int{
	if(current_floor < queue[0]){
		return UP
	} else{
		return DOWN
	}

}

func Should_elevator_stop(current_floor int, queue * [QUEUE_SIZE]int) bool{
	if(current_floor == queue[0]){
		queue_remove_multiple_floors(current_floor, queue)
		return true
	}
	return false
}

func Pop_queue(queue * [QUEUE_SIZE]int) (result int){

	result = queue[0]

	for i := 1; i < QUEUE_SIZE; i++ {

		queue[i-1] = queue[i]
		if(queue[i] == -1){break}
	}

	queue[QUEUE_SIZE-1] = -1

	queue_remove_multiple_floors(result, queue)

	return result
}

func Get_queue_json(queue [QUEUE_SIZE]int)(queue_encoded []byte){

	queue_encoded, _ = json.Marshal(queue)
	return queue_encoded
}


func queue_insert_to_pos(insert_floor int, insert_type Elev_button_type_t, position int, queue * [QUEUE_SIZE]int){

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

func queue_remove_multiple_floors(floor int, queue * [QUEUE_SIZE]int){


	previndex :=0

	for i := 0; i < QUEUE_SIZE; i++ {

		queue[previndex] = queue[i]

		if(queue[i]==floor){continue}

		previndex++
	}

	for i := previndex; i< QUEUE_SIZE; i++ {queue[i] = -1}


}