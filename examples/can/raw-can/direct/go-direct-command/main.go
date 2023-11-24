package main

import (
	"fmt"
	"os/exec"
)

func execute(cmd string) {
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	output := string(out[:])
	fmt.Println(output)
}

func main() {
	// get first ten CAN objects + header
	execute("candump | head -n 11")
}
