package main

import( 
	"driver_module"
	//"network1"
	. "fmt"
	"time"
	"queue_module"
	//. "debug_elevator"
	)

var get_next_chan = make(chan int, 1)
var send_next_chan = make(chan int, 1)
var add_queue_chan = make(chan int, 4)
var stop_elevator_chan = make(chan int, 1)

func main(){
	
	cfloor := driver_module.Elev_init();
	state := 0;
	var direction int;
	

	//network1.Network()
	go sensors()
	go queue_thread()

	for(true){
		time.Sleep(30 * time.Millisecond)

		switch (state){
		case 0:

			get_next_chan <- cfloor
			direction = <- send_next_chan
			state = 2

			if(direction == queue_module.UP){
				driver_module.Elev_start_engine(true)
			} else{
				driver_module.Elev_start_engine(false)
			}
			

		case 2:
			
			cfloor = <- stop_elevator_chan
			driver_module.Elev_stop_engine();
			state = 0;
				
		}
	}
}


func queue_thread(){

	queue := queue_module.Init_queue()
	current_floor := 0
	var current_floor_sensor int
	var jejeh int
	req := false
	ggg(queue, "")

	for{

		time.Sleep(30 * time.Millisecond)
		current_floor_sensor = driver_module.Elev_get_floor_sensor_signal()

		select{

		case current_floor = <- get_next_chan:
			req = true

		case jejeh = <- add_queue_chan:
			queue_module.Queue_insert(jejeh,driver_module.BUTTON_COMMAND, current_floor, &queue)
			ggg(queue, "")

		default:

			if(req){
				jejeh = queue_module.One_direction(current_floor, &queue)

				if(jejeh != -1){
					send_next_chan <- jejeh
					req = false
					
				}
			}
			if current_floor_sensor != -1{
				
				current_floor = current_floor_sensor

				if(queue_module.Should_elevator_stop(current_floor_sensor, &queue)){
					stop_elevator_chan <- current_floor_sensor
					ggg(queue, "pop")
				}
			}
		}
	}
}

func sensors(){

	for{

		
		time.Sleep(30 * time.Millisecond)
		if driver_module.Elev_get_button_signal(2, 0) {
			add_queue_chan <- 0

		}else if driver_module.Elev_get_button_signal(2, 1) {
			add_queue_chan <- 1

		}else if driver_module.Elev_get_button_signal(2, 2) {
			add_queue_chan <- 2

		}else if driver_module.Elev_get_button_signal(2, 3) {
			add_queue_chan <- 3

		}
	}
}



	
func ggg(queue [queue_module.QUEUE_SIZE]int, ekstra string){
	Println("queue " + ekstra + ": ")
	for i:=0;i<len(queue);i++{
		//Debug_message(string(queue[i]), "queue")
		Println(queue[i])
	}
	Println("\n")
}