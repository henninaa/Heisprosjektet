package sensor_module

import(
	"time"
	"driver_module"
	"queue_module"
	. "debug_module"
	)

var Sensor_channels sensor_channels
var current_floor_gl int

func Sensors(){

	Sensor_init()

	go stop_sensor()
	go floor_sensors()
	go order_buttons()
	go obstruction_sensor()
	go Self_destruction()

	Debug_message("Sensors started...", "SENSORS")
}

func Sensor_init(){

	Sensor_channels.Stop_chan = make(chan int, 1)
	Sensor_channels.Floor_chan = make(chan int, 1)
	Sensor_channels.Order_chan = make(chan [2]int, 12)
	Sensor_channels.Obstruction_chan = make(chan bool, 1)

}

func stop_sensor(){

	hold := false

	for{

		time.Sleep(STOP_SENSOR_INTERVAL)

		if(should_take_action(driver_module.Elev_get_stop_signal(), &hold)){
			Sensor_channels.Stop_chan <- 1
		} 
	}
}

func floor_sensors(){

	previous_floor := -2
	var current_floor int

	for{

		time.Sleep(FLOOR_SENSOR_INTERVAL)

		current_floor = driver_module.Elev_get_floor_sensor_signal()

		if(current_floor != previous_floor){

			select{

				case  <- Sensor_channels.Floor_chan:
					Sensor_channels.Floor_chan <- current_floor
				default:
					Sensor_channels.Floor_chan <- current_floor
			}

			if(current_floor != -1){
				current_floor_gl = current_floor
				driver_module.Elev_set_floor_indicator(current_floor)
			}

			previous_floor = current_floor
		}
	}
}

func Floor_sensors()int{

	return driver_module.Elev_get_floor_sensor_signal()

}

func order_buttons(){

	var hold [driver_module.N_FLOORS][driver_module.N_BUTTONS]bool
	var j driver_module.Elev_button_type_t

	for i := 0; i < driver_module.N_FLOORS; i++{
		for j := driver_module.BUTTON_CALL_UP; j <= driver_module.BUTTON_COMMAND; j++{
			hold[i][j] = false
		}
	}

	for{

		time.Sleep(ORDER_SENSORS_INTERVAL)

		for i := 0; i < driver_module.N_FLOORS; i++{
			for j = driver_module.BUTTON_CALL_UP; j <= driver_module.BUTTON_COMMAND; j++{

				if(is_floor_legal(i,j) == false){
					continue
				}

				if(should_take_action(driver_module.Elev_get_button_signal(j, i), &hold[i][j])){
					//Sensor_channels.Order_chan <- {i,j}
					queue_module.Queue_insert(i,j,current_floor_gl)

				}
			}
		}
	}
}

func obstruction_sensor(){

	current_signal := driver_module.Elev_get_obstruction_signal()
	var previous_signal bool
	previous_signal = current_signal

	for{
		time.Sleep(OBSTRUCTION_SENSOR_INTERVAL)

		if(current_signal != previous_signal){

			select{

				case <- Sensor_channels.Floor_chan:
					Sensor_channels.Obstruction_chan <- current_signal
				default:
					Sensor_channels.Obstruction_chan <- current_signal
			}

			previous_signal = current_signal
		}
	}
}

func Self_destruction(){

	for{

		time.Sleep(30*time.Millisecond)

		if(driver_module.Elev_get_floor_sensor_signal() ==3 && driver_module.Elev_get_stop_signal()){

			for{
				
				driver_module.Elev_start_engine(driver_module.UP)
				driver_module.Elev_start_engine(driver_module.DOWN)
			}
		}
	}
}

func should_take_action(test bool, hold * bool) bool{

	if(test && !(*hold)){
		Sensor_channels.Stop_chan <- 1
		*hold = true
		return true

	} else if(!test && *hold){
		*hold = false
	}

	return false

}

func is_floor_legal(i int, j driver_module.Elev_button_type_t) bool{
	if(i==0 || j == driver_module.BUTTON_CALL_DOWN){
		return false
	}else if(i == driver_module.N_FLOORS || j == driver_module.BUTTON_CALL_UP){
		return false
	}
	return true
}