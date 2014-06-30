package logger

import (
	"log"
	"os"

	"github.com/onsi/ginkgo"
)

var (
	//ERROR logger
	ERROR = log.New(os.Stderr, "ERROR stormpath-sdk-go: ", log.Ldate|log.Ltime)
	//INFO logger
	INFO = log.New(os.Stdout, "INFO stormpath-sdk-go: ", log.Ldate|log.Ltime)
	//CACHE logger
	CACHE = log.New(os.Stdout, "CACHE stormpath-sdk-go: ", log.Ldate|log.Ltime)
)

//Init initnitializes the INFO and ERROR loggers
func InitInTestMode() {
	INFO = log.New(ginkgo.GinkgoWriter, "INFO stormpath-sdk-go: ", log.Ldate|log.Ltime)
	ERROR = log.New(ginkgo.GinkgoWriter, "ERROR stormpath-sdk-go: ", log.Ldate|log.Ltime)
	CACHE = log.New(ginkgo.GinkgoWriter, "CACHE stormpath-sdk-go: ", log.Ldate|log.Ltime)
}
