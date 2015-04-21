package FSM_module

import(
	"driver_module"
	)

const(
	STOP_E 			= 0x00
	ASCEND_E 		= 0x01
	DESCEND_E		= 0x02
	OPEN_DOOR_E		= 0x03
	CLOSE_DOOR_E	= 0x04
	NEW_DIRECTION_E = 0x05
	REACHED_FLOOR_E	= 0x06
	)


func Event_generator(){

	var state = IDLE

	select{

	case <- externalChan.reached_floor:

		event_reached_floor(&state)

	case <- externalChan.stop:

		event_stop(&state)

	case <- externalChan.ascend:

		event_ascend(&state)

	case <- externalChan.descend:

		event_descend(&state)

	case <- externalChan.open_door:

		event_open_door(&state)

	case <- external_chan.close_door:

		event_close_door(&state)

	case direction := <- externalChan.new_direction:

		event_new_direction(&state, direction)
	}
}

func event_reached_floor(state * int){

	state_machine(*state, REACHED_FLOOR_E)
}

func event_stop(state * int){

	driver_module.Elev_stop_engine()
	state_machine(*state, STOP_E)

}

func event_ascend(state * int){

	driver_module.Elev_start_engine(driver_module.UP)
	state_machine(*state, ASCEND_E)
}

func event_descend(state * int){

	driver_module.Elev_start_engine(driver_module.DOWN)
	state_machine(*state, DESCEND_E)
}

func event_open_door(state * int){

	driver_module.Elev_set_door_light(1)
	state_machine(*state, OPEN_DOOR_E)
}

func event_close_door(state * int){

	driver_module.Elev_set_door_light(0)
	state_machine(*state, CLOSE_DOOR_E)
}

func event_new_direction(state * int, direction int){

	if(direction == driver_module.UP){
		externalChan.ascend <- 1
	} else if(direction == driver_module.DOWN){
		externalChan.descend <- 1
	}
	
	state_machine(*state, NEW_DIRECTION_E)
}