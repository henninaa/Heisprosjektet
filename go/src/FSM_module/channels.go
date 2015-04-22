package FSM_module

type internal_channels struct{

	open_door chan int
	close_door chan int
	ascend chan int
	descend chan int

}

func (internal_chan * internal_channels) internal_channels_FSM_init(){

	internal_chan.open_door = make(chan int)
	internal_chan.close_door = make(chan int)
	internal_chan.ascend = make(chan int)
	internal_chan.descend = make(chan int)
}

type External_channels struct{

	Reached_floor chan int
	Stop chan int
	New_direction chan int

}

func (external_chan * External_channels) External_channels_FSM_init(){

	external_chan.Reached_floor = make(chan int)
	external_chan.Stop  = make(chan int)
	external_chan.New_direction = make(chan int)
	
}
