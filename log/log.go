package log

import (
	"fmt"
	"log"
)

const (
	dbg = 0x01
	inf = 0x02
	wrn = 0x04
	err = 0x08

	LogLevelDebug = dbg | inf | wrn | err
	LogLevelInfo = inf | wrn | err
	LogLevelWarn = wrn | err
	LogLevelError = err
)

var logLevel uint32 = LogLevelInfo

func out(logType string, message string) {
	log.Printf("%s,\"%s\"\n", logType, message)
}

func SetLogLevel(level uint32) {
	logLevel = level
}

func D(format string, args ...interface{}) {
	if (logLevel & dbg) != 0 {
		out("DBG", fmt.Sprintf(format, args...))
	}
}

func De(error error) {
	if (logLevel & dbg) != 0 {
		out("DBG", fmt.Sprintf("%+v", error))
	}
}

func I(format string, args ...interface{}) {
	if (logLevel & inf) != 0 {
		out("INF", fmt.Sprintf(format, args...))
	}
}

func Ie(error error) {
	if (logLevel & inf) != 0 {
		out("INF", fmt.Sprintf("%+v", error))
	}
}

func W(format string, args ...interface{}) {
	if (logLevel & wrn) != 0 {
		out("WRN", fmt.Sprintf(format, args...))
	}
}

func We(error error) {
	if (logLevel & wrn) != 0 {
		out("WRN", fmt.Sprintf("%+v", error))
	}
}

func E(format string, args ...interface{}) {
	if (logLevel & err) != 0 {
		out("ERR", fmt.Sprintf(format, args...))
	}
}

func Ee(error error) {
	if (logLevel & err) != 0 {
		out("WRN", fmt.Sprintf("%+v", error))
	}
}
