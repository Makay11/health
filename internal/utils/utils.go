package utils

import (
	"log"
	"os"
)

const loggerFlags = log.Ldate | log.Lmicroseconds | log.Lshortfile

var Logger = log.New(os.Stderr, "[health] ", loggerFlags)
