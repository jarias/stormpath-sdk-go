package stormpath_test

import (
	"log"
	"os"
	"runtime"

	"github.com/garyburd/redigo/redis"
	"github.com/hashicorp/logutils"
	. "github.com/jarias/stormpath-sdk-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	uuid "github.com/nu7hatch/gouuid"

	"testing"
)

var (
	app     *Application
	cred    Credentials
	account *Account
	tenant  *Tenant
)

func TestStormpath(t *testing.T) {
	runtime.GOMAXPROCS(4)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Stormpath Suite")
}

func randomName() string {
	uuid, _ := uuid.NewV4()
	return uuid.String()
}

func newTestApplication() *Application {
	return NewApplication("app-" + randomName())
}

func newTestGroup() *Group {
	return NewGroup("group-" + randomName())
}

func newTestDirectory() *Directory {
	return NewDirectory("directory-" + randomName())
}

func newTestAccount() *Account {
	name := randomName()
	return NewAccount(name+"@test.org", "1234567z!A89", name, name)
}

func initLogInTestMode() {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: "DEBUG",
		Writer:   GinkgoWriter,
	}

	log.SetOutput(filter)
}

var _ = BeforeSuite(func() {
	var err error
	cred, err = NewDefaultCredentials()
	if err != nil {
		panic(err)
	}

	stormpathBaseURL := os.Getenv("STORMPATH_BASE_URL")
	if stormpathBaseURL != "" {
		BaseURL = stormpathBaseURL
	}

	cacheEnabled := os.Getenv("CACHE_ENABLED")
	if cacheEnabled == "true" {
		redisServer := os.Getenv("REDIS_SERVER")
		redisConn, err := redis.Dial("tcp", redisServer+":6379")
		if err != nil {
			panic(err)
		}

		Init(cred, RedisCache{redisConn})
	} else {
		Init(cred, nil)
	}
	initLogInTestMode()

	tenant, err = CurrentTenant()
	if err != nil {
		panic(err)
	}

	app = newTestApplication()

	err = tenant.CreateApplication(app)
	if err != nil {
		panic(err)
	}

	account = NewAccount("test@test.org", "1234567z!A89", "test", "test")
	app.RegisterAccount(account)
})

var _ = AfterSuite(func() {
	if app != nil {
		app.Purge()
	}
})
