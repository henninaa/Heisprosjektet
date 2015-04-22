package queue_module

import(
	"driver_module"
	)

type Queue_type struct{

	queue queue_list
	order_lights [4][3]bool
	backup []queue_backup

	}

type queue_list struct{
	list [QUEUE_SIZE]Queue_post
}

type Queue_post struct{

	Floor int
	Button_type driver_module.Elev_button_type_t

}

type queue_backup struct{

	IP string
	queue queue_list
	floor int

}

func (queue * Queue_type) queue_type_init(){

	var j driver_module.Elev_button_type_t

	queue.queue.list = Init_queue()

	for i :=0; i<driver_module.N_FLOORS; i++{
		for j = 0; j<driver_module.N_BUTTONS; j++{
			queue.order_lights[i][j] = false
			driver_module.Elev_set_button_lamp(j,i,0)
		}
	}
}

func (backup * queue_backup) init(ip string){

	backup.queue.list = Init_queue()
	backup.IP = ip
}
