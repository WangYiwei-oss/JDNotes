package main

import (
	"fmt"
	"time"
)

func main() {
	count := 0
	ticker := time.NewTicker(time.Second)
	for {
		fmt.Println(time.Now())
		<-ticker.C
		count++
		if count == 10 {
			ticker.Stop()
			break
		}
	}
}
