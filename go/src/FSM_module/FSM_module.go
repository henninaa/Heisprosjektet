package FSM_module

import(
	"driver_module"
	"sensor_module"
	"queue_module"
	"time"
	. "debug_module"
	)

var direction int
var current_floor int
var state _state
var obstruction = false

func FSM(){

	var current_position int
	direction = -1

	state = initialize

	Debug_message("FSM started...", "FSM")
	
	for{

		time.Sleep(FSM_INTERVAL)

		stop_control()
		obstruction_control()

		if(obstruction){continue}

		switch (state){
		case initialize:

			current_floor = <- sensor_module.Sensor_channels.Floor_chan

			Debug_message(string(current_floor), "FSM")

			if (current_floor != -1){
				driver_module.Elev_stop_engine()
				state = still
				Debug_message("gar til stille", "FSM")

			} else if (direction == -1){
				direction = driver_module.DOWN
				driver_module.Elev_start_engine(direction)
			}

		case still:

			direction = queue_module.One_direction(current_floor)

			if (queue_module.Should_elevator_stop(current_floor)){
				driver_module.Elev_set_door_open_lamp(true)
				time.Sleep(DOOR_OPEN_TIME)
				state = door_open
				Debug_message("gar til door_open", "FSM")

			} else if (direction != -1){
				driver_module.Elev_start_engine(direction)
				state = moving
				Debug_message("gar til moving", "FSM")
			}

		case door_open:
			
			if(false == obstruction){
				driver_module.Elev_set_door_open_lamp(false)
				state = still
				Debug_message("gar til stille", "FSM")
			}

		case moving:
			select{
			case current_position = <- sensor_module.Sensor_channels.Floor_chan:

				if(current_position != -1){
					current_floor = current_position

					if(queue_module.Should_elevator_stop(current_floor)){
						driver_module.Elev_stop_engine()
						direction = -1
						driver_module.Elev_set_door_open_lamp(true)
						time.Sleep(DOOR_OPEN_TIME)
						state = door_open
						Debug_message("gar til door_open", "FSM")
					}
				}
			default: 
				continue
			}
		}

	}

}

func stop_control(){

	select{

	case <- sensor_module.Sensor_channels.Stop_chan:

		if(direction != -1){
			driver_module.Elev_stop_engine()
		}
		driver_module.Elev_set_stop_lamp(true)
		driver_module.Elev_set_door_open_lamp(true)
		<- sensor_module.Sensor_channels.Stop_chan
		driver_module.Elev_set_door_open_lamp(false)
		driver_module.Elev_set_stop_lamp(false)

		if(direction != -1){
			driver_module.Elev_start_engine(direction)
		}

	default:
		return
	}

}

func obstruction_control(){

	select{
	case msg := <- sensor_module.Sensor_channels.Obstruction_chan:

		if (msg == true){
			obstruction = true
			driver_module.Elev_stop_engine()
			Debug_message("treff", "obstruction_control")
		}else if(obstruction == true){
			obstruction = false
			if(direction != -1){
				driver_module.Elev_start_engine(direction)
			}
		}

	default:
		return
	}

}