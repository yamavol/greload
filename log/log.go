package log

import (
	"fmt"
	"strings"
)

var logLevel int = LevelInfo
var maskPatterns = []string{}

const (
	LevelDebug int = -4
	LevelInfo  int = 0
	LevelWarn  int = 4
	LevelError int = 8
)

func SetLogLevel(level int) {
	logLevel = level
}

func Debug(args ...interface{}) {
	if logLevel <= LevelDebug {
		fmt.Println(args...)
	}
}

func Debugf(format string, args ...interface{}) {
	if logLevel <= LevelDebug {
		fmt.Println(fmt.Sprintf(format, args...))
	}
}

func Info(args ...interface{}) {
	if logLevel <= LevelInfo {
		fmt.Println(args...)
	}
}

func Infof(format string, args ...interface{}) {
	if logLevel <= LevelInfo {
		fmt.Println(fmt.Sprintf(format, args...))
	}
}

func Error(args ...interface{}) {
	if logLevel <= LevelError {
		fmt.Println(args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if logLevel <= LevelError {
		fmt.Println(fmt.Sprintf(format, args...))
	}
}

func Warn(args ...interface{}) {
	if logLevel <= LevelWarn {
		fmt.Println(args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if logLevel <= LevelWarn {
		fmt.Println(fmt.Sprintf(format, args...))
	}
}

func AddMask(pattern string) {
	maskPatterns = append(maskPatterns, pattern)
}

func DebugSecuref(format string, args ...interface{}) {
	if logLevel <= LevelDebug {
		s := fmt.Sprintf(format, args...)
		for _, mask := range maskPatterns {
			s = strings.Replace(s, mask, "***", -1)
		}
		Debugf(s)
	}
}
