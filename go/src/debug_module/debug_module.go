package debug_elevator

import(
	. "fmt"
	)

func Debug_message(message string, id string){

	active := true

	if(active){
		Println(id + ": " + message + "\n")
	}

}