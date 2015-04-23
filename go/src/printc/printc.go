package printc

import (
	"fmt"
	"strconv"
)

//-----------------------------------------------//

const (
	LOG_NONE 		int 	= 		iota
	LOG_ERROR 		int 	= 		iota
	LOG_WARNING		int 	= 		iota
	LOG_ALL			int 	= 		iota

	COLOR_BLACK 	int 	= 		30
	COLOR_RED 		int 	= 		31	//Error!
	COLOR_GREEN 	int 	= 		32	//Success
	COLOR_YELLOW 	int 	= 		33	//FSM_module
	COLOR_BLUE 		int 	= 		34	//Bank_module
	COLOR_MAGENTA 	int 	= 		35	//Network_module
	COLOR_CYAN 		int 	= 		36	//Queue_module
	COLOR_WHITE 	int 	= 		37	
)

var logLevel = LOG_ALL

func SetLogLevel(newLogLevel int) {
	logLevel = newLogLevel
}

//-----------------------------------------------//

func Data(values ... interface{}) {

	if logLevel >= LOG_ALL {

		for _, value := range values {
			fmt.Print(value)
			fmt.Print(" ")
		}

		fmt.Println("")
	}
}

func Data_with_color(color int, values ... interface{}) {

	if logLevel >= LOG_ALL {

		fmt.Print("\x1b[" + strconv.Itoa(color) + ";1m")

		for _, value := range values {
			fmt.Print(value)
			fmt.Print(" ")
		}

		fmt.Println("\x1b[0m") // Reset color
	}
}

//-----------------------------------------------//

func Error(values ... interface{}) {
	
	if logLevel >= LOG_ERROR {
		Data_with_color(COLOR_RED, values)
	}
}

func Warning(values ... interface{}) {
	
	if logLevel >= LOG_WARNING {
		Data_with_color(COLOR_YELLOW, values)
	}
}