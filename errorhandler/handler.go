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

///must be called deferred
func HandlePanic(){
	if r := recover(); r != nil{
		debug.PrintStack()
		log.Println("PANIC: ", r)
	}
}


