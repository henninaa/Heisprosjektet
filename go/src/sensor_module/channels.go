package sensor_module

import("queue_module")

type External_channels struct {

	Stop_chan chan int
	Floor_chan chan int
	Order_chan chan queue_module.Queue_post
	Obstruction_chan chan bool

}

func (sensor_chan * External_channels) Init(){

	sensor_chan.Stop_chan = make(chan int, 1)
	sensor_chan.Floor_chan = make(chan int, 1)
	sensor_chan.Order_chan = make(chan queue_module.Queue_post, 40)
	sensor_chan.Obstruction_chan = make(chan bool, 1)

}