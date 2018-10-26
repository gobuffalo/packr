package plog

import (
	"fmt"

	"github.com/gobuffalo/logger"
	"github.com/sirupsen/logrus"
)

var Default = logger.New(logger.ErrorLevel)

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
