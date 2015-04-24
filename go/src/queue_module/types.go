package queue_module

import(
	"driver_module"
	)

type Queue_type struct{

	Queue Queue_list
	Order_lights [4][3]bool
	Backup []queue_backup

	}

type Queue_list struct{
	List [QUEUE_SIZE]Queue_post
}

type Queue_post struct{

	Floor int
	Button_type driver_module.Elev_button_type_t

}

type queue_backup struct{

	IP string
	queue Queue_list
	floor int

}

func (queue * Queue_type) Init(current_floor int){

	var j driver_module.Elev_button_type_t

	queue.Queue.List = Init_queue()

	for i :=0; i<driver_module.N_FLOORS; i++{
		for j = 0; j<driver_module.N_BUTTONS; j++{
			queue.Order_lights[i][j] = false
			driver_module.Elev_set_button_lamp(j,i,0)
		}
	}

	queue.Queue.Get_previous_internal_queue(current_floor)
}

func (backup * queue_backup) init(ip string){

	backup.queue.List = Init_queue()
	backup.IP = ip
}
