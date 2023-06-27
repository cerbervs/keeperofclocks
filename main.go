package main

import "github.com/cerbervs/keeperofclocks/logging"

func main() {
	timeLogger := logging.TimeLogger{}
	timeLogger.Init("test")
	defer timeLogger.File.Close()

	print(
		timeLogger.Info().Time.String() + "\n" + timeLogger.Info().Username + "\n" + timeLogger.Info().Filename + "\n",
	)
}