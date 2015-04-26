package FSM_module

import(
	"time"
	)

const(
	DOOR_TIMER = 500 * time.Millisecond
	FSM_INTERVAL				 	= 30 * time.Millisecond
	OBSTRUCTION_CONTROL_INTERVAL 	= 30 * time.Millisecond
	BREAKDOWN_TIMER 					= 10 * time.Second
)

//-------states
const (
	moving = 0 
	idle = 1
	initialize = 2
	door_open = 3
	)

//-------event


const EVENT_INTERVAL = 20* time.Millisecond

const(
	STOP_E 			= 0x00
	ASCEND_E 		= 0x01
	DESCEND_E		= 0x02
	OPEN_DOOR_E		= 0x03
	CLOSE_DOOR_E	= 0x04
	NEW_ORDER_E		= 0x14
	NEW_DIRECTION_UP_E = 12
	NEW_DIRECTION_DOWN_E = 11
	RIGHT_FLOOR_E	= 0x06
	TOO_FAR_DOWN_E	= 0x07
	TOO_FAR_UP_E	= 0x08
	)