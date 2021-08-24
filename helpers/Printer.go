package helpers

import (
	"fmt"
	"time"
)

const COUNT_TIME = 10

func StartTimePrinter() {

	count := 0
	ticker := time.NewTicker(time.Duration(COUNT_TIME) * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				count++
				fmt.Print(count * COUNT_TIME)
				fmt.Print(" seconds have passed")
				fmt.Println()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
