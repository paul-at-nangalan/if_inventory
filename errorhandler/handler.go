package errorhandler

import (
	"log"
	"runtime/debug"
)

func PanicOnErr(err error){
	if err != nil{
		debug.PrintStack()
		log.Panicln("ERROR: ", err)
	}
}


