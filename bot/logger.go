package main

import (
	"fmt"
	"log"
	"time"
)

func toLog(message string) {
	fmt.Printf("%s : %s\n", time.Now().Format(TIME_LAYOUT), message)
}

func toLogFatal(message string) {
	log.Fatalf("%s : %s\n", time.Now().Format(TIME_LAYOUT), message)
}
