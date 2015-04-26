package FSM_module

import(
	"driver_module"
	"time"
	"printc"
	)

var external_chan External_channels


func Event_generator(external_chan_in External_channels){

	external_chan = external_chan_in

	var internal_chan internal_channels
	internal_chan.internal_channels_FSM_init()

	var state = idle
	printc.Data_with_color(printc.COLOR_CYAN, "Event_generator running...")
	for{

		time.Sleep(EVENT_INTERVAL)

		select{

		case <- external_chan.Right_floor:
			printc.Data_with_color(printc.COLOR_RED, "EVENTCHANNEL --> external_chan.Right_floor")
			event_right_floor(&state, internal_chan)

		case <- internal_chan.stop:
			printc.Data_with_color(printc.COLOR_RED, "EVENTCHANNEL --> internal_chan.stop")
			event_stop(&state, internal_chan)

		case <- internal_chan.ascend:
			printc.Data_with_color(printc.COLOR_RED, "EVENTCHANNEL --> internal_chan.ascend")
			event_ascend(&state, internal_chan)

		case <- internal_chan.descend:
			printc.Data_with_color(printc.COLOR_RED, "EVENTCHANNEL --> internal_chan.descend")
			event_descend(&state, internal_chan)

		case <- internal_chan.open_door:
			printc.Data_with_color(printc.COLOR_RED, "EVENTCHANNEL --> internal_chan.open_door")
			event_open_door(&state, internal_chan)

		case <- internal_chan.close_door:
			printc.Data_with_color(printc.COLOR_RED, "EVENTCHANNEL --> internal_chan.close_door")
			event_close_door(&state, internal_chan, external_chan)

		case <- external_chan.New_direction_up:
			printc.Data_with_color(printc.COLOR_RED, "EVENTCHANNEL --> external_chan.New_dir_up")
			event_new_direction_up(&state, internal_chan)

		case <- external_chan.New_direction_down:
			printc.Data_with_color(printc.COLOR_RED, "EVENTCHANNEL --> external_chan.New_dir_down")
			event_new_direction_down(&state, internal_chan)

		case <- external_chan.Stop:
			printc.Data_with_color(printc.COLOR_RED, "EVENTCHANNEL --> external_chan.Stop")
			event_stop(&state, internal_chan)

		case <- external_chan.New_order:

			event_new_order(&state, internal_chan)

		case <- external_chan.Too_far_up:

			event_too_far_up(&state, internal_chan)

		case <- external_chan.Too_far_down:

			event_too_far_down(&state, internal_chan)

		case <- internal_chan.breakdown_timer:

			event_engine_error(internal_chan)

	
		
		}
	}
}

func event_right_floor(state * int, internal_chan internal_channels){

	state_machine(state, RIGHT_FLOOR_E, internal_chan)
	
	printc.Data_with_color(printc.COLOR_CYAN, "reached EVENT")

	select{
	case <- internal_chan.breakdown_timer_abort:
		internal_chan.breakdown_timer_abort <- true
	default:
		internal_chan.breakdown_timer_abort <- true
	}

}

func event_stop(state * int, internal_chan internal_channels){

	driver_module.Elev_stop_engine()
	state_machine(state, STOP_E, internal_chan)
	internal_chan.open_door <- 1
	printc.Data_with_color(printc.COLOR_CYAN, "stop EVENT")

}

func event_ascend(state * int, internal_chan internal_channels){

	driver_module.Elev_start_engine(driver_module.UP)
	state_machine(state, ASCEND_E, internal_chan)
	external_chan.New_dir <- driver_module.UP
}

func event_descend(state * int, internal_chan internal_channels){


	driver_module.Elev_start_engine(driver_module.DOWN)
	state_machine(state, DESCEND_E, internal_chan)
	external_chan.New_dir <- driver_module.DOWN
}

func event_open_door(state * int, internal_chan internal_channels){

	driver_module.Elev_set_door_open_lamp(true)
	state_machine(state, OPEN_DOOR_E, internal_chan)
	go door_timer(internal_chan)
}

func event_close_door(state * int, internal_chan internal_channels, external_chan External_channels){

	driver_module.Elev_set_door_open_lamp(false)
	state_machine(state, CLOSE_DOOR_E, internal_chan)
	external_chan.Get_new_action <- 1
}

func event_new_direction_up(state * int, internal_chan internal_channels){
	
	state_machine(state, NEW_DIRECTION_UP_E, internal_chan)

	printc.Data_with_color(printc.COLOR_CYAN, "new direction EVENT")
}

func event_new_direction_down(state * int, internal_chan internal_channels){
	
	state_machine(state, NEW_DIRECTION_DOWN_E, internal_chan)

	printc.Data_with_color(printc.COLOR_CYAN, "new direction EVENT")
}

func door_timer(internal_chan internal_channels){

	select{

	case <- time.After(DOOR_TIMER):

		internal_chan.close_door <- 1 

	}

}

func breakdown_timer(internal_chan internal_channels){
	printc.Data_with_color(printc.COLOR_CYAN, "Ja men dao var me her igjen")
	
	select{
	case <- internal_chan.breakdown_timer_abort:
	default:
	}

	select{

	case <- time.After(BREAKDOWN_TIMER):

		internal_chan.breakdown_timer <- true
		printc.Data_with_color(printc.COLOR_CYAN, "No har da skjedd noke merkelig med heisen, han skulle jo ha vori her! BAD KOOOOODE!!!! CALL MAINANECECENEMENNENE" )

	case <- internal_chan.breakdown_timer_abort:
		printc.Data_with_color(printc.COLOR_CYAN, "Her er det ikke nøysamd med return")

	}

}

func event_new_order(state * int, internal_chan internal_channels){
	
	state_machine(state, NEW_ORDER_E, internal_chan)

	printc.Data_with_color(printc.COLOR_CYAN, "new direction EVENT")
}

func event_engine_error(internal_chan internal_channels){
	external_chan.Engine_error <- 1
	driver_module.Elev_stop_engine()
}

func event_too_far_up(state * int, internal_chan internal_channels){

	state_machine(state, TOO_FAR_UP_E, internal_chan)
}

func event_too_far_down(state * int, internal_chan internal_channels){

	state_machine(state, TOO_FAR_DOWN_E, internal_chan)
}