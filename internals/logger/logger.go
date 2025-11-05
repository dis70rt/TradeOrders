package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func Init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
		TimestampFormat:        "15:04:05.000",    
		PadLevelText:           true,              
		DisableLevelTruncation: true,
	})
}

func Info(msg string) {
	log.Info(msg)
}

func WithError(err error) *logrus.Entry {
	return log.WithError(err)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}
