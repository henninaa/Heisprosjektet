package FSM_module

import("printc")


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

	printc.Data_with_color(printc.COLOR_CYAN, "state: ", event)

	switch(event){

	case NEW_ORDER_E:

		external_chan.Get_new_action <- 1
		
	case NEW_DIRECTION_DOWN_E:

		*state = moving
		internal_chan.descend <- 1
		go breakdown_timer(internal_chan)

	case NEW_DIRECTION_UP_E:

		*state = moving
		internal_chan.ascend <- 1
		go breakdown_timer(internal_chan)

	case RIGHT_FLOOR_E:

		*state = door_open
		internal_chan.open_door <- 1


	}

}

func door_open_state(state * int, event int, internal_chan internal_channels){

	switch(event){

	case CLOSE_DOOR_E:

		*state = idle

		external_chan.Get_new_action  <- 1
	}
}

func moving_state(state * int, event int, internal_chan internal_channels){

	switch(event){
		
	case RIGHT_FLOOR_E:

		internal_chan.stop <- 1
		*state = door_open
	
	case TOO_FAR_DOWN_E:
		internal_chan.ascend <- 1

	case TOO_FAR_UP_E:
		internal_chan.descend <- 1

	}
}

func init_state(state * int ,event int, internal_chan internal_channels){

	switch(event){
		
	case RIGHT_FLOOR_E:

		internal_chan.stop <- 1
		*state = idle
	}
}