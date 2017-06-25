/*
  trader API Engine
*/

package logger

import (
	. "config"
	"github.com/sirupsen/logrus"
	"github.com/rifflock/lfshook"
	// "github.com/zbindenren/logrus_mail"
	"os"
	"path/filepath"
	"runtime"
)

var (
	// 日志文件
	debug_file = Root + "/log/debug.log"
	info_file  = Root + "/log/info.log"
	error_file = Root + "/log/error.log"
	panic_file = Root + "/log/panic.log"
)

func init() {
	os.Mkdir(Root+"/log/", 0777)
}

var Log *logrus.Logger

func NewLogger() *logrus.Logger {
	if Log != nil {
		return Log
	}

	Log = logrus.New()
	textFormatter := new(logrus.TextFormatter)
	// textFormatter.DisableColors = true
	textFormatter.FullTimestamp = true
	textFormatter.TimestampFormat = "20060102 15:04:05"
	Log.Formatter = textFormatter
	if Config["loglevel"] == "debug" {
		Log.Level = logrus.DebugLevel
	} else {
		Log.Level = logrus.InfoLevel
	}

	Log.Hooks.Add(lfshook.NewHook(lfshook.PathMap{
		logrus.DebugLevel: debug_file,
		logrus.InfoLevel:  info_file,
		logrus.ErrorLevel: error_file,
		logrus.PanicLevel: panic_file,
	}))

	// if you do not need authentication for your smtp host
	// hook, err := logrus_mail.NewMailAuthHook("trader", "smtp.ym.163.com", 994, "haobtc@blocktip.com", "78623269@qq.com", "haobtc@blocktip.com", "pfE8pmQUUK00")
	// if err == nil {
	// 	Log.Hooks.Add(hook)
	// }

	return Log
}

func Debugln(args ...interface{}) {
	// 加上文件调用和行号
	_, callerFile, line, ok := runtime.Caller(1)
	if ok {
		args = append([]interface{}{"[", filepath.Base(callerFile), "]", line}, args...)
	}

	NewLogger().Debugln(args...)
}

func Infoln(args ...interface{}) {
	// 加上文件调用和行号
	_, callerFile, line, ok := runtime.Caller(1)
	if ok {
		args = append([]interface{}{"[", filepath.Base(callerFile), "]", line}, args...)
	}
	NewLogger().Infoln(args...)
}

func Errorln(args ...interface{}) {
	// 加上文件调用和行号
	_, callerFile, line, ok := runtime.Caller(1)
	if ok {
		args = append([]interface{}{"[", filepath.Base(callerFile), "]", line}, args...)
	}
	NewLogger().Errorln(args...)
}

func Panicln(args ...interface{}) {
	// 加上文件调用和行号
	_, callerFile, line, ok := runtime.Caller(1)
	if ok {
		args = append([]interface{}{"[", filepath.Base(callerFile), "]", line}, args...)
	}
	NewLogger().Panicln(args...)
}
