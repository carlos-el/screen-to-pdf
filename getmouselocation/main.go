package main

import (
	"fmt"
	"time"

	"github.com/go-vgo/robotgo"
)

func main() {
	fmt.Printf("Returning mouse location every second. Halt the program to stop...\n")
	for {
		time.Sleep(1 * time.Second)
		x, y := robotgo.Location()
		fmt.Printf("%d, %d\n", x, y)
	}
}
