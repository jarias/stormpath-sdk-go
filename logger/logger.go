package logger

import (
	"log"
	"os"

	"github.com/onsi/ginkgo"
)

var (
	//ERROR logger
	ERROR = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
	//INFO logger
	INFO = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	//CACHE logger
	CACHE = log.New(os.Stdout, "CACHE: ", log.Ldate|log.Ltime)
)

//Init initnitializes the INFO and ERROR loggers
func InitInTestMode() {
	INFO = log.New(ginkgo.GinkgoWriter, "INFO: ", log.Ldate|log.Ltime)
	ERROR = log.New(ginkgo.GinkgoWriter, "ERROR: ", log.Ldate|log.Ltime)
	CACHE = log.New(ginkgo.GinkgoWriter, "CACHE: ", log.Ldate|log.Ltime)
}
