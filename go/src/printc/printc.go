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
	COLOR_RED 		int 	= 		31
	COLOR_GREEN 	int 	= 		32
	COLOR_YELLOW 	int 	= 		33
	COLOR_BLUE 		int 	= 		34
	COLOR_MAGENTA 	int 	= 		35
	COLOR_CYAN 		int 	= 		36
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

func DataWithColor(color int, values ... interface{}) {

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
		DataWithColor(COLOR_RED, values)
	}
}

func Warning(values ... interface{}) {
	
	if logLevel >= LOG_WARNING {
		DataWithColor(COLOR_YELLOW, values)
	}
}