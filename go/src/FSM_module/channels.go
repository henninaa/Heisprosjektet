package FSM_module

type internal_channels struct{

	open_door chan int
	close_door chan int
	ascend chan int
	descend chan int
	Stop chan int
	dir chan int
	start_moving chan int

}

func (internal_chan * internal_channels) internal_channels_FSM_init(){

	internal_chan.open_door = make(chan int,2)
	internal_chan.close_door = make(chan int,2)
	internal_chan.ascend = make(chan int,2)
	internal_chan.descend = make(chan int,2)
	internal_chan.Stop  = make(chan int,2)
	internal_chan.dir = make(chan int,2)
	internal_chan.start_moving = make(chan int,2)
}

type External_channels struct{

	Reached_floor chan int
	New_direction chan int
	Stop chan int
	Req_direction chan int
}

func (external_chan * External_channels) Init(){

	external_chan.Reached_floor = make(chan int,2)
	external_chan.New_direction = make(chan int,2)
	external_chan.Stop  = make(chan int,2)
	external_chan.Req_direction = make(chan int,2)
}
