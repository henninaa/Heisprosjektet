package FSM_module

import(
	"time"
	)

const(
	DOOR_OPEN_TIME 					= 3000 * time.Millisecond
	FSM_INTERVAL				 	= 30 * time.Millisecond
	OBSTRUCTION_CONTROL_INTERVAL 	= 30 * time.Millisecond
	BREAKDOWN_TIMER 					= 10 * time.Second
)