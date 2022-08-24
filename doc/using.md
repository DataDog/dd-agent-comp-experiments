# Using Components and Bundles

## Component Dependencies

Component dependencies are automatically determined from the arguments to a component constructor.
Most components have a few dependencies, and use a struct named `dependencies` to represent them:

```go
type dependencies struct {
    fx.In

    Lc fx.Lifecycle
    Params internal.BundleParams
    Config config.Module
    Log log.Module
    // ...
}

func newThing(deps dependencies) Component {
    t := &thing{
        log: deps.Log,
        ...
    }
    deps.Lc.Append(fx.Hook{OnStart: t.start})
    return t
}
```

## Testing

Tests for a component should use `pkg/util/comptest`, which is a thin wrapper around [`fxtest`](https://pkg.go.dev/go.uber.org/fx/fxtest), to create apps in which to test the component.
This approach focuses testing on the API surface of the component against which other components will be built.
Per-function unit tests are, of course, also great where appropriate!

Here's an example testing a component with a mocked dependency on `other` and on the `forwarder` bundle:

```go
func TestMyComponent(t *testing.T) {
    var comp Component
    var other other.Component
    comptest.FxTest(t,
        Module,               // use the real version of this component
        other.MockModule,     // use the mock version of another component in this bundle
        forwarder.MockBundle, // all forwarder components, mocked
        fx.Populate(&comp),   // get the instance of this component
        fx.Populate(&other),  // get the (mock) instance of the other component
    ).WithRunningApp(func() {
        // cast `other` to its mock interface to call mock-specific methods on it
        config.(config.Mock).SetConfig('foo', 'bar')             // Arrange (from core.MockBundle)
        other.(other.Mock).SetSomeValue(10)                      // Arrange
        comp.DoTheThing()                                        // Act
        require.Equal(t, 20, other.(other.Mock).GetSomeResult()) // Assert
    })
}
```

If the component has a mock implementation, it is a good idea to test that mock implementation as well.
