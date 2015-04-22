package sensor_module

import("queue_module")

type External_channels struct {

	stop_chan chan int
	floor_chan chan int
	order_chan chan queue_module.Queue_post
	obstruction_chan chan bool

}

func (sensor_chan * External_channels) Sensor_init(){

	sensor_chan.stop_chan = make(chan int, 1)
	sensor_chan.floor_chan = make(chan int, 1)
	sensor_chan.order_chan = make(chan queue_module.Queue_post)
	sensor_chan.obstruction_chan = make(chan bool, 1)

}