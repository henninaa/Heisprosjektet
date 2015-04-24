package queue_module

import (
	"os"
	"math"
	"printc"
	"encoding/gob"
	"driver_module"
	"network_module"
)


const QUEUE_SIZE = 12

func (queue_class * Queue_type) Insert_to_own_queue(post Queue_post, current_floor int){

	insert_floor := post.Floor
	insert_type := post.Button_type

	queue_class.Queue.insert(insert_floor, insert_type, current_floor, true)

	queue_class.Queue.save_internal_orders_to_file()
	printc.Data_with_color(printc.COLOR_MAGENTA, "I HAVE JUST WITTEN TO YOUR STUPID FILE!!!")

}

func (queue_class * Queue_type) Insert_to_backup_queue(insert_floor int, insert_type driver_module.Elev_button_type_t, ip string){

	for i:= 0; i<len(queue_class.Backup); i++{
		if(queue_class.Backup[i].IP == ip){
			queue_class.Backup[i].queue.insert(insert_floor, insert_type, queue_class.Backup[i].floor, false)
			return
		}
	}
	
	queue_class.add_backup(ip)
	queue_class.Insert_to_backup_queue(insert_floor, insert_type, ip)

}

func (queue_class * Queue_type) Annex_backup(IP string, current_floor int){

	var backup_ext_orders []Queue_post

	for i:= 0; i<len(queue_class.Backup); i++{
		if(queue_class.Backup[i].IP == IP){
			
			for j := range queue_class.Backup[i].queue.List{

				if(queue_class.Backup[i].queue.List[j].Button_type != driver_module.BUTTON_COMMAND){
					backup_ext_orders = append(backup_ext_orders, queue_class.Backup[i].queue.List[j])
				}
			}

			queue_class.remove_backup(i)
			break
		}
	}

	for _, order := range backup_ext_orders{

		IP = queue_class.Get_lowest_cost_ip(order, current_floor)

		if(IP == "self"){
			queue_class.Queue.insert(order.Floor, order.Button_type, current_floor, true)
		}

	}

}

func (queue * Queue_type) New_backup_floor(floor int, ip string){

		for i:= 0; i<len(queue.Backup); i++{
		if(queue.Backup[i].IP == ip){
			queue.Backup[i].floor = floor
			return
		}
	}
	
	queue.add_backup(ip)
	queue.New_backup_floor(floor, ip)

}

func Init_queue()(queue [12]Queue_post){

	for i := 0; i < QUEUE_SIZE; i++ {
			queue[i].Floor = -1
			queue[i].Button_type = driver_module.BUTTON_COMMAND
	}

	return queue

}

func (queue_class * Queue_type) Get_new_direction(current_floor int) int{
	if(queue_class.Queue.List[0].Floor==-1 ){
		return -1
	}else if(current_floor < queue_class.Queue.List[0].Floor){
		return driver_module.UP
	} else{
		return driver_module.DOWN
	}

}

func (queue_class * Queue_type) Should_elevator_stop(current_floor int, direction driver_module.Elev_button_type_t, post * Queue_post) bool{

	if(current_floor == queue_class.Queue.List[0].Floor && queue_class.Queue.List[0].Button_type == driver_module.BUTTON_COMMAND){
		turn_off_lights(current_floor, queue_class.Queue.List[0].Button_type)
		post.Floor = current_floor
		post.Button_type = queue_class.Queue.List[0].Button_type
		queue_class.Queue.queue_remove_multiple_floors(queue_class.Queue.List[0], true)
		
		return true
	}else if(current_floor == queue_class.Queue.List[0].Floor && direction == driver_module.BUTTON_COMMAND){
		turn_off_lights(current_floor, queue_class.Queue.List[0].Button_type)
		post.Floor = current_floor
		post.Button_type = queue_class.Queue.List[0].Button_type
		queue_class.Queue.queue_remove_multiple_floors(queue_class.Queue.List[0], true)
		
		return true
	} else if(current_floor == queue_class.Queue.List[0].Floor && (queue_class.Queue.List[0].Button_type == direction || queue_class.Queue.List[1].Floor == -1)){
		turn_off_lights(current_floor, queue_class.Queue.List[0].Button_type)
		post.Floor = current_floor
		post.Button_type = queue_class.Queue.List[0].Button_type
		queue_class.Queue.queue_remove_multiple_floors(queue_class.Queue.List[0], true)

		return true
	}
	return false
}

func (queue * Queue_type) Remove_post_from_backup_queue(post Queue_post, ip string){


	for i:= 0; i<len(queue.Backup); i++{
		if(queue.Backup[i].IP == ip){

			queue.Backup[i].queue.queue_remove_multiple_floors(post, false)
			

		}
	}
}

func (queue * Queue_type) Get_lowest_cost_ip(insert_post Queue_post, current_floor int) (lowest_cost_ip string){
	lowest_cost := 999999
	cost := -1

	for i := range queue.Backup{

		cost = queue.Backup[i].queue.get_insertion_cost(insert_post.Floor, insert_post.Button_type, queue.Backup[i].floor)

		printc.Data_with_color(printc.COLOR_BLACK, "				IP: ", queue.Backup[i].IP)
		printc.Data_with_color(printc.COLOR_BLACK, "				COST: ", cost)
		
		if(cost < lowest_cost){
			lowest_cost_ip = queue.Backup[i].IP
			lowest_cost = cost
		}

	}
	printc.Data_with_color(printc.COLOR_BLACK, "Antall: ", len(queue.Backup))

	cost = queue.Queue.get_insertion_cost(insert_post.Floor, insert_post.Button_type, current_floor)

	if(cost <= lowest_cost){
		lowest_cost_ip = "self"
	}
	
	printc.Data_with_color(printc.COLOR_BLACK, "				SelfIP: ")
	printc.Data_with_color(printc.COLOR_BLACK, "				COST: ", cost)

	printc.Data_with_color(printc.COLOR_BLACK, "best ip is: ", lowest_cost_ip)

	return lowest_cost_ip
}

func Convert_mail_to_queue_post(mail network_module.Mail)(order Queue_post){

	order.Floor = mail.Msg.Floor
	order.Button_type = mail.Msg.Dir

	return order
}


/*func Convert_mail_to_backup_post(mail network_module.Mail)(post Queue_post){

	return post
}
*/
func (queue * Queue_type) Remove_backup_with_IP(ip string){
	
	for i:= 0; i<len(queue.Backup); i++{
		if(queue.Backup[i].IP == ip){
			queue.Backup = append(queue.Backup[:i], queue.Backup[i+1:]...)
		}
	}
}

func (internal_orders * Queue_list) Get_previous_internal_queue(current_floor int) {
	var backup_from_file Queue_list
	backup_from_file.get_internal_orders_from_file()

	printc.Data_with_color(printc.COLOR_MAGENTA, "READING FROM FILE!!!!!!")

	for _, order := range backup_from_file.List{

		if(order.Button_type == driver_module.BUTTON_COMMAND && order.Floor != -1){

			internal_orders.insert(order.Floor, order.Button_type, current_floor, true)
		}
	}
}

func (queue * Queue_type) remove_backup(pos int){

	queue.Backup = append(queue.Backup[:pos], queue.Backup[pos+1:]...)

}

func (queue * Queue_type) add_backup(ip string)(){

	var new_backup queue_backup

	new_backup.IP = ip
	new_backup.floor = -1
	new_backup.queue.List = Init_queue()

	pos := len(queue.Backup)

	queue.Backup = append(queue.Backup[:pos], new_backup)
}


func (queue * Queue_list) insert_to_pos(insert_post Queue_post, position int){

	var swap Queue_post
	swap = insert_post
	var swap_tmp Queue_post
	

	for i := position; i < QUEUE_SIZE; i++ {
		
		if(queue.List[i].Floor == -1){
			queue.List[i] = swap
			break
		}

		swap_tmp = queue.List[i]
		queue.List[i] = swap
		swap = swap_tmp

	}

}

func (queue * Queue_list) queue_remove_multiple_floors(post Queue_post, local_order bool){

	previndex :=0
	find :=-2
	local_order = true
	
	for i := 0; i < QUEUE_SIZE; i++ {
		printc.Data_with_color(printc.COLOR_MAGENTA, "QUEUE_SIZE: ", QUEUE_SIZE, "\ni: ", i, "\nqueue.List[i].Floor: ", queue.List[i].Floor, "find: ", find)

		queue.List[previndex] = queue.List[i]

		if(i == QUEUE_SIZE-1){
			if(queue.List[i].Floor == post.Floor && (queue.List[i].Button_type == post.Button_type || queue.List[i].Button_type == driver_module.BUTTON_COMMAND)){
				if(local_order){turn_off_lights(queue.List[i].Floor, queue.List[i].Button_type)}
				continue
			}else{
					previndex++
					continue
			}
		}

		if(queue.List[i].Floor == post.Floor){

			if(queue.List[i].Button_type == post.Button_type){
				if(local_order){turn_off_lights(queue.List[i].Floor, queue.List[i].Button_type)}
				printc.Data_with_color(printc.COLOR_CYAN, "Treff: 1")
				find = i +1
				//if(i<4){queue.check_for_delete_all(i)}
				continue
			} else if(queue.List[i].Button_type == driver_module.BUTTON_COMMAND){
				printc.Data_with_color(printc.COLOR_CYAN, "Treff: 2")
				if(local_order){turn_off_lights(queue.List[i].Floor, queue.List[i].Button_type)}
				continue
			} else if(find == -2){
			}else if(queue.List[i].Floor < queue.List[find].Floor && queue.List[i].Button_type == driver_module.BUTTON_CALL_UP){
				printc.Data_with_color(printc.COLOR_CYAN, "Treff: 4")
				if(local_order){turn_off_lights(queue.List[i].Floor, queue.List[i].Button_type)}
				continue
			} else if(queue.List[i].Floor > queue.List[find].Floor && queue.List[i].Button_type == driver_module.BUTTON_CALL_DOWN){
				printc.Data_with_color(printc.COLOR_CYAN, "Treff: 5")
				if(local_order){turn_off_lights(queue.List[i].Floor, queue.List[i].Button_type)}
				continue
			}
		}

		previndex++
	}
	if(previndex < QUEUE_SIZE){for i := previndex; i< QUEUE_SIZE; i++ {queue.List[i].Floor = -1}}
	

	queue.save_internal_orders_to_file()

}

func (queue * Queue_list) check_for_delete_all(first_hit int){

	if(first_hit<3 && queue.List[3].Floor==-1 && queue.List[0].Floor==queue.List[1].Floor){
		if(first_hit < 2 && queue.List[2].Floor==-1){
			for i:=0;i<2;i++{
				turn_off_lights(queue.List[i].Floor, queue.List[i].Button_type)
				queue.List[i].Floor=-1
			}
		}else if(first_hit == 2 && queue.List[1].Floor==queue.List[2].Floor){
			for i:=0;i<3;i++{
				turn_off_lights(queue.List[i].Floor, queue.List[i].Button_type)
				queue.List[i].Floor=-1
			}
		}
	}
}

func turn_on_light(insert_floor int, insert_type driver_module.Elev_button_type_t){

	driver_module.Elev_set_button_lamp(insert_type, insert_floor, 1)

}

func turn_off_lights(floor int, itype driver_module.Elev_button_type_t){

	driver_module.Elev_set_button_lamp(itype, floor, 0)

}

func (queue * Queue_list) insert(insert_floor int, insert_type driver_module.Elev_button_type_t, current_floor int, local_order bool){

	var input_post Queue_post
	input_post.Floor = insert_floor
	input_post.Button_type = insert_type

	if(insert_type != driver_module.BUTTON_COMMAND && local_order == false){local_order = true}

	if(local_order){turn_on_light(insert_floor, insert_type)}

	printc.Data_with_color(printc.COLOR_GREEN, "Queue_insert: Got queue insert ", insert_floor, current_floor)

	prev := current_floor
	var direction driver_module.Elev_button_type_t

	for i := range queue.List {

		if(queue.List[i] == input_post){
			break

		}else if(queue.List[i].Floor == -1){
			queue.List[i] = input_post
			break

		}else if(prev < queue.List[i].Floor){
			direction = driver_module.BUTTON_CALL_UP

			if(insert_floor < queue.List[i].Floor && insert_floor > prev){
				if(insert_type == direction || insert_type == driver_module.BUTTON_COMMAND){
				queue.insert_to_pos(input_post, i)
				break
				}	
			}


		} else{
			direction = driver_module.BUTTON_CALL_DOWN
		
			if(insert_floor > queue.List[i].Floor && insert_floor < prev){
				if(insert_type == direction || insert_type == driver_module.BUTTON_COMMAND){
				queue.insert_to_pos(input_post, i)
				break
				}	
			}
		}

		prev = queue.List[i].Floor

	}
	printc.Data_with_color(printc.COLOR_GREEN, "Queue_insert: Ferdig!")
}

func (queue * Queue_list) get_insertion_cost(insert_floor int, insert_type driver_module.Elev_button_type_t, current_floor int)(int){


	var input_post Queue_post
	input_post.Floor = insert_floor
	input_post.Button_type = insert_type
	cost:=0

	prev := current_floor
	var direction driver_module.Elev_button_type_t

	for i := range queue.List {

		if(queue.List[i].Floor == input_post.Floor){
			if((queue.List[i].Button_type != input_post.Button_type) && queue.List[i].Button_type!=driver_module.BUTTON_COMMAND && input_post.Button_type!=driver_module.BUTTON_COMMAND){

				cost += 1

			} else if(queue.List[i].Button_type == input_post.Button_type){

				break
			}

		}else if(queue.List[i].Floor == -1){
			break

		}else if(prev < queue.List[i].Floor){
			direction = driver_module.BUTTON_CALL_UP

			if(insert_floor < queue.List[i].Floor && insert_floor > prev){
				if(insert_type == direction || insert_type == driver_module.BUTTON_COMMAND){
				break
				}	
			}


		} else{
			direction = driver_module.BUTTON_CALL_DOWN
		
			if(insert_floor > queue.List[i].Floor && insert_floor < prev){
				if(insert_type == direction || insert_type == driver_module.BUTTON_COMMAND){
				break
				}	
			}
		}

		cost += int(math.Abs(float64(prev - queue.List[i].Floor)))
		prev = queue.List[i].Floor

	}

	cost += int(math.Abs(float64(prev - insert_floor)))

	return cost
}

func (internal_queue * Queue_list) save_internal_orders_to_file(){
	file, err := os.Create("internal_orders.gob")

	if err != nil {
		panic(err)		
	}

	enc := gob.NewEncoder(file)

	err = enc.Encode(&internal_queue)

	if err != nil {
		printc.Data_with_color(printc.COLOR_RED,"Encode error while writing to file: ", err)
		os.Exit(1)
	}

	file.Close()

}

func (internal_queue * Queue_list) get_internal_orders_from_file(){
		
	if _,err := os.Open("internal_orders.gob"); os.IsNotExist(err) {
		printc.Data_with_color(printc.COLOR_MAGENTA,"Have some smartass studass deleted our file? You moron!!!")
	}else{
		file, err := os.Open("internal_orders.gob")

		dec := gob.NewDecoder(file)

		err = dec.Decode(&internal_queue)
		printc.Data_with_color(printc.COLOR_MAGENTA,"I AM A STUPID READFUNCTION WHO CANNOT READ!!")
		if err != nil {
			printc.Data_with_color(printc.COLOR_MAGENTA,"Decode error while reading from file: ", err)
			return
		}
		printc.Data_with_color(printc.COLOR_RED, "Decoded file: ", internal_queue)
		file.Close()
	}
}
