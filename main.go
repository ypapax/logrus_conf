package logrus_conf

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type FilesConf struct {
	FileNameSuffix string
	LogLevels      []logrus.Level
	file           *os.File
}

func AllLevelFiles(dir, appName string, level logrus.Level) error {
	ff := []FilesConf{
		{
			FileNameSuffix: "panic.log",
			LogLevels: []logrus.Level{
				logrus.PanicLevel,
			},
		},
		{
			FileNameSuffix: "fatal.log",
			LogLevels: []logrus.Level{
				logrus.FatalLevel,
			},
		},
		{
			FileNameSuffix: "error.log",
			LogLevels: []logrus.Level{
				logrus.ErrorLevel,
			},
		},
		{
			FileNameSuffix: "warn.log",
			LogLevels: []logrus.Level{
				logrus.WarnLevel,
			},
		},
		{
			FileNameSuffix: "info.log",
			LogLevels: []logrus.Level{
				logrus.InfoLevel,
			},
		},
		{
			FileNameSuffix: "debug.log",
			LogLevels: []logrus.Level{
				logrus.DebugLevel,
			},
		},
		{
			FileNameSuffix: "trace.log",
			LogLevels: []logrus.Level{
				logrus.TraceLevel,
			},
		},
	}
	return Files(dir, appName, level, ff)
}

func Files(dir, appName string, level logrus.Level, ff []FilesConf) error {
	for i, f := range ff {
		os.Remove(f.FileNameSuffix)
		fullFileName := filepath.Join(dir, fmt.Sprintf("%+v.%+v", appName, f.FileNameSuffix))
		file, err := os.Create(fullFileName)
		if err != nil {
			err := errors.WithStack(err)
			return err
		}
		ff[i].file = file
		logrus.Infof("log file %+v for levels: %+v", f.FileNameSuffix, f.LogLevels)
	}

	logrus.SetReportCaller(true)
	logrus.SetOutput(ioutil.Discard)
	for _, f := range ff {
		mr := io.MultiWriter(os.Stderr, f.file)
		logrus.AddHook(&WriterHook{
			Writer:    mr,
			LogLevels: f.LogLevels,
		})
	}
	logrus.SetLevel(level)
	return nil
}

// WriterHook is a hook that writes logs of specified LogLevels to specified Writer
type WriterHook struct {
	Writer    io.Writer
	LogLevels []logrus.Level
}

// Fire will be called when some logging function is called with current hook
// It will format log entry to string and write it to appropriate writer
func (hook *WriterHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write([]byte(line))
	return err
}

// Levels define on which log levels this hook would trigger
func (hook *WriterHook) Levels() []logrus.Level {
	return hook.LogLevels
}
