package FSM_module

import(
	"driver_module"
	)




func Event_generator(external_chan External_channels){

	var internal_chan internal_channels
	internal_chan.internal_channels_FSM_init()

	var state = idle

	select{

	case <- external_chan.Reached_floor:

		event_reached_floor(&state, internal_chan)

	case <- external_chan.Stop:

		event_stop(&state, internal_chan)

	case <- internal_chan.ascend:

		event_ascend(&state, internal_chan)

	case <- internal_chan.descend:

		event_descend(&state, internal_chan)

	case <- internal_chan.open_door:

		event_open_door(&state, internal_chan)

	case <- internal_chan.close_door:

		event_close_door(&state, internal_chan)

	case direction := <- external_chan.New_direction:

		event_new_direction(&state, direction, internal_chan)
	}
}

func event_reached_floor(state * int, internal_chan internal_channels){

	state_machine(state, REACHED_FLOOR_E, internal_chan)
}

func event_stop(state * int, internal_chan internal_channels){

	driver_module.Elev_stop_engine()
	state_machine(state, STOP_E, internal_chan)

}

func event_ascend(state * int, internal_chan internal_channels){

	driver_module.Elev_start_engine(driver_module.UP)
	state_machine(state, ASCEND_E, internal_chan)
}

func event_descend(state * int, internal_chan internal_channels){

	driver_module.Elev_start_engine(driver_module.DOWN)
	state_machine(state, DESCEND_E, internal_chan)
}

func event_open_door(state * int, internal_chan internal_channels){

	driver_module.Elev_set_door_open_lamp(true)
	state_machine(state, OPEN_DOOR_E, internal_chan)
}

func event_close_door(state * int, internal_chan internal_channels){

	driver_module.Elev_set_door_open_lamp(false)
	state_machine(state, CLOSE_DOOR_E, internal_chan)
}

func event_new_direction(state * int, direction int, internal_chan internal_channels){

	if(direction == driver_module.UP){
		internal_chan.ascend <- 1
	} else if(direction == driver_module.DOWN){
		internal_chan.descend <- 1
	}
	
	state_machine(state, NEW_DIRECTION_E, internal_chan)
}