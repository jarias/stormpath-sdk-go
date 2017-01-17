#Contributing

Pull request are more than welcome, please follow this sample workflow, make sure you work out of
the develop branch.

- Fork
- Clone `git clone YOUR_USERNAME/stormpath-sdk-go`
- Checkout develop branch `git checkout -t origin/develop`
- Create a feature branch `git checkout -b YOUR_FEATURE_OR_BUG_FIX`
- Create a PR to jarias/develop `hub pull-request -b jarias/stormpath-sdk-go:develop`

Please make sure you add tests ;)

Development requirements:

- Go 1.7+
- [Testify](https://github.com/stretchr/testify) `go get github.com/stretchr/testify/assert`
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

Thanks.
