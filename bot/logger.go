package main

import (
	"fmt"
	"log"
	"time"
)

func toLog(message string) {
	fmt.Printf("%s : %s\n", time.Now().Format(TimeLayout), message)
}

func toLogF(format string, a ...interface{}) {
	fm := fmt.Sprintf(format, a...)
	fmt.Printf("%s : %s\n", time.Now().Format(TimeLayout), fm)
}

func toLogFatal(message string) {
	log.Fatalf("%s : %s\n", time.Now().Format(TimeLayout), message)
}

func toLogFatalF(format string, a ...interface{}) {
	fm := fmt.Sprintf(format, a...)
	log.Fatalf("%s : %s\n", time.Now().Format(TimeLayout), fm)
}
