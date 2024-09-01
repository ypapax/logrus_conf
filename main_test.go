package logrus_conf_test

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/ypapax/logrus_conf"
)

func TestAllLevelFiles(t *testing.T) {
	as := assert.New(t)
	if err := logrus_conf.AllLevelFiles("/tmp", "logrus_conf_test", logrus.TraceLevel); as.NoError(err) {
		return
	}
}

func TestPrepareFromEnv(t *testing.T) {
	logrus_conf.PrepareFromEnv("test")
	logrus.Infof("hello")
}
