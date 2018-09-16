package plog

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var Default = func() *logrus.Logger {
	l := logrus.New()
	l.SetOutput(os.Stdout)
	// l.SetLevel(logrus.DebugLevel)
	l.SetLevel(logrus.InfoLevel)
	// l.Formatter = &logrus.JSONFormatter{}
	return l
}()

func Debug(t interface{}, m string, args ...interface{}) {
	if len(args)%2 == 1 {
		args = append(args, "")
	}
	f := logrus.Fields{}
	for i := 0; i < len(args); i += 2 {
		k := args[i]
		v := args[i+1]
		f[fmt.Sprint(k)] = v
	}
	e := Default.WithFields(f)
	e.Debugf("%T#%s", t, m)
}
