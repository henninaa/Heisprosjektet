package sensor_module

type sensor_channels struct {

	Stop_chan chan int
	Floor_chan chan int
	Order_chan chan [2]int
	Obstruction_chan chan bool

}

