Go SDK for the [Stormpath](http://stormpath.com/) API

Develop:

[![Build Status](https://travis-ci.org/jarias/stormpath-sdk-go.svg?branch=develop)](https://travis-ci.org/jarias/stormpath-sdk-go) [![codecov.io](http://codecov.io/github/jarias/stormpath-sdk-go/coverage.svg?branch=develop)](http://codecov.io/github/jarias/stormpath-sdk-go?branch=develop)

Master:

[![Build Status](https://travis-ci.org/jarias/stormpath-sdk-go.svg?branch=master)](https://travis-ci.org/jarias/stormpath-sdk-go) [![codecov.io](http://codecov.io/github/jarias/stormpath-sdk-go/coverage.svg?branch=master)](http://codecov.io/github/jarias/stormpath-sdk-go?branch=master)

# Usage

## Core

```go get github.com/jarias/stormpath-sdk-go```

```go
import "github.com/jarias/stormpath-sdk-go"
import "fmt"

//Load the configuration according to the StormPath framework spec 
//See: https://github.com/stormpath/stormpath-sdk-spec/blob/master/specifications/config.md
clientConfig, err := stormpath.LoadConfiguration()

if err != nil {
    stormpath.Logger.Panicf("[ERROR] Couldn't load Stormpath client configuration: %s", err)
}

//Init the client with the loaded config and no specific cache, 
//note that if the cache is enabled via config the default local cache would be used
stormpath.Init(clientConfig, nil)

//Get the current tenant
tenant, _ := stormpath.CurrentTenant()

//Get the tenat applications
apps, _ := tenant.GetApplications(stormpath.MakeApplicationCriteria().NameEq("test app"))

//Get the first application
app := apps.Items[0]

//Authenticate a user against the app
account, _ := app.AuthenticateAccount("username", "password")

fmt.Println(account)
```

## Web

See `web/example/example.go`

Features:

* Cache with a sample local in-memory implementation
* Almost 100% of the Stormpath API implemented
* Load credentials via properties file or env variables
* Load client configuration according to Stormpath framework spec
* Requests are authenticated via Stormpath SAuthc1 algorithm only
* Web extension according to the [Stormpath Spec](https://github.com/stormpath/stormpath-framework-spec)

# Debugging

If you need to trace all requests done to stormpath you can enable debugging in the logs
by setting the environment variable STORMPATH_LOG_LEVEL=DEBUG the default level is ERROR.

# Contributing

Pull request are more than welcome, all pull requests should be from and directed to the ```develop``` branch **NOT** ```master```.

Please make sure you add tests ;)

Development requirements:

- Go 1.6+
- [Testify](https://github.com/stretchr/testify) ```go get github.com/stretchr/testify/assert```
- An [Stormpath](https://stormpath.com) account (for integration testing)

Running the test suite

Env variables:

```
export STORMPATH_API_KEY_ID=XXXX
export STORMPATH_API_KEY_SECRET=XXXX
```

```
go test . -cover -covermode=atomic
```

I'm aiming at 85% test coverage not yet met but thats the goal.

# License

Copyright 2014, 2015, 2016 Julio Arias

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
