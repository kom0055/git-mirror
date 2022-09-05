package utils

import (
	"log"
	"runtime/debug"
)

func Guard() {

	if r := recover(); r != nil {
		log.Printf("recover panic: %+v,  stack: %+v",
			r, BytesToString(debug.Stack()))
	}

}

func GoRoutine(method func()) {
	if method == nil {
		return
	}
	go func() {

		defer Guard()
		method()
	}()

}
