package FSM_module


func State_machine(state * int, event int){

	switch(state){

	case idle:

		idle_state(&state,event)
	case door_open:

		door_open_state(&state,event)
	case moving:

		moving_state(&state,event)

	}
}


func idle_state(state * int,event int){

	switch(event){
		
	case NEW_DIRECTION_E:

		state = moving

	case STOP_E:

		state = door_open
		externalChan.open_door <- 1

	}

}

func door_open_state(state * int, event int){

	switch(event){

	case CLOSE_DOOR_E:

		externalChan.close_door <- 1
		state = idle
	}
}

func moving_state(state * int, event int){

	switch(event){
		
	case STOP_E:

		externalChan.open_door <- 1
		state = door_open
	}
}

