package sensor_module

type sensor_channels struct {

	stop_chan chan int
	current_floor_chan chan int
	order_chan chan [2]int
	obstruction_chan chan bool

}

