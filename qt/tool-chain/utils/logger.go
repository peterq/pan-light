package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func init() {
	Log.Out = os.Stderr
	Log.Level = logrus.InfoLevel
}

func Debug() {
	Log.Level = logrus.DebugLevel
}
