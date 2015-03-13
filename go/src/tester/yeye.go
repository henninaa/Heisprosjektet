package main

import( 
	"driver_elev"
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
	driver_elev.Elev_init();
	driver_elev.Elev_stop_engine()
	time.Sleep(1000 * time.Millisecond)

	state := 0;
	var direction int;
	cfloor := 0;

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

			if(direction == UP){
				driver_elev.Elev_start_enginge(true)
			} else{
				driver_elev.Elev_start_enginge(false)
			}
			

		case 2:
			
			cfloor<- stop_elevator_chan
			driver_elev.Elev_stop_engine();
			state = 0;
				
		}
	}
}


func queue_thread(){

	queue := queue_module.Init_queue()
	current_floor := 0
	var jejeh int
	req := false
	ggg(queue, "")

	for{
		time.Sleep(30 * time.Millisecond)
		select{

		case current_floor <- get_next_chan:

			req = true

		case jejeh = <- add_queue_chan:
			queue_module.Queue_insert(jejeh,driver_elev.BUTTON_COMMAND, current_floor, &queue)
			ggg(queue, "")

		default:

			if(req){
				
				jejeh = queue_module.One_direction(current_floor, &queue)

				if(jejeh != -1){
					send_next_chan <- jejeh
					req = false
					ggg(queue, "pop")
					
				}
			}
		}
	}
}

func sensors(){

	var current_floor int
	for{

		current_floor = Elev_get_floor_signal()

		time.Sleep(30 * time.Millisecond)
		if driver_elev.Elev_get_button_signal(2, 0) {
			add_queue_chan <- 0

		}else if driver_elev.Elev_get_button_signal(2, 1) {
			add_queue_chan <- 1

		}else if driver_elev.Elev_get_button_signal(2, 2) {
			add_queue_chan <- 2

		}else if driver_elev.Elev_get_button_signal(2, 3) {
			add_queue_chan <- 3

		}else if current_floor != -1{
			if(queue_module.Should_elevator_stop(current_floor)){
				stop_elevator_chan <- current_floor
			}
		}
	}
}



	
func ggg(queue [QUEUE_SIZE]int, ekstra string){
	Println("queue " + ekstra + ": ")
	for i:=0;i<len(queue);i++{
		//Debug_message(string(queue[i]), "queue")
		Println(queue[i])
	}
	Println("\n")
}