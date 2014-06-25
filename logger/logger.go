package logger

import (
	"io/ioutil"
	"log"
	"os"
)

var (
	//ERROR logger
	ERROR *log.Logger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
	//INFO logger
	INFO *log.Logger = log.New(ioutil.Discard, "INFO: ", log.Ldate|log.Ltime)
)

//Init initnitializes the INFO and ERROR loggers
func Init(verbose bool) {
	if verbose {
		INFO = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	}
}
