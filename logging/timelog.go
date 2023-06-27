package logging

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
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
	t.File = createOrOpenFile("timelog.csv")
}

func (t *TimeLogger) Info() TimeLoggerInfo {
	return TimeLoggerInfo{
		Time:     time.Now(),
		Username: t.Username,
		Filename: t.File.Name(),
	}
}

func (t *TimeLogger) Error(err error) {
	f := createOrOpenFile("errorlog.csv")
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
	errLog := log.New(os.Stderr, "ERROR: ", log.LstdFlags|log.Ldate|log.Ltime)

	path, err := filepath.Abs("./")
	outPath := filepath.Join(path, "logs")

	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		os.MkdirAll(outPath, os.FileMode(0777))
	}
	if err != nil {
		errLog.Fatal(err)
	}

	f, err := os.OpenFile(
		filepath.Join(outPath, file),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		os.FileMode(0777),
	)
	if err != nil {
		errLog.Fatal(err)
	}

	return f
}
