package logging

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type Logging struct {
	*logrus.Logger
}

var once sync.Once
var logInstance *Logging

func NewLogger(level logrus.Level) *Logging {
	if logInstance == nil {
		once.Do(func() {
			logger := logrus.New()
			logger.SetLevel(level)
			logInstance = &Logging{logger}
		})
	}
	return logInstance
}

func GetLogger() *Logging {
	return logInstance
}
