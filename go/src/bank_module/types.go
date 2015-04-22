package bank_module

import(
	"time"
	"queue_module"
	)

type elevator_type struct{

	floor int
	direction int
	moving bool
}

type queue_backup_post struct{

	post queue_module.Queue_post
	IP string

}

func (elevator * elevator_type) elevator_type_init(){

	elevator.floor = -1
	direction = -1
	moving = false
}

const ELEVATOR_MAIN_CONTROL_INTERVAL = 20 * time.Millisecond