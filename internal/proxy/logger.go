package proxy

import (
	"fmt"
	"github.com/apex/log"
	"strings"
)

type LogWrapper struct{}

func (w LogWrapper) Printf(format string, v ...interface{}) {
	s := strings.TrimSuffix(fmt.Sprintf(format, v...), "\n")
	log.Errorf(s)
}
