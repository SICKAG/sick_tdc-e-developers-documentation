/* created 12.09.2023 */
/* program for directly accessing current AIN value in raw format and scaled format */

package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func execute(cmd string) float64 {
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		fmt.Printf("%s", err)
		return 0
	}
	output := string(out)

	/* spaces removed */
	output = strings.TrimSpace(output)

	outputFloat, err := strconv.ParseFloat(output, 32)
	if err != nil {
		fmt.Println("Problem parsing: ", err)
	}
	return outputFloat
}

func main() {
	/* infinite for loop; sleeps for 5 seconds before getting new value */
	/* for AIN F */

	/* set current or voltage here */
	//mode := "voltage"
	mode := "current"

	/* set if accounting for voltage divisor */
	account := true
	//account := false

	for {
		command := "cat /sys/bus/iio/devices/iio:device1/in_voltage5_raw"
		raw := execute(command)
		command2 := "cat /sys/bus/iio/devices/iio:device1/in_voltage5_scale"
		scale := execute(command2)

		avg_clean := raw * scale
		err := 0.0059*avg_clean - 0.013
		avg_comp := avg_clean + err

		switch mode {
		case "voltage":
			if account {
				if avg_comp < 0 {
					avg_comp = 0
				}
				avg_comp = (avg_comp * 27.5) / 7.5
			}
			if !account {
				avg_comp = avg_clean
			}
			fmt.Printf("Current voltage: %f\n", avg_comp)
			break

		case "current":
			if account {
				if avg_comp < 0 {
					avg_comp = 0
				}
				avg_comp = (avg_comp * 27.5) / 7.5
				avg_comp = avg_comp * 1000 / 100.2
			}
			if !account {
				avg_comp = avg_clean
			}
			fmt.Printf("Current: %f\n", avg_comp)
			break
		}

		time.Sleep(1 * time.Second)
	}
}
