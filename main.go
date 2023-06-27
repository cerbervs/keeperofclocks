package main

import (
	"errors"
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

	print(
		timeLogger.Info().Time.String() + "\n" + timeLogger.Info().Username + "\n" + timeLogger.Info().Filename + "\n",
	)

	r, err := timeLogger.Reader.Read()
	if err != nil {
		timeLogger.Error(errors.New("read " + err.Error()))
	}
	if len(r) == 0 {
		timeLogger.Writer.Write(
			[]string{"username", "time", "time_to_leave"},
		)
	}

	timeLogger.Writer.Write(
		[]string{
			timeLogger.Username,
			timeLogger.Info().Time.String(),
			time.Now().Local().Add(time.Hour * 8).String(),
		},
	)

	startTime := time.Now().Local()
	resetTime := startTime
	endTime := startTime.Add(time.Hour * 8)
	start := true
	for {
		fmt.Println("Press enter to start timer.")
		fmt.Scanln()

		if start {
			quitChan := make(chan bool)
			timeChan := make(chan time.Time)

			// start timer
			go runner(quitChan, timeChan)

			// record local start time for later calcs
			err := timeLogger.Writer.Write(
				[]string{
					timeLogger.Username,
					time.Now().Local().String(),
					startTime.Local().Add(time.Hour * 8).String(),
				},
			)
			if err != nil {
				timeLogger.Error(errors.New("start true? " + err.Error()))
			}

			timeLogger.Writer.Flush()

			fmt.Scanln()
			quitChan <- true
			for {
				select {
				case <-timeChan:
					err = timeLogger.Writer.Write(
						[]string{
							timeLogger.Username,
							startTime.Local().String(),
							// INFO: Adds the time elapsed from the start time to the start time.
							endTime.Local().
								Add(time.Now().Local().Sub(resetTime.Local())).
								Local().
								String(),
						},
					)
					if err != nil {
						timeLogger.Error(errors.New("after endtime" + err.Error()))
					}
					timeLogger.Writer.Flush()

					resetTime = time.Now().Local()
				default:
					continue
				}
				break
			}
			start = false
		} else {
			fmt.Println("You can clock out at " + endTime.Local().Format("00:00:00"))
			start = true
			continue
		}
	}
}

func runner(quitChan chan bool, timeChan chan time.Time) time.Time {
	for {
		select {
		case <-quitChan:
			timeChan <- time.Now()
		}
	}
}