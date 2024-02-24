package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Entry

func InitLogs() *os.File {
	f, err := os.OpenFile("errors.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		logrus.Error(err)
	}
	logrus.SetReportCaller(true)

	logrus.SetFormatter(&logrus.TextFormatter{})

	logrus.SetOutput(f)

	logrus.SetLevel(logrus.DebugLevel)

	Logger = logrus.NewEntry(logrus.StandardLogger())

	return f
}

func LoggerDesc(desc string) *logrus.Entry {
	Logger = Logger.WithField("desc", desc)
	return Logger
}

func GetLoggerDesc() string {
	desc := ""
	if v, ok := Logger.Data["desc"].(string); ok {
		desc = v
	}
	return desc
}
