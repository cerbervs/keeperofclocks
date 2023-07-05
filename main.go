package main

import (
	"fmt"
	"os"
	"time"

	"github.com/cerbervs/keeperofclocks/logging"
)

func main() {
	timeLogger := logging.TimeLogger{}
	timeLogger.Init(os.Getenv("USER"), "timelog.csv", "errorlog.log")
	defer timeLogger.File.Close()
	defer timeLogger.ErrorFile.Close()

	timeLogger.Writer.Write(
		[]string{
			"OPENED",
			timeLogger.Username,
			timeLogger.Info().Time.String(),
			time.Time{}.String(),
		},
	)

	fmt.Println("Press enter to start timer.")

	var (
		i         string
		startTime time.Time
		endTime   time.Time
	)

	fmt.Scanln(&i)

	startTime = time.Now().Local()
	endTime = startTime.Add(time.Hour * 8).Local()

	timeLogger.Writer.Write(
		[]string{
			"CLOCK IN",
			timeLogger.Username,
			startTime.String(),
			endTime.String(),
		},
	)

	fmt.Println(
		"Clocked in @ " + startTime.Format(
			"03:04PM",
		) + "; can clock out @ " + endTime.Format(
			"03:04PM",
		),
	)

	startBreak := true
	for {
		if startBreak {
			err := timeLogger.Writer.Write(
				[]string{
					"BREAK",
					timeLogger.Username,
					startTime.String(),
					endTime.String(),
				},
			)
			if err != nil {
				panic(err)
			}
			timeLogger.Writer.Flush()

			quitChan := make(chan bool, 1)
			timeChan := make(chan time.Time, 1)

			fmt.Println(
				"\nPress enter to:\tclock out,\nor r to:\tclock out and quit,\nor q to:\tquit without saving.",
			)
			fmt.Scanln(&i)
			if i == "r" {
				timeLogger.Writer.Write(
					[]string{
						"CLOCK OUT",
						timeLogger.Username,
						"End " + time.Now().Local().Format("03:04PM"),
						"End " + time.Now().Local().Format("03:04PM"),
					},
				)
				timeLogger.Writer.Flush()
				os.Exit(0)
			} else if i == "q" {
				os.Exit(0)
			}

			breakStartTime := time.Now().Local()

			go func() {
				runner(quitChan, timeChan)
			}()

			fmt.Scanln(&i)
			quitChan <- true

			resetTime := <-timeChan
			timeToLeave := endTime.Add(resetTime.Sub(breakStartTime)).Local()

			fmt.Println("You clocked back in at " + resetTime.Local().Format("03:04PM"))
			err = timeLogger.Writer.Write(
				[]string{
					"CLOCK IN",
					timeLogger.Username,
					startTime.String(),
					timeToLeave.String(),
				},
			)
			if err != nil {
				panic(err)
			}
			timeLogger.Writer.Flush()

			startBreak = false
			endTime = timeToLeave
			break
		} else {
			fmt.Println("\n\n=========\nYou can clock out at " + endTime.Local().Format("03:04PM") + "\n=========\n\n")
			startBreak = true
			continue
		}
	}
}

func runner(quitChan chan bool, timeChan chan time.Time) {
	for {
		quit := <-quitChan
		if quit {
			timeChan <- time.Now()
			break
		}
		time.Sleep(time.Second * 1)
	}
}
