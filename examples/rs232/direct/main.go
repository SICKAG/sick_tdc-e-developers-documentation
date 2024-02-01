package main

import (
	"fmt"
	"math/rand"
	"os/exec"
	"sync"
	"time"
)

func execute(cmd string) {
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	output := string(out)
	if cmd == "cat /dev/ttymxc5 | head -n 1" {
		fmt.Println("Data received: ", output)
	}
}

func main() {

	var wg sync.WaitGroup
	wg.Add(2)

	/* Goroutine for fetching data from device */
	go func() {
		defer wg.Done()
		/* ttymxc5 is used for rs-232 data */
		/* implemented to fetch the last data line that is sent to the device */
		command := "cat /dev/ttymxc5 | head -n 1"
		for {
			execute(command)
		}
	}()

	/* Goroutine for sending data to device */

	go func() {
		defer wg.Done()
		datasend := []string{"hello", "world", "go", "thank_you", "fun", "rs232", "test-data", "smile", "good-day", "good-evening", "good-night", "SICK", "Mobilisis", "TDC-E"}
		for {
			/* sends data every three seconds */
			command := "echo " + datasend[rand.Intn(len(datasend))] + " > /dev/ttymxc5"
			execute(command)
			time.Sleep(3 * time.Second)
		}
	}()

	wg.Wait()
}
