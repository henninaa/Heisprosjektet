package FSM_module

type internal_channels struct{

	open_door chan int
	close_door chan int
	ascend chan int
	descend chan int
	stop chan int
	breakdown_timer chan bool
	breakdown_timer_abort chan bool
}

func (internal_chan * internal_channels) internal_channels_FSM_init(){

	internal_chan.open_door = make(chan int,2)
	internal_chan.close_door = make(chan int,2)
	internal_chan.ascend = make(chan int,2)
	internal_chan.descend = make(chan int,2)
	internal_chan.stop  = make(chan int,2)
	internal_chan.breakdown_timer = make(chan bool, 10)
	internal_chan.breakdown_timer_abort = make(chan bool, 10)
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
	Engine_error chan int
	Abort_light_show chan bool
	Too_far_down chan int
	Too_far_up chan int
	New_dir chan int
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
	external_chan.Engine_error = make(chan int, 2)
	external_chan.Abort_light_show = make(chan bool)
	external_chan.Too_far_up = make(chan int, 2)
	external_chan.Too_far_down = make(chan int, 2)
	external_chan.New_dir = make(chan int, 2)
}
