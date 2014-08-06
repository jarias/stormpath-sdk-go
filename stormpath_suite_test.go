package stormpath_test

import (
	"os"

	"github.com/garyburd/redigo/redis"
	. "github.com/jarias/stormpath-sdk-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var (
	app     *Application
	cred    *Credentials
	account *Account
	tenant  *Tenant
)

func TestStormpath(t *testing.T) {
	InitInTestMode()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Stormpath Suite")
}

var _ = BeforeSuite(func() {
	var err error
	cred, err = NewDefaultCredentials()
	if err != nil {
		panic(err)
	}
	redisServer := os.Getenv("REDIS_SERVER")
	redisConn, err := redis.Dial("tcp", redisServer+":6379")
	if err != nil {
		panic(err)
	}
	Client = NewStormpathClient(cred, RedisCache{redisConn})

	tenant, err = CurrentTenant()
	if err != nil {
		panic(err)
	}

	app = NewApplication("test-app")

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
