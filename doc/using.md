# Using Components and Bundles

TODO

### Apps (binaries)

Apps represent the final binaries defined by the agent.
Each "flavor" of agent is defined in a different sub-package of `cmd/`.

Apps are formulaic and should not contain any complex logic.
Their job is to parse command-line options, set up an `fx` App, and run it.

For parameterized components, the app must also supply the parameters:

```go
app = fx.New(
    foo.Module,
    fx.Supply(&foo.ModuleParams{
        MaxFoos: 13,
    }),
    ...
)
```

### Component Dependencies

Component dependencies are automatically determined from the arguments to a component constructor.
For example, a component that depends on the log component will have a `logs.Component` in its argument list:

```go
import "github.com/djmitche/dd-agent-comp-experiments/comp/core/log"

func newThing(..., log log.Component, ...) Component {
    return &thing{
        log: log,
        ...
    }
}
```

### Component Bundles

Many components naturally gather into larger areas of the agent codebase, such as DogStatsD.
In many cases, these components are not intended for use outside of that area.
These components should be defined in an `internal/` package, and included in a `Module` definition in the parent package.
This single `Module` can then be included by apps that require that functionality.

TODO: expand, include BundleParams, bundle_test.go, etc.

### Testing

Testing for a component should use `fxtest` to create the component.
This focuses testing on the API surface of the component against which other components will be built.
Per-function unit tests are, of course, also great where appropriate!

Here's an example testing a component with a mocked dependency on `other`:

```go
func TestMyComponent(t *testing.T) {
    var comp Component
    var other other.Component
    app := fxtest.New(t,
        Module,              // use the real version of this component
        other.MockModule,    // use the mock version of other
        fx.Populate(&comp),  // get the instance of this component
        fx.Populate(&other), // get the (mock) instance of the other component
    )

    // start and, at completion of the test, stop the components
    defer app.RequireStart().RequireStop()

    // cast `other` to its mock interface to call mock-specific methods on it
    other.(other.Mock).SetSomeValue(10)                      // Arrange
    comp.DoTheThing()                                        // Act
    require.Equal(t, 20, other.(other.Mock).GetSomeResult()) // Assert
}
```

If the component has a mock implementation, it is a good idea to test that mock implementation as well.
