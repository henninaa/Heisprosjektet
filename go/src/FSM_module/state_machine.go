package FSM_module


func state_machine(state * int, event int, internal_chan internal_channels){

	switch(*state){

	case initialize:

		init_state(state,event, internal_chan)
	case idle:

		idle_state(state,event, internal_chan)
	case door_open:

		door_open_state(state,event, internal_chan)
	case moving:

		moving_state(state,event, internal_chan)

	}
}


func idle_state(state * int,event int, internal_chan internal_channels){

	switch(event){
		
	case NEW_DIRECTION_E:

		*state = moving
		internal_chan.start_moving <- 1

	case STOP_E:

		*state = door_open
		internal_chan.open_door <- 1

	}

}

func door_open_state(state * int, event int, internal_chan internal_channels){

	switch(event){

	case CLOSE_DOOR_E:

		*state = idle
	}
}

func moving_state(state * int, event int, internal_chan internal_channels){

	switch(event){
		
	case STOP_E:


		*state = door_open
	}
}

func init_state(state * int ,event int, internal_chan internal_channels){

	switch(event){
		
	case REACHED_FLOOR_E:

		internal_chan.Stop <- 1
		*state = idle
	}
}