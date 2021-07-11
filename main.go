package logrus_conf

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const logDirEnvVarName = "LOG_DIR"

const (
	logLevelEnvName = "LOG_LEVEL"
	defaultLogLevel = logrus.InfoLevel
)

func ParseLogLevelFromEnv() (logrus.Level, error) {
	ll := os.Getenv(logLevelEnvName)
	if len(ll) == 0 {
		logrus.Infof("no %+v env, so set default: %+v", logLevelEnvName, defaultLogLevel)
		return defaultLogLevel, nil
	}
	lv, err := logrus.ParseLevel(ll)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return lv, nil
}

func GetLogDir() string {
	logPath := os.Getenv(logDirEnvVarName)
	if len(logPath) == 0 {
		logrus.Infof("missing %+v env var, so skipping writing to files", logDirEnvVarName)
		return ""
	}
	return logPath
}


func LogToFilesTraceForDev(appName string) error {
	logPath := GetLogDir()
	logrus.Tracef("logPath: %+v", logPath)
	if len(logPath) == 0 {
		return nil
	}
	lvl, err := ParseLogLevelFromEnv()
	if err != nil {
		return errors.Wrapf(err, "couldn't get log level from env")
	}
	if err := AllLevelFiles(logPath, appName, lvl); err != nil {
		return errors.WithStack(err)
	}
	return nil
}


func PrepareFromEnv(appName string) error {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	lv, err := ParseLogLevelFromEnv()
	if err != nil {
		return errors.WithStack(err)
	}
	log.Printf("log level from env: %+v\n", lv)
	Prepare(uint32(lv))
	if err := LogToFilesTraceForDev(appName); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func Prepare(logLevel uint32) {
	customFormatter := logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	}
	customFormatter.TimestampFormat = "2006-01-02 15:04:05.999999999 -0700"
	logrus.SetFormatter(&customFormatter)
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.Level(logLevel))
}

type FilesConf struct {
	Name      string
	LogLevels []logrus.Level
	file      *os.File
}

const ext = ".log"

func AllLevelFiles(dir, appName string, level logrus.Level) error {
	ff := []FilesConf{
		{
			Name: "panic",
			LogLevels: []logrus.Level{
				logrus.PanicLevel,
			},
		},
		{
			Name: "fatal",
			LogLevels: []logrus.Level{
				logrus.FatalLevel,
			},
		},
		{
			Name: "error",
			LogLevels: []logrus.Level{
				logrus.ErrorLevel,
			},
		},
		{
			Name: "warn",
			LogLevels: []logrus.Level{
				logrus.WarnLevel,
			},
		},
		{
			Name: "info",
			LogLevels: []logrus.Level{
				logrus.InfoLevel,
			},
		},
		{
			Name: "debug",
			LogLevels: []logrus.Level{
				logrus.DebugLevel,
			},
		},
		{
			Name: "trace",
			LogLevels: []logrus.Level{
				logrus.TraceLevel,
			},
		},
	}
	return Files(dir, appName, level, ff)
}

func Files(dir, appName string, level logrus.Level, ff []FilesConf) error {
	for i, f := range ff {
		os.Remove(f.Name)
		fullFileName := filepath.Join(dir, fmt.Sprintf("%+v.%+v%s", f.Name, appName, ext))
		file, err := os.Create(fullFileName)
		if err != nil {
			err := errors.WithStack(err)
			return err
		}
		ff[i].file = file
		logrus.Infof("log file %+v for levels: %+v", fullFileName, f.LogLevels)
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
