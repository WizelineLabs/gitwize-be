package utils

import (
	"log"
	"runtime"
	"time"
)

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("\n%s took %s", name, elapsed)
}

func Trace() {
	if pc, file, line, ok := runtime.Caller(1); ok {
		fn := runtime.FuncForPC(pc)
		log.Printf("Entering: %s:%s:%d\n", file, fn.Name(), line)
		return
	}

	log.Printf("Entering: ?:?:0\n")
}

func GetFuncName() string {
	if pc, _, _, ok := runtime.Caller(1); ok {
		fn := runtime.FuncForPC(pc)
		return fn.Name()
	}

	return "?"
}
