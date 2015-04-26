package sensor_module

import(
	"time"
	"driver_module"
	"queue_module"
	"printc"
	)

func Sensors(sensor_chan External_channels){

	go stop_sensor(sensor_chan.Stop_chan)
	go floor_sensors(sensor_chan.Floor_chan)
	go order_buttons(sensor_chan.Order_chan, sensor_chan.Activate_orders, sensor_chan.Deactivate_orders)
	go obstruction_sensor(sensor_chan.Obstruction_chan)
	//go Self_destruction()

	printc.Data_with_color(printc.COLOR_GREEN, "Sensor_module started...")
}



func stop_sensor(stop_chan chan int){

	hold := false

	for{

		time.Sleep(STOP_SENSOR_INTERVAL)

		if(should_take_action(driver_module.Elev_get_stop_signal(), &hold)){
			stop_chan <- 1
		} 
	}
}

func floor_sensors(floor_chan chan int){

	previous_position := -2
	current_floor := -2
	var current_position int

	for{

		time.Sleep(FLOOR_SENSOR_INTERVAL)

		current_position = driver_module.Elev_get_floor_sensor_signal()

		if(current_position != previous_position){

			select{

				case  <- floor_chan:
					floor_chan <- current_position
				default:
					floor_chan <- current_position
			}

			if(current_position != -1){
				current_floor = current_position
				driver_module.Elev_set_floor_indicator(current_floor)
				//fmt.Printf("%d", current_floor)
			}

			previous_position = current_position
		}
	}
}

func Floor_sensors()int{

	return driver_module.Elev_get_floor_sensor_signal()

}

/*func order_buttons_redundant(){
	var hold [driver_module.N_FLOORS][driver_module.N_BUTTONS]bool
	var j driver_module.Elev_button_type_t
	var got_order bool
	for i := 0; i < driver_module.N_FLOORS; i++{
		for j := driver_module.BUTTON_CALL_UP; j <= driver_module.BUTTON_COMMAND; j++{
			hold[i][j] = false
		}
	}
	for{
		time.Sleep(ORDER_SENSORS_INTERVAL)
		for i := 0; i < driver_module.N_FLOORS; i++{
			for j = driver_module.BUTTON_CALL_UP; j <= driver_module.BUTTON_COMMAND; j++{
				if(is_floor_legal(i,j) == true){
					//Debug_message("Got order illegal", "order_buttons")
					continue
				}
				got_order = driver_module.Elev_get_button_signal(j, i)
				if(should_take_action_array(got_order, hold, i, int(j))){
					//Sensor_channels.Order_chan <- {i,j}
					Debug_message("Got order", "order_buttons")
					queue_module.Queue_insert(i,j,current_floor_gl)
				}
			}
		}
	}
}*/

func order_buttons(order_chan chan queue_module.Queue_post, activate_orders chan bool, deactivate_orders chan bool){

	var hold [driver_module.N_BUTTONS][driver_module.N_FLOORS]bool

	for i:= 0; i<driver_module.N_BUTTONS; i++{
		for j:= 0; j<driver_module.N_FLOORS; j++{
		
		hold[i][j] = false
		}
	}

	for{

		select{
		case <- deactivate_orders:
			<- activate_orders
		default:
			time.Sleep(ORDER_SENSORS_INTERVAL)
			check_command_orders(&(hold[0]), order_chan)
			check_up_orders(&(hold[1]), order_chan)
			check_down_orders(&(hold[2]), order_chan)

		}




	}
}

func obstruction_sensor(obstruction_chan chan bool){

	current_signal := driver_module.Elev_get_obstruction_signal()
	var previous_signal bool
	previous_signal = current_signal

	for{
		time.Sleep(OBSTRUCTION_SENSOR_INTERVAL)

		current_signal = driver_module.Elev_get_obstruction_signal()

		if(current_signal != previous_signal){

			select{

				case <- obstruction_chan:
					obstruction_chan <- current_signal
				default:
					obstruction_chan <- current_signal
			}

			previous_signal = current_signal
		}
	}
}

func self_destruction(){

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

func check_command_orders(hold * [driver_module.N_FLOORS]bool, order_chan chan queue_module.Queue_post){

	got_order := false
	var button_type driver_module.Elev_button_type_t
	var post queue_module.Queue_post

	button_type = driver_module.BUTTON_COMMAND

	for i :=0; i < driver_module.N_FLOORS; i ++{
		
		got_order = driver_module.Elev_get_button_signal(button_type, i)

		if(should_take_action(got_order, &(hold[i]))){

			post.Floor = i
			post.Button_type = button_type
			order_chan <- post


			
		}
	}
}
	
func check_up_orders(hold * [driver_module.N_FLOORS]bool, order_chan chan queue_module.Queue_post){

	got_order := false
	var button_type driver_module.Elev_button_type_t
	var post queue_module.Queue_post

	button_type = driver_module.BUTTON_CALL_UP

	for i :=0; i < driver_module.N_FLOORS -1; i ++{
		
		got_order = driver_module.Elev_get_button_signal(button_type, i)

		if(should_take_action(got_order, &(hold[i]))){

			post.Floor = i
			post.Button_type = button_type
			order_chan <- post
			
		}
	}
}

func check_down_orders(hold * [driver_module.N_FLOORS]bool, order_chan chan queue_module.Queue_post){

	got_order := false
	var button_type driver_module.Elev_button_type_t
	var post queue_module.Queue_post

	button_type = driver_module.BUTTON_CALL_DOWN

	for i :=1; i < driver_module.N_FLOORS; i ++{
		
		got_order = driver_module.Elev_get_button_signal(button_type, i)

		if(should_take_action(got_order, &(hold[i]))){

			post.Floor = i
			post.Button_type = button_type
			order_chan <- post
			
		}
	}
}