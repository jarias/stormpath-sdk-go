#Contributing

Pull request are more than welcome, all pull requests should be from and directed to the ```develop``` branch **NOT** ```master```.

Please make sure you add tests ;)

Development requirements:

- Go 1.4+
- [Ginkgo](https://onsi.github.io/ginkgo/) ```go get github.com/onsi/ginkgo/ginkgo```
- [Gomega](http://onsi.github.io/gomega/) ```go get github.com/onsi/gomega```
- An [Stormpath](https://stormpath.com) account (for integration testing)

Running the test suite

Env variables:

export STORMPATH_API_KEY_ID=XXXX
export STORMPATH_API_KEY_SECRET=XXXX

```
ginkgo -r --randomizeAllSpecs --randomizeSuites --failOnPending --cover
```

I'm aiming at 85% test coverage not yet met but thats the goal.

Thanks.
