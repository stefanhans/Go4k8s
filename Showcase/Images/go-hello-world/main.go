package main

import (
	"fmt"
	"time"
)

func main() {
	version := "0.0.1"
	fmt.Printf("\n%s: Go's Hello World Version %s\n\n", time.Now().Format(time.RFC1123), version)

	//time.Sleep(time.Hour)
}
