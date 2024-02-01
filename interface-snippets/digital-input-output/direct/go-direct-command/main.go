package main

import (
	"fmt"
	"math/rand"
	"os/exec"
	"time"
)

func execute(cmd string) {
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		fmt.Printf("%s", err)
		fmt.Printf("Command output:\n%s\n", out)
	}

	output := string(out[:])
	fmt.Println(output)
}

// ON, OFF; LED is on A (496)
var commands = [2]string{"echo 1 > /sys/class/gpio/gpio496/value", "echo 0 > /sys/class/gpio/gpio496/value"}

func main() {
	// infinite loop
	for {
		// executes a random command from the slice
		// go versions before 1.20 have to use rand.Seed(time.Now().Unix())
		execute(commands[rand.Intn(len(commands))])

		time.Sleep(5 * time.Second)
	}
}
