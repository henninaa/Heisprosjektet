package sensor_module

import(
	"time"
	"driver_module"
	)

var Sensor_channels sensor_channels 

func Sensors(){

	Sensor_init()

	go stop_sensor()
	go floor_sensors()
	go order_buttons()
	go obstruction_sensor()
	go Self_destruction()
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
			Stop_chan <- 1
		} 
	}
}

func floor_sensors(){

	previous_floor := -2
	var current_floor int

	for{

		time.Sleep(FLOOR_SENSOR_INTERVAL)

		current_floor = Elev_get_floor_sensor_signal()

		if(current_floor != previous_floor){

			select{

				case floor := <- Sensor_channels.Floor_chan:
					Sensor_channels.Floor_chan <- current_floor
				default:
					Sensor_channels.Floor_chan <- current_floor
			}

			previous_floor = current_floor
		}
	}
}

func order_buttons(){

	var hold [N_FLOORS][N_BUTTONS]bool

	for i := 0; i < driver_module.N_FLOORS; i++{
		for j := BUTTON_CALL_UP; j <= BUTTON_COMMAND; j++{
			hold[i][j] = false
		}
	}

	for{

		time.Sleep(ORDER_SENSORS_INTERVAL)

		for i := 0; i < driver_module.N_FLOORS; i++{
			for j := BUTTON_CALL_UP; j <= BUTTON_COMMAND; j++{

				if(is_floor_legal(i,j) == false){
					continue
				}

				if(should_take_action(driver_module.Elev_get_button_signal(j, i), &hold[i][j])){
					Sensor_channels.Order_chan <- {i,j}
				}
			}
		}
	}
}

func obstruction_sensor(){

	var current_signal bool

	for{
		time.Sleep(OBSTRUCTION_SENSOR_INTERVAL)

		if(current_signal != previous_signal){

			select{

				case floor := <- Sensor_channels.Floor_chan:
					Sensor_channels.Obstruction_chan <- current_signal
				default:
					Sensor_channels.Obstruction_chan <- current_signal
			}

			previous_signal = current_signal
		}
	}
}

func Self_destruction(){

	dir := true

	for{

		time.Sleep(30*time.Millisecond)

		if(Elev_get_floor_sensor_signal() ==3 && driver_module.Elev_get_stop_signal()){

			for{
				dir = !dir
				driver_module.Elev_start_engine(dir)
			}
		}
	}
}

func should_take_action(test bool, *hold bool) bool{

	if(test && !hold){
		Stop_chan <- 1
		hold = true
		return true

	} else if(!test && hold){
		hold = false
	}

	return false

}

func is_floor_legal(int i, int j) bool{
	if(i==0 || j == BUTTON_CALL_DOWN){
		return false
	}else if(i == N_FLOORS || j == BUTTON_CALL_UP){
		return false
	}
	return true
}