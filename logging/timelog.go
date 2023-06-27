package logging

import (
	"encoding/csv"
	"log"
	"os"
	"time"
)

type Logger interface {
	Init()
	Info()
	Error()
	GetWriter()
	GetReader()
}

// TimeLogger is a struct that implements the Logger interface.  
// INFO: You only need to fill in the `User` field.
type TimeLogger struct {
	Username string
	Reader   *csv.Reader
	Writer   *csv.Writer
	File     *os.File
}

// TimeLoggerInfo struct  
type TimeLoggerInfo struct {
	Time     time.Time
	Username string
	Filename string
}

// WARNING: Make sure that you handle closing the file when you are done with it.
func (t *TimeLogger) Init(username string) {
	t.Username = username
	t.Reader = csv.NewReader(t.File)
	t.Writer = csv.NewWriter(t.File)
	t.File = createOrOpenFile("logs/timelog.csv")
}

func (t *TimeLogger) Info() TimeLoggerInfo {
	return TimeLoggerInfo{
		Time:     time.Now(),
		Username: t.Username,
		Filename: t.File.Name(),
	}
}

func (t *TimeLogger) Error(err error) {
	f := createOrOpenFile("logs/errorlog.csv")
	defer f.Close()

	logger := log.New(f, t.Username+": ", log.LstdFlags|log.Ldate|log.Ltime)
	logger.Println("Error " + err.Error())
}

func (t *TimeLogger) GetWriter() *csv.Writer {
	return t.Writer
}

func (t *TimeLogger) GetReader() *csv.Reader {
	return t.Reader
}

func createOrOpenFile(file string) *os.File {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0777))
	if err != nil {
		panic(err.Error())
	}
	return f
}