package util

import (
	"github.com/Sirupsen/logrus"
	"log"
	"os"
)

func CheckErr(err error, msg ...string) {
	if err != nil {
		log.Fatal(err, msg)
	}
}

// Log some common log for tests
var Log = &logrus.Logger{
	Level:     logrus.DebugLevel,
	Out:       os.Stdout,
	Formatter: newJSONFormatter(),
}

func newJSONFormatter() *logrus.JSONFormatter {
	return &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.999",
	}
}
