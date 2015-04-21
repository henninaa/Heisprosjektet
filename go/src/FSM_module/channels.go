package FSM_module

type internal_channels struct{

	open_door chan int
	close_door chan int

}

func (internal_chan * internal_channels) internal_channels_FSM_init(){

	open_door = make(chan int)
	close_door = make(chan int)
}

type External_channels struct{

	reached_floor chan int
	stop chan int
	ascend chan int
	descend chan int
	new_direction chan int

}

func (external_chan * External_channels) External_channels_FSM_init(){

	reached_floor = make(chan int)
	stop  = make(chan int)
	ascend = make(chan int)
	descend = make(chan int)
	new_direction = make(chan int)
	
}
