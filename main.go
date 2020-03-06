package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

const (
	meterFilled = "▓"
	meterBlank  = "░"
)

type stats []stat
type stat struct {
	prefix string
	value  float64
}

func main() {
	for {
		c, err := cpu.Percent(0, false)
		failOnErr(err)

		m, err := mem.VirtualMemory()
		failOnErr(err)

		s, err := mem.SwapMemory()
		failOnErr(err)

		st := stats{
			stat{
				prefix: "RAM", value: m.UsedPercent,
			},
			stat{
				prefix: "SWAP", value: s.UsedPercent,
			},
		}

		if len(c) > 1 {
			for i, val := range c {
				st = append(st, stat{
					prefix: fmt.Sprintf("CPU %d", i+1),
					value:  val,
				})
			}
		} else {
			st = append(st, stat{
				prefix: "CPU",
				value:  c[0],
			})
		}

		status, err := status(st, "|")
		failOnErr(err)

		fmt.Println("\r", status)
		time.Sleep(time.Second * 5)
	}
}

func status(s stats, sep string) (combinedStatus string, err error) {
	statsLength := len(s)
	if statsLength == 0 {
		return "", errors.New("no stats supplied")
	}

	meterMax := 150 / statsLength
	var blankMeter string
	for i := 0; i < meterMax; i++ {
		blankMeter += meterBlank
	}

	for i, stat := range s {
		meter := blankMeter
		gauge := (stat.value / 100) * float64(meterMax)
		meter = strings.Replace(meter, meterBlank, meterFilled, int(gauge))
		status := fmt.Sprintf("%s: %s (%05.2f%s)", stat.prefix, meter, stat.value, "%")
		if i < (statsLength - 1) {
			status += fmt.Sprintf(" %s ", sep)
		}
		combinedStatus += status
	}

	return
}

func failOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
