Go SDK for the [Stormpath](http://stormpath.com/) API

Develop:

[![Build Status](https://travis-ci.org/jarias/stormpath-sdk-go.svg?branch=develop)](https://travis-ci.org/jarias/stormpath-sdk-go) [![codecov.io](http://codecov.io/github/jarias/stormpath-sdk-go/coverage.svg?branch=develop)](http://codecov.io/github/jarias/stormpath-sdk-go?branch=develop)

Master:

[![Build Status](https://travis-ci.org/jarias/stormpath-sdk-go.svg?branch=master)](https://travis-ci.org/jarias/stormpath-sdk-go) [![codecov.io](http://codecov.io/github/jarias/stormpath-sdk-go/coverage.svg?branch=master)](http://codecov.io/github/jarias/stormpath-sdk-go?branch=master)

# Usage

```go get github.com/jarias/stormpath-sdk-go```

```go
import "github.com/jarias/stormpath-sdk-go"
import "fmt"

//This would look for env variables first STORMPATH_API_KEY_ID and STORMPATH_API_KEY_SECRET if empty
//then it would look for os.Getenv("HOME") + "/.config/stormpath/apiKey.properties" for the credentials
credentials, _ := stormpath.NewDefaultCredentials()

//Init Whithout cache
stormpath.Init(credentials, nil)

//Get the current tenant
tenant, _ := stormpath.CurrentTenant()

//Get the tenat applications
apps, _ := tenant.GetApplications(MakeApplicationCriteria().NameEq("test app"))

//Get the first application
app := apps.Items[0]

//Authenticate a user against the app
account, _ := app.AuthenticateAccount("username", "password")

fmt.Println(account)
```

Features:

* Cache with a sample Redis implementation
* Almost 100% of the Stormpath API implemented
* Load credentials via properties file or env variables
* Requests are authenticated via Stormpath SAuthc1 algorithm

# Debugging

If you need to trace all requests done to stormpath you can enable debugging in the logs
by setting the environment variable STORMPATH_LOG_LEVEL=DEBUG the default level is ERROR.

# Contributing

Pull request are more than welcome, all pull requests should be from and directed to the ```develop``` branch **NOT** ```master```.

Please make sure you add tests ;)

Development requirements:

- Go 1.4+
- [Ginkgo](https://onsi.github.io/ginkgo/) ```go get github.com/onsi/ginkgo/ginkgo```
- [Gomega](http://onsi.github.io/gomega/) ```go get github.com/onsi/gomega```
- An [Stormpath](https://stormpath.com) account (for integration testing)
- Redis (there is a Docker compose file to easily start up redis)

Running the test suite

Env variables:

```
export STORMPATH_API_KEY_ID=XXXX
export STORMPATH_API_KEY_SECRET=XXXX
export REDIS_SERVER=localhost
```

```
ginkgo -r --randomizeAllSpecs --randomizeSuites --failOnPending --cover
```

I'm aiming at 85% test coverage not yet met but thats the goal.

# License

Copyright 2014, 2015 Julio Arias

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
