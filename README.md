Go SDK for the [Stormpath](http://stormpath.com/) API

[![Build Status](https://drone.io/github.com/jarias/stormpath-sdk-go/status.png)](https://drone.io/github.com/jarias/stormpath-sdk-go/latest) [![Coverage Status](https://coveralls.io/repos/jarias/stormpath-sdk-go/badge.png?branch=develop)](https://coveralls.io/r/jarias/stormpath-sdk-go?branch=develop)

#Usage

```go get github.com/jarias/stormpath-sdk-go```

```go
import "github.com/jarias/stormpath-sdk-go"
import "fmt"

//This would look for env variables first STORMPATH_API_KEY_ID and STORMPATH_API_KEY_SECRET if empty
//then it would look for os.Getenv("HOME") + "/.config/stormpath/apiKey.properties" for the credentials
credentials := stormpath.NewDefaultCredentials()

//Whithout cache
stormpath.Client = stormpath.NewStormpathClient(credentials, nil)

//Get the current tenant
tenant := stormpath.CurrentTenant()

//Get the tenat applications
apps := tenant.GetApplications(stormpath.NewDefaultPageRequest(), stormpath.NewEmptyFilter())

//Get the first application
app := apps.Items[0]

//Authenticate a user against the app
accountRef, _ := app.AuthenticateAccount("username", "password")

//Print the account information
account, _ := accountRef.GetAccount()
fmt.Println(account)
```

Features:

* Cache with a sample Redis implementation
* Almost 100% of the Stormpath API implemented
* Load credentials via properties file or env variables
* Requests are authenticated via Stormpath SAuthc1 algorithm

#License

Copyright 2014 Julio Arias

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
