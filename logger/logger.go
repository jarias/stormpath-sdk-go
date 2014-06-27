package logger

import (
	"log"
	"os"

	"github.com/onsi/ginkgo"
)

var (
	//ERROR logger
	ERROR *log.Logger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
	//INFO logger
	INFO *log.Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
)

//Init initnitializes the INFO and ERROR loggers
func InitInTestMode() {
	INFO = log.New(ginkgo.GinkgoWriter, "INFO: ", log.Ldate|log.Ltime)
	ERROR = log.New(ginkgo.GinkgoWriter, "ERROR: ", log.Ldate|log.Ltime)
}
