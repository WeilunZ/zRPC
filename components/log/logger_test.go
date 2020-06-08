package log

import (
	"log"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	myLogger := &logger{
		Logger: log.New(os.Stdout, "", log.Lshortfile|log.LstdFlags),
		options: &Options{
			Path:         "../rpc.log",
			FrameLogPath: "../frame.log",
			Level:        DEBUG,
		},
	}
	myLogger.Infof("test info, msg : %s", "info")
	myLogger.Warningf("test warning msg: %s", "warning")
	myLogger.Errorf("test error msg: %s", "error")
}
