package FSM_module

import(
	"driver_module"
	"fmt"
	"time"
	"printc"
	)

var external_chan External_channels


func Event_generator(external_chan_in External_channels){

	external_chan = external_chan_in

	var internal_chan internal_channels
	internal_chan.internal_channels_FSM_init()

	var state = idle

	fmt.Println("Eventgeneratoren kjorer...")
	for{

		time.Sleep(EVENT_INTERVAL)

		select{

		case <- external_chan.Right_floor:
printc.DataWithColor(printc.COLOR_RED, "Right_floor EVENTCHANNEL")
			event_right_floor(&state, internal_chan)

		case <- internal_chan.stop:
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

		case <- external_chan.New_direction_up:
printc.DataWithColor(printc.COLOR_RED, "new dir UP EVENTCHANNEL")
			event_new_direction_up(&state, internal_chan)

		case <- external_chan.New_direction_down:
printc.DataWithColor(printc.COLOR_RED, "new dir DOWN EVENTCHANNEL")
			event_new_direction_down(&state, internal_chan)

		case <- external_chan.Stop:
printc.DataWithColor(printc.COLOR_RED, "extern stop EVENTCHANNEL")
			event_stop(&state, internal_chan)

		case <- external_chan.New_order:

			event_new_order(&state, internal_chan)
		
		}
	}
}

func event_right_floor(state * int, internal_chan internal_channels){

	state_machine(state, RIGHT_FLOOR_E, internal_chan)
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
	go door_timer(internal_chan)
}

func event_close_door(state * int, internal_chan internal_channels, external_chan External_channels){

	driver_module.Elev_set_door_open_lamp(false)
	state_machine(state, CLOSE_DOOR_E, internal_chan)
	external_chan.Get_new_action <- 1
}

func event_new_direction_up(state * int, internal_chan internal_channels){
	
	state_machine(state, NEW_DIRECTION_UP_E, internal_chan)

	printc.DataWithColor(printc.COLOR_CYAN, "new direction EVENT")
}

func event_new_direction_down(state * int, internal_chan internal_channels){
	
	state_machine(state, NEW_DIRECTION_DOWN_E, internal_chan)

	printc.DataWithColor(printc.COLOR_CYAN, "new direction EVENT")
}

func door_timer(internal_chan internal_channels){

	select{

	case <- time.After(DOOR_TIMER):

		internal_chan.close_door <- 1 

	}

}

func event_new_order(state * int, internal_chan internal_channels){
	
	state_machine(state, NEW_ORDER_E, internal_chan)

	printc.DataWithColor(printc.COLOR_CYAN, "new direction EVENT")
}