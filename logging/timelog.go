package logging

import (
	"encoding/csv"
	"errors"
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
// INFO: You only need to fill in the `User` field and provide a filename
type TimeLogger struct {
	File      *os.File
	ErrorFile *os.File
	Username  string
	Reader    *csv.Reader
	Writer    *csv.Writer
}

// TimeLoggerInfo struct  
type TimeLoggerInfo struct {
	Time     time.Time
	Username string
	Filename string
}

// WARNING: Make sure that you handle closing the file when you are done with it.
func (t *TimeLogger) Init(username string, outfile string, errfile string) {
	of, err := createOrOpenFile(outfile)
	ef, oferr := createOrOpenFile(errfile)
	if err != nil {
		t.Error(errors.New("init " + err.Error()))
	}
	if oferr != nil {
		t.Error(errors.New("init " + oferr.Error()))
	}

	t.File = of
	t.ErrorFile = ef
	t.Username = username
	t.Reader = csv.NewReader(t.File)
	t.Writer = csv.NewWriter(t.File)
}

func (t *TimeLogger) Info() TimeLoggerInfo {
	return TimeLoggerInfo{
		Time:     time.Now(),
		Username: t.Username,
		Filename: t.File.Name(),
	}
}

func (t *TimeLogger) Error(err error) {
	if err != nil {
		return
	}
	logger := log.New(t.ErrorFile, t.Username+": ", log.LstdFlags|log.Ldate|log.Ltime)
	logger.Println("Error " + err.Error())
}

func (t *TimeLogger) GetWriter() *csv.Writer {
	return t.Writer
}

func (t *TimeLogger) GetReader() *csv.Reader {
	return t.Reader
}

func createOrOpenFile(file string) (*os.File, error) {
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
		errLog.Println(err)
		return nil, err
	}

	return f, nil
}