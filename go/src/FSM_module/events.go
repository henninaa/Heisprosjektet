package FSM_module

import(
	"driver_module"
	"fmt"
	"time"
	"printc"
	)




func Event_generator(external_chan External_channels){

	var internal_chan internal_channels
	internal_chan.internal_channels_FSM_init()

	var state = initialize

	event_descend(&state, internal_chan)

	fmt.Println("Eventgeneratoren kjorer...")
	for{

		time.Sleep(EVENT_INTERVAL)

		select{

		case <- external_chan.Reached_floor:
printc.DataWithColor(printc.COLOR_RED, "Reached_floor EVENTCHANNEL")
			event_reached_floor(&state, internal_chan)

		case <- internal_chan.Stop:
printc.DataWithColor(printc.COLOR_RED, "intern stop EVENTCHANNEL")
			event_stop(&state, internal_chan)

		case <- internal_chan.ascend:
printc.DataWithColor(printc.COLOR_RED, "ascend EVENTCHANNEL")
			event_ascend(&state, internal_chan)

		case <- internal_chan.descend:
printc.DataWithColor(printc.COLOR_RED, "descend EVENTCHANNEL")
			event_descend(&state, internal_chan)

		case <- internal_chan.open_door:
printc.DataWithColor(printc.COLOR_RED, "open-door EVENTCHANNEL")
			event_open_door(&state, internal_chan)

		case <- internal_chan.close_door:
printc.DataWithColor(printc.COLOR_RED, "close door EVENTCHANNEL")
			event_close_door(&state, internal_chan, external_chan)

		case direction := <- external_chan.New_direction:
printc.DataWithColor(printc.COLOR_RED, "new dir EVENTCHANNEL")
			event_new_direction(&state, direction, internal_chan)

		case <- external_chan.Stop:
printc.DataWithColor(printc.COLOR_RED, "extern stop EVENTCHANNEL")
			event_stop(&state, internal_chan)
		
		case <- internal_chan.start_moving:
printc.DataWithColor(printc.COLOR_RED, "start moving EVENTCHANNEL")
			start_moving_event(&state, internal_chan)
		}
	}
}

func event_reached_floor(state * int, internal_chan internal_channels){

	state_machine(state, REACHED_FLOOR_E, internal_chan)
	printc.DataWithColor(printc.COLOR_CYAN, "reached EVENT")
}

func event_stop(state * int, internal_chan internal_channels){

	driver_module.Elev_stop_engine()
	state_machine(state, STOP_E, internal_chan)
	internal_chan.open_door <- 1
	printc.DataWithColor(printc.COLOR_CYAN, "stop EVENT")

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
	go timer(internal_chan)
}

func event_close_door(state * int, internal_chan internal_channels, external_chan External_channels){

	driver_module.Elev_set_door_open_lamp(false)
	state_machine(state, CLOSE_DOOR_E, internal_chan)
	external_chan.Req_direction <- 1
}

func start_moving_event(state * int, internal_chan internal_channels){

	 direction := <- internal_chan.dir	
			driver_module.Elev_start_engine(direction)
			state_machine(state, DESCEND_E, internal_chan)
	
		//	fmt.Println("her er no ghalt")
	
}

func event_new_direction(state * int, direction int, internal_chan internal_channels){

	if(direction == -1){return}

	select{
	case <- internal_chan.dir:
		internal_chan.dir <- direction
	default:
		internal_chan.dir <- direction
	}

	fmt.Println("direction: ", direction)
	
	state_machine(state, NEW_DIRECTION_E, internal_chan)

	printc.DataWithColor(printc.COLOR_CYAN, "new direction EVENT")
}

func timer(internal_chan internal_channels){

	select{

	case <- time.After(3 * time.Second):

		internal_chan.close_door <- 1 

	}

}