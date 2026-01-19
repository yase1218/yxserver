package global

import (
	"os"
)

var (
	exitChan   = make(chan struct{})
	signalChan = make(chan os.Signal, 1)
)
