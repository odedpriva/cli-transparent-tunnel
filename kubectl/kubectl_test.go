package kubectl

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestKubeCtl_getKubeConfig(t *testing.T) {
	k := KubeCtl{
		log: &logrus.Logger{},
	}
	k.getKubeConfig()
}
