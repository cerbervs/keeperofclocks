package main

import (
	"fmt"
	"os"
	"time"

	"github.com/cerbervs/keeperofclocks/logging"
)

func main() {
	timeLogger := logging.TimeLogger{}
	timeLogger.Init(os.Getenv("USER"))
	defer timeLogger.File.Close()

	print(
		timeLogger.Info().Time.String() + "\n" + timeLogger.Info().Username + "\n" + timeLogger.Info().Filename + "\n",
	)

	var i string
	fmt.Scanln(&i)

	r, err := timeLogger.Reader.ReadAll()
	if err != nil {
		timeLogger.Error(err)
	}
	if r != nil {
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

	start := true
	for {
		if start {
			quitChan := make(chan bool, 1)
			timeChan := make(chan time.Time, 1)

			// start timer
			go runner(quitChan, timeChan)

			// record local start time for later calcs
			startTime := time.Now().Local()

			err := timeLogger.Writer.Write(
				[]string{
					timeLogger.Username,
					time.Now().Local().String(),
					startTime.Local().Add(time.Hour * 8).String(),
				},
			)
			if err != nil {
				timeLogger.Error(err)
			}

			fmt.Scanln()
			quitChan <- true
			endTime := <-timeChan
			err = timeLogger.Writer.Write(
				[]string{
					timeLogger.Username,
					time.Now().Local().String(),
					// INFO: Gets the time elapsed from the local start time and adds it to the current time.
					startTime.Add(endTime.Local().Sub(startTime.Local())).Local().String(),
				},
			)
			if err != nil {
				timeLogger.Error(err)
			}

			start = false
		} else {
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
		default:
			// run continuously
		}
	}
}