package log

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

const (
	TRACE = iota
	DEBUG
	INFO
	WARNING
	ERROR
	FATAL
)

const DefaultLogPath = "../zRPC.log"

var (
	DefaultLog *logger
)

type Log interface {
	Trace(v ...interface{})
	Debug(v ...interface{})
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	Fatal(v ...interface{})
	Tracef(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warningf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
}

type logger struct {
	*log.Logger
	options *Options
}

type level int

type Options struct {
	Path         string `default:"../log/zRPC.log"`
	FrameLogPath string `default:"../log/frame.log"`
	Level        level  `default:"debug"`
}

func init() {
	file, err := os.OpenFile(DefaultLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file: ", err)
	}
	DefaultLog = &logger{
		Logger: log.New(file, "", log.LstdFlags|log.Lshortfile),
		options: &Options{
			Level: DEBUG,
		},
	}
}

func (level level) ToString() string {
	switch level {
	case TRACE:
		return "trace"
	case DEBUG:
		return "debug"
	case INFO:
		return "info"
	case WARNING:
		return "warn"
	case ERROR:
		return "error"
	case FATAL:
		return "fatal"
	default:
		return "unknown"
	}
}

func (o *Options) WithPath(path string) func(*Options) {
	return func(options *Options) {
		options.Path = path
	}
}

func (o *Options) WithFrameLogPath(frameLogPath string) func(*Options) {
	return func(options *Options) {
		options.FrameLogPath = frameLogPath
	}
}

func (o *Options) WithLevel(level level) func(*Options) {
	return func(options *Options) {
		options.Level = level
	}
}

// Trace print trace log
func Trace(v ...interface{}) {
	DefaultLog.Trace(v...)
}

// Tracef print a formatted trace log
func Tracef(format string, v ...interface{}) {
	DefaultLog.Tracef(format, v...)
}

func (log *logger) Trace(v ...interface{}) {
	if log.options.Level > TRACE {
		return
	}
	data := log.Prefix() + fmt.Sprint(v...)
	Output(log, 4, "[TRACE] ", data)
}

func (log *logger) Tracef(format string, v ...interface{}) {
	if log.options.Level > TRACE {
		return
	}
	data := log.Prefix() + fmt.Sprintf(format, v...)
	Output(log, 4, "[TRACE] ", data)
}

// Debug print debug log
func Debug(v ...interface{}) {
	DefaultLog.Debug(v...)
}

// Debugf print a formatted debug log
func Debugf(format string, v ...interface{}) {
	DefaultLog.Debugf(format, v...)
}

func (log *logger) Debug(v ...interface{}) {
	if log.options.Level > DEBUG {
		return
	}
	data := log.Prefix() + fmt.Sprint(v...)
	Output(log, 4, "[DEBUG] ", data)
}

func (log *logger) Debugf(format string, v ...interface{}) {
	if log.options.Level > DEBUG {
		return
	}
	data := log.Prefix() + fmt.Sprintf(format, v...)
	Output(log, 4, "[DEBUG] ", data)
}

// Info print info log
func Info(v ...interface{}) {
	DefaultLog.Info(v...)
}

// Infof print a formatted info log
func Infof(format string, v ...interface{}) {
	DefaultLog.Infof(format, v...)
}

func (log *logger) Info(v ...interface{}) {
	if log.options.Level > INFO {
		return
	}
	data := log.Prefix() + fmt.Sprint(v...)
	Output(log, 4, "[INFO] ", data)
}

func (log *logger) Infof(format string, v ...interface{}) {
	if log.options.Level > INFO {
		return
	}
	data := log.Prefix() + fmt.Sprintf(format, v...)
	Output(log, 4, "[INFO] ", data)
}

// Warning print warning log
func Warning(v ...interface{}) {
	DefaultLog.Warning(v...)
}

// Warningf print a formatted warning log
func Warningf(format string, v ...interface{}) {
	DefaultLog.Warningf(format, v...)
}

func (log *logger) Warning(v ...interface{}) {
	if log.options.Level > WARNING {
		return
	}
	data := log.Prefix() + fmt.Sprint(v...)
	Output(log, 4, "[WARNING] ", data)
}

func (log *logger) Warningf(format string, v ...interface{}) {
	if log.options.Level > WARNING {
		return
	}
	data := log.Prefix() + fmt.Sprintf(format, v...)
	Output(log, 4, "[WARNING] ", data)
}

// Error print error log
func Error(v ...interface{}) {
	DefaultLog.Error(v...)
}

// Errorf print a formatted error log
func Errorf(format string, v ...interface{}) {
	DefaultLog.Errorf(format, v...)
}

func (log *logger) Error(v ...interface{}) {
	if log.options.Level > ERROR {
		return
	}
	data := log.Prefix() + fmt.Sprint(v...)
	Output(log, 4, "[ERROR] ", data)
}

func (log *logger) Errorf(format string, v ...interface{}) {
	if log.options.Level > ERROR {
		return
	}
	data := log.Prefix() + fmt.Sprintf(format, v...)
	Output(log, 4, "[ERROR] ", data)
}

// Fatal print fatal log
func Fatal(v ...interface{}) {
	DefaultLog.Fatal(v...)
}

// Fatalf print a formatted fatal log
func Fatalf(format string, v ...interface{}) {
	DefaultLog.Fatalf(format, v...)
}

func (log *logger) Fatal(v ...interface{}) {
	if log.options.Level > FATAL {
		return
	}
	data := log.Prefix() + fmt.Sprint(v...)
	Output(log, 4, "[FATAL] ", data)
}

func (log *logger) Fatalf(format string, v ...interface{}) {
	if log.options.Level > FATAL {
		return
	}
	data := log.Prefix() + fmt.Sprintf(format, v...)
	Output(log, 4, "[FATAL] ", data)
}

// call Output to write log
func Output(log *logger, calldepth int, prefix string, data string) {
	var buffer bytes.Buffer
	buffer.WriteString(prefix)
	buffer.WriteString(data)
	log.Output(calldepth, buffer.String())
}
