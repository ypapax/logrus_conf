package logrus_conf

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
)

func Files(appName string, level logrus.Level) error {
	type f struct {
		Name string
		File *os.File
	}

	ff := []f{
		{Name: fmt.Sprintf("/tmp/%+v.errors.log", appName)},
		{Name: fmt.Sprintf("/tmp/%+v.info.log", appName)},
	}
	for i, f := range ff {
		os.Remove(f.Name)
		file, err := os.Create(f.Name)
		if err != nil {
			logrus.Error(err)
			return err
		}
		ff[i].File = file
		logrus.Infof("log file %+v", f.Name)
	}
	logrus.SetReportCaller(true)
	logrus.SetOutput(ioutil.Discard)
	errors := io.MultiWriter(os.Stderr, ff[0].File)
	logrus.AddHook(&WriterHook{
		Writer: errors,
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		},
	})

	info := io.MultiWriter(os.Stderr, ff[1].File)
	logrus.AddHook(&WriterHook{
		Writer: info,
		LogLevels: []logrus.Level{
			logrus.InfoLevel,
			logrus.DebugLevel,
			logrus.TraceLevel,
		},
	})
	logrus.SetLevel(level)
	logrus.Tracef("log level: %+v", logrus.GetLevel())
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
