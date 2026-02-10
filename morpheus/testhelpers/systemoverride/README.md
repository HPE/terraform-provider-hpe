# System Override

This package allows redirection of different test cases to different systems. The goal is to make it easy to control on which system a test will be executed.

To use the package in a test add an invocation to `systemoverride.ParseFlags()` into the `TestMain()` function of the package you want to develop tests for.
Then in each test function call `systemoverride.GetPreferred()` which will return either the default specified system or an overrided value passed in through CLI flags.
To get different parameters per system use the `systemoverride.GetParameters()` function, this function takes in the name of the test system as well as a slice of `systemoverride.SystemTestParameters`, each `SystemTestParameters` struct contains the name of the system for which the parameters are valid.


## CLI Flag overrides

To change which system a test is being executed on without making changes to the code it is possible to pass custom overrides through CLI flags when executing tests.

The flag for setting custom overrides is `-system-override`.

### Example usage

To force `TestAccMorpheusTaskExampleConditionalOk` to run on the `feature` system:
```bash
go test ./... -v -system-override TestAccMorpheusTaskExampleConditionalOk=feature
```

To force all tests to run on the `zodiac` system:
```bash
go test ./... -v -system-override all=feature
```

It is also possible to provide multiple custom overrides by specifying the `-system-override` flag multiple times.