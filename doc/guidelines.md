# Component Guidelines

What follows are draft guidelines for writing components.
This repository follows these guidelines.

If these guidelines are adopted, then the expectation is that all components would follow them.
This reduces cognitive load when using components -- no need to remember one component's peculiarities.
It also allows Agent-wide changes, where we make the same formulaic change to each component.
If a situation arises that contradicts the guidelines, then we can update the guidelines (and change all affected components).

## Component Definition

A component is defined in a dedicated package under `comp/`, with the following defined in `component.go`:

 * Extensive package-level documentation.
   This should define, as precisely as possible, the behavior of the component, acting as a contract on which users of the component may depend.

   The documentation (both package-level and method-level) should also include:

   * Precise information about which interface methods can be called during the setup phase, and which must only be called after the component is started.
   * Precise information about data ownership of passed values and returned values.
     By default, any mutable value returned by a component will not be modified after it is returned.
     Similarly, any mutable value passed to a component will not be later modified either by the component or the caller.
     Any deviation from these defaults should be clearly documented.
   * Precise information about goroutines and blocking (does the method block? is a callback invoked in a dedicated goroutine? what happens if a channel is full?).
     By default, all methods are assumed to return without blocking.

 * `pkg.Component` -- the type implemented by the component.
   This can be an empty interface, but is the type by which other components will find this one via `fx` and so must still appear in a function signature.
   It should have a formulaic doc string like `// Component is the component type.`, deferring documentation to the package docs.
   All interface methods should be exported and thoroughly documented.

 * `pkg.Module` -- an `fx.Option` that can be included in an `fx.App` to make this component available.
   To improve logging, use `fx.Module(comppath, ..)`, where `comppath` is the root-relative package path, e.g., `"comp/core/log"`.
   It should have a formulaic doc string like `// Module defines the fx options for this component.`

Components should not be nested; that is, no component's Go path should be a prefix of another component's Go path.

### Implementation

The completed `component.go` looks like this:

```go
// Package foo ... (detailed doc comment for the component)
package config

// team: some-team-name

const componentname = "comp/foo" // ... (should match Go package name)

// Component is the component type.
type Component interface {
	// Foo is ... (detailed doc comment)
	Foo(key string) string
}

// Module defines the fx options for this component.
var Module = fx.Module(
    componentName,
    fx.Provide(newFoo),
)
```

The Component interface is implemented in another file by an unexported type with a sensible name such as `launcher` or `provider` or `foo`.

```go
package config

type foo {
    foos []string
}

// Foo implements Component#Foo.
func (f *foo) Foo(key string) string { ... }

type dependencies struct {
    fx.In

    Log log.Component
    Config config.Component
    // ...
}

type provides struct {
    fx.Out

    Component
}

func newFoo(deps dependencies) provides { ...  }
```

The constructor `newFoo` is an `fx` constructor, so it can refer to other types and expect them to be automatically supplied.
For very simple constructors, listing the dependencies inline is OK, but most will want to use the `dependencies` pattern shown above.
As an `fx` constructor, it can also take an `fx.Lifetime` argument and set up OnStart or OnStop hooks.

The constructor can return either `provides` if it is infallible, or `(provides, error)` if it could fail.
Failure will crash the agent with a suitable message.
For simple constructors that return only the Component type, omitting the `provides` struct and just returning `Component` is perfectly fine.

Within the body of the constructor, it may call methods on other components, as long as that component allows calls to the methods during the setup phase.

### Parameterized Components

Some components require parameters before they are instantiated.
For example, `comp/core/config` requires the path to the configuration file so that it can be ready to answer config requests as soon as it is instantiated.
Other components may wish to provide different implementations depending on these parameters; for example, `comp/core/health` need not monitor anything if not in a running agent.

To support this, components can define a `pkg.ModuleParams` type and allow apps to (optionally) supply it.

```go
// ModuleParams are the parameters to Module.
type ModuleParams struct {
	// ConfFilePath is the path to the configuration file.
	ConfFilePath string
}
```

..and accept that type in the constructor, but optionally:

```go
type dependenices {
    fx.In

    // ...
    Params *ModuleParams `optional:"true"`
}
func newFoo(deps dependencies) { ...  }
```

The dependency should be optional to more easily support tests for other components depending on this one, which will typically want default behaviors.

### Testing Support

To support testing, components can optionally provide a mock implementation, with the following in `component.go`.

 * `pkg.Mock` -- the type implemented by the mock version of the component.
   This should embed `pkg.Component`, and provide additional exported methods for manipulating the mock for use by other packages.

 * `pkg.MockModule` -- an `fx.Option` that can be included in a test `App` to get the component's mock implementation.

```go

// Mock implements mock-specific methods.
type Mock interface {
    // Component methods are included in Mock.
    Component

    // AddedFoos returns the foos added by AddFoo calls on the mock implementation.
    AddedFoos() []Foo
}

// MockModule defines the fx options for the mock component.
var MockModule = fx.Module(
    "comp/foo",
    fx.Provide(newMockFoo),
)
```

The `newMockFoo` constructor should create an implementation of the Mock interface.

#### Other Fx Types

It's fine to provide other, unexported `fx` types in `pkg.Module`, if that is helpful.
Because they are unexported, they will be invisible to users of the component.

## Using Components

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
