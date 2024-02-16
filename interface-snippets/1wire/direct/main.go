package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func execute(cmd string) {
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	output := string(out)
	//fmt.Println(output)

	/* Fetching and calculating temperature */

	splitter := strings.Split(output, " ")
	temp := splitter[len(splitter)-1]
	tempVal := strings.Split(temp, "=")[1]
	tempVal = strings.TrimRight(tempVal, "\n")

	tempFloat, err := strconv.ParseFloat(tempVal, 32)
	if err != nil {
		fmt.Println("Error parsing temperature:", err)
		return
	}

	fmt.Println("Temperature in °C: ", tempFloat/1000)

}

func main() {
	command := "cat /sys/bus/w1/devices/10-00080366ca4a/w1_slave"
	for {
		execute(command)
		time.Sleep(5 * time.Second)
	}
}
