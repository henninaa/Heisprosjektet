package FSM_module


func State_machine(state * int, event int, internal_chan interal_channels){

	switch(state){

	case idle:

		idle_state(&state,event, internal_chan)
	case door_open:

		door_open_state(&state,event, internal_chan)
	case moving:

		moving_state(&state,event, internal_chan)

	}
}


func idle_state(state * int,event int, internal_chan interal_channels){

	switch(event){
		
	case NEW_DIRECTION_E:

		state = moving

	case STOP_E:

		state = door_open
		internal_chan.open_door <- 1

	}

}

func door_open_state(state * int, event int, internal_chan interal_channels){

	switch(event){

	case CLOSE_DOOR_E:

		internal_chan.close_door <- 1
		state = idle
	}
}

func moving_state(state * int, event int, internal_chan interal_channels){

	switch(event){
		
	case STOP_E:

		internal_chan.open_door <- 1
		state = door_open
	}
}

