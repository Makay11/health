package health

import (
	"log"
	"os"
)

const loggerFlags = log.Ldate | log.Lmicroseconds | log.Lshortfile

var logger = log.New(os.Stderr, "", loggerFlags)
