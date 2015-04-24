package driver_module

import(
	"printc"
)

const N_BUTTONS = 3
const N_FLOORS = 4
const(
	UP = 	0x01
	DOWN = 	0x02
	STILL = 0x03
	)

var lamp_channel_matrix = [N_FLOORS][N_BUTTONS]int{
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4}}

var button_channel_matrix = [N_FLOORS][N_BUTTONS]int{
	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4}}

type Elev_button_type_t int;
const (
	BUTTON_CALL_UP = 0
	BUTTON_CALL_DOWN = 1
	BUTTON_COMMAND = 2
)

func Elev_init()(current_floor int){

	if io_init() {
		printc.Data_with_color(printc.COLOR_GREEN, "Elevator driver initialized...")
	}

	return 0
}

func Elev_get_floor_sensor_signal() int {

	if (io_read_bit(SENSOR_FLOOR1)){
		return 0;
	}else if (io_read_bit(SENSOR_FLOOR2)){
		return 1;
	}else if (io_read_bit(SENSOR_FLOOR3)){
		return 2;
	}else if (io_read_bit(SENSOR_FLOOR4)){
		return 3;
	}else{
		return -1;
	}
}

func Elev_get_button_signal(button Elev_button_type_t, floor int) bool{
	
	if (io_read_bit(button_channel_matrix[floor][button])){
		return true;
	}else{
		return false;
	}
}

func Elev_set_speed(speed int){

	if (speed > 0){
		io_clear_bit(MOTORDIR);
	}else if (speed < 0){
		io_set_bit(MOTORDIR);
		speed = -speed;
	}
	
	io_write_analog(MOTOR, speed);
}

func Elev_start_engine(direction int){

	if (direction == UP){
		io_clear_bit(MOTORDIR)
	}else {
		io_set_bit(MOTORDIR)
	}
	
	io_write_analog(MOTOR, MOTOR_SPEED);
}

func Elev_stop_engine(){
	io_write_analog(MOTOR, 0);
}

func Elev_set_door_open_lamp(is_open bool){
	if(is_open){
		io_set_bit(LIGHT_DOOR_OPEN);
	}else{
		io_clear_bit(LIGHT_DOOR_OPEN);
	}
}

func Elev_get_obstruction_signal() bool {
	return io_read_bit(OBSTRUCTION);
}

func Elev_get_stop_signal() bool {
	return io_read_bit(STOP);
}

func Elev_set_stop_lamp(stop bool){
	if(stop){
		io_set_bit(LIGHT_STOP);
	}else{
		io_clear_bit(LIGHT_STOP);
	}
}




func Elev_set_floor_indicator(floor int){

	if ((floor & 0x02) != 0){
		io_set_bit(LIGHT_FLOOR_IND1);
	}else{
		io_clear_bit(LIGHT_FLOOR_IND1);
	}
	if ((floor & 0x01) != 0){
		io_set_bit(LIGHT_FLOOR_IND2);
	}else{
		io_clear_bit(LIGHT_FLOOR_IND2);
	}
}



func Elev_set_button_lamp(button Elev_button_type_t, floor int, value int) {
	
	if (value == 1){
		printc.Data_with_color(printc.COLOR_RED, "Floor: ",floor,"Button:",button)
		io_set_bit(lamp_channel_matrix[floor][button]);
	}else{
		io_clear_bit(lamp_channel_matrix[floor][button]);
	}
}

func Convert_dir_to_button(direction int)(Elev_button_type_t){

	if(direction == UP){
		return BUTTON_CALL_UP
	} else if(direction == DOWN){
		return BUTTON_CALL_DOWN
	} else{
		return BUTTON_COMMAND
	}

}