package utils

import (
	"fmt"
	"log"
	"runtime"
	"time"
)

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("\n%s took %s", name, elapsed)
}

func Trace() string {
	var trace string

	if pc, file, line, ok := runtime.Caller(1); ok {
		fn := runtime.FuncForPC(pc)
		trace = fmt.Sprintf("Entering: %s:%s:%d\n", file, fn.Name(), line)
	} else {
		trace = "Entering: ?:?:0\n"
	}
	log.Println(trace)
	return trace
}

func GetFuncName() string {
	if pc, _, _, ok := runtime.Caller(1); ok {
		fn := runtime.FuncForPC(pc)
		return fn.Name()
	}

	return "?"
}
