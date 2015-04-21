package bank_module

type elevator_type struct{

	floor int
	direction int
	moving bool
}

func (elevator * elevator_type) elevator_type_init(){

	elevator.floor = -1
	direction = -1
	moving = false
}