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


func Event_generator(external_chan ExtarnalChannels){

	var internal_chan internal_channels
	internal_chan.internal_channels_FSM_init()

	var state = IDLE

	select{

	case <- external_chan.reached_floor:

		event_reached_floor(&state, internal_chan)

	case <- external_chan.stop:

		event_stop(&state, internal_chan)

	case <- external_chan.ascend:

		event_ascend(&state, internal_chan)

	case <- external_chan.descend:

		event_descend(&state, internal_chan)

	case <- internal_chan.open_door:

		event_open_door(&state, internal_chan)

	case <- internal_chan.close_door:

		event_close_door(&state, internal_chan)

	case direction := <- external_chan.new_direction:

		event_new_direction(&state, direction, internal_chan)
	}
}

func event_reached_floor(state * int, internal_chan internal_channels){

	state_machine(*state, REACHED_FLOOR_E, internal_chan)
}

func event_stop(state * int, internal_chan internal_channels){

	driver_module.Elev_stop_engine()
	state_machine(*state, STOP_E, internal_chan)

}

func event_ascend(state * int, internal_chan internal_channels){

	driver_module.Elev_start_engine(driver_module.UP)
	state_machine(*state, ASCEND_E, internal_chan)
}

func event_descend(state * int, internal_chan internal_channels){

	driver_module.Elev_start_engine(driver_module.DOWN)
	state_machine(*state, DESCEND_E, internal_chan)
}

func event_open_door(state * int, internal_chan internal_channels){

	driver_module.Elev_set_door_light(1)
	state_machine(*state, OPEN_DOOR_E, internal_chan)
}

func event_close_door(state * int, internal_chan internal_channels){

	driver_module.Elev_set_door_light(0)
	state_machine(*state, CLOSE_DOOR_E, internal_chan)
}

func event_new_direction(state * int, direction int, internal_chan internal_channels){

	if(direction == driver_module.UP){
		external_chan.ascend <- 1
	} else if(direction == driver_module.DOWN){
		external_chan.descend <- 1
	}
	
	state_machine(*state, NEW_DIRECTION_E, internal_chan)
}