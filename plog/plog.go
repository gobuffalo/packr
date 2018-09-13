package plog

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger interface for a logger to be used
// with genny. Logrus is 100% compatible.
type Logger interface {
	Debugf(string, ...interface{})
	Debug(...interface{})
	Infof(string, ...interface{})
	Info(...interface{})
	Printf(string, ...interface{})
	Print(...interface{})
	Warnf(string, ...interface{})
	Warn(...interface{})
	Errorf(string, ...interface{})
	Error(...interface{})
	Fatalf(string, ...interface{})
	Fatal(...interface{})
}

var Default = func() Logger {
	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.SetLevel(logrus.DebugLevel)
	// l.SetLevel(logrus.InfoLevel)
	return l
}()

var Debugf = Default.Debugf
var Debug = Default.Debug
var Infof = Default.Infof
var Info = Default.Info
var Printf = Default.Printf
var Print = Default.Print
var Warnf = Default.Warnf
var Warn = Default.Warn
var Errorf = Default.Errorf
var Error = Default.Error
var Fatalf = Default.Fatalf
var Fatal = Default.Fatal
