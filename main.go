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

	/* r, err := timeLogger.Reader.ReadAll()
	if err != nil {
		timeLogger.Error(errors.New("read " + err.Error()))
	} else if r == nil || len(r) == 0 {
		timeLogger.Writer.Write(
			[]string{"username", "time", "time_to_leave"},
		)
	} */

	timeLogger.Writer.Write(
		[]string{
			timeLogger.Username,
			timeLogger.Info().Time.String(),
			time.Now().Local().Add(time.Hour * 8).Format("03:04PM"),
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

	fmt.Println(
		"Clocked in @ " + startTime.Format(
			"03:04PM",
		) + "; can clock out @ " + endTime.Local().
			Format("03:04PM"),
	)

	start := true
	var resetTime time.Time
	for {
		if start {
			err := timeLogger.Writer.Write(
				[]string{
					timeLogger.Username,
					startTime.String(),
					endTime.String(),
				},
			)
			if err != nil {
				panic(err)
			}
			timeLogger.Writer.Flush()

			quitChan := make(chan bool)
			timeChan := make(chan time.Time)

			fmt.Println(
				"\nPress enter to:\tclock out,\nor r to:\tclock out and quit,\nor q to:\tquit without saving.",
			)
			fmt.Scanln(&i)
			if i == "r" {
				timeLogger.Writer.Write(
					[]string{
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

			go runner(quitChan, timeChan)
			fmt.Scanln(&i)

			quitChan <- true

			for {
				resetTime = <-timeChan
				timeToLeave := endTime.Local().
					Add(resetTime.Sub(breakStartTime.Local())).
					Local()

				fmt.Println("You clocked back in at " + resetTime.Local().Format("03:04PM"))
				err = timeLogger.Writer.Write(
					[]string{
						timeLogger.Username,
						startTime.String(),
						timeToLeave.String(),
					},
				)
				if err != nil {
					panic(err)
				}
				timeLogger.Writer.Flush()

				start = false
				endTime = timeToLeave
				break
			}
		} else {
			fmt.Print("You can clock out at " + endTime.Local().Format("03:04PM"))
			start = true
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
