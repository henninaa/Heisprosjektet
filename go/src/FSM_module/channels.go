package FSM_module

type internal_channels struct{

	open_door chan int
	close_door chan int
	ascend chan int
	descend chan int
	stop chan int
}

func (internal_chan * internal_channels) internal_channels_FSM_init(){

	internal_chan.open_door = make(chan int,2)
	internal_chan.close_door = make(chan int,2)
	internal_chan.ascend = make(chan int,2)
	internal_chan.descend = make(chan int,2)
	internal_chan.stop  = make(chan int,2)
}

type External_channels struct{

	Right_floor chan int
	New_direction_up chan int
	New_direction_down chan int
	Stop chan int
	New_order chan int
	Get_new_direction chan int
	Get_new_action chan int
	Get_should_stop chan int
}

func (external_chan * External_channels) Init(){

	external_chan.Right_floor = make(chan int,2)
	external_chan.New_direction_up = make(chan int,2)
	external_chan.New_direction_down = make(chan int,2)
	external_chan.New_order = make(chan int,2)
	external_chan.Stop  = make(chan int,2)
	external_chan.Get_new_action = make(chan int,2)
	external_chan.Get_should_stop = make(chan int,2)
	external_chan.Get_new_direction = make(chan int,2)
}
