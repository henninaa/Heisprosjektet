package FSM_module

import(
	"sensor_module"
	"queue_module"
	)



func FSM(){

	var direction int
	var current_floor int
	var state _state
	state = init

	stop_control()
	obstruction_control()

	for{

		switch (state){
		case init:

			current_floor = <- sensor_module.Sensor_channels.Floor_chan

			if (current_floor != -1){
				state = still
			}

		case still:

			direction = queue_module.One_direction(current_floor)

		case moving:

			


		}

	}




}


func stop_control(){



}

func obstruction_control(){



}