# Defining Components and Bundles

This file describes the mechanics of implementing components and bundles.

This guidelines in this file are quite prescriptive, with the intent of making all components "look the same".
This reduces cognitive load when using components -- no need to remember one component's peculiarities.
It also allows Agent-wide changes, where we make the same formulaic change to each component.
If a situation arises that contradicts the guidelines, then we can update the guidelines (and change all affected components).
In fact, many of these prescriptions can be easily verified by linters.

## Defining a Component

A component is defined in a dedicated package under `comp/<bundlename>/...`, where `<bundlename>` names the bundle that contains the component.
The package must have the following defined in `component.go`:

 * Extensive package-level documentation.
   This should define, as precisely as possible, the behavior of the component, acting as a contract on which users of the component may depend.
   See the "Documentation" section below for details.

 * A team-name comment of the form `// team: <teamname>`.
   This is used to generate CODEOWNERS information.

 * `componentName` -- the Go path of the component, relative to the repository root, e.g., `comp/core/health`.

 * `Component` -- the interface type implemented by the component.
   This is the type by which other components will require this one via `fx`.
   It can be an empty interface, if there is no need for any methods.
   It should have a formulaic doc string like `// Component is the component type.`, deferring documentation to the package docs.
   All interface methods should be exported and thoroughly documented.

 * `Module` -- an `fx.Option` that can be included in the bundle's `Module` or an `fx.App` to make this component available.
   To assist with debugging, use `fx.Module(componentName, ..)`.
   This item should have a formulaic doc string like `// Module defines the fx options for this component.`

Components should not be nested; that is, no component's Go path should be a prefix of another component's Go path.

### Implementation

The completed `component.go` looks like this:

```go
// Package foo ... (detailed doc comment for the component)
package config

// team: some-team-name

const componentName = "comp/foo"

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

The constructor can return either `provides`, if it is infallible, or `(provides, error)`, if it could fail.
An returned error will crash the agent at startup with a suitable message.
For simple constructors that return only the Component type, omitting the `provides` struct and just returning `Component` is perfectly fine.

The constructor may call methods on other components, as long as the called method's documentation indicates it is OK.

### Documentation

The documentation (both package-level and method-level) should include everything a user of the component needs to know.
In particular, any assumptions that might lead to panics if violated by the user should be clearly documented.

Treat extensive "how to use this component without introducing bugs" documentation as a code smell: simplifying the usage will improve the robustness of the Agent.

Include:

* Precise information on when each method may be called.
  Can methods be called concurrently?
  Are some methods invalid before the component has started?
  Such assumptions are difficult to verify, so where possible try to make every method callable concurrently, at all times.

* Precise information about data ownership of passed values and returned values.
  Users can assume that any mutable value returned by a component will not be modified by the user or the component after it is returned.
  Similarly, any mutable value passed to a component will not be later modified either by the component or the caller.
  Any deviation from these defaults should be clearly documented.
  It can be surprisingly hard to avoid mutating data -- for example, `append(..)` surprisingly mutates its first argument.
  It is also hard to detect these bugs, as they are often intermittent, cause silent data corruption, or introduce rare data races.
  Where performance is not an issue, prefer to copy mutable input and outputs to avoid any potential bugs.

* Precise information about goroutines and blocking.
  Users can assume that methods do not block indefinitely, so blocking methods should be documented as such.
  Methods that invoke callbacks should be clear about how the callback is invoked: is it OK for the callback to block?

* Precise information about channels.
  Is the channel buffered?
  What happens if the channel is not read from quickly enough, or if reading stops?
  Can the channel be closed by the sender, and if so, what does that mean?

### Testing Support

To support testing, components can optionally provide a mock implementation, with the following in `component.go`.

 * `Mock` -- the type implemented by the mock version of the component.
   This should embed `pkg.Component`, and provide additional exported methods for manipulating the mock for use by other packages.

 * `MockModule` -- an `fx.Option` that can be included in a test `App` to get the component's mock implementation.

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

## Defining a Bundle

A bundle is defined in a dedicated package name `comp/<bundlename>`.
The package must have the following defined in `bundle.go`:

 * Extensive package-level documentation.
   This should define:

     * The purpose of the bundle
     * What components are and are not included in the bundle.
       Components might be omitted in the interest of binary size, as discussed in the [overview](./components.md).
     * Which components are automatically instantiated.
     * Which other _bundles_ this bundle depends on.
       Bundle dependencies are always expressed at a bundle level.

 * A team-name comment of the form `// team: <teamname>`.
   This is used to generate CODEOWNERS information.

 * `componentName` -- the Go path of the component, relative to the repository root, e.g., `comp/core/health`.

 * `BundleParams` -- the type of the bundle's parameters (see below).
   This item should have a formulaic doc string like `// BundleParams defines the parameters for this bundle.`

 * `Bundle` -- an `fx.Option` that can be included in an `fx.App` to make this bundle's components available.
   To assist with debugging, use `fx.Module(componentName, ..)`.
   Use `fx.Invoke(func(componentpkg.Component) {})` to instantiate components automatically.
   This item should have a formulaic doc string like `// Module defines the fx options for this component.`

Typically, a bundle will automatically instantiate the top-level components that represent the bundle's purpose.
For example, the trace-agent bundle `comp/trace` might automatically instantiate `comp/trace/agent`.

### Bundle Parameters

Apps can provide some intialization-time parameters to bundles.
These parameters are limited to two kinds:

 * Parameters specific to the app, such as whether to start a network server; and
 * Parameters from the environment, such as command-line options.

Anything else is runtime configuration and should be handled vi `comp/core/config` or another mechanism.

To avoid Go package cycles, the `BundleParams` type must be defined in the bundle's internal package, and re-exported from the bundle package:

```go
// --- comp/<bundlename>/internal/params.go ---

// BundleParams defines the parameters for this bundle.
type BundleParams struct {
    ...
}

// --- comp/<bundlename>/bundle.go ---
import ".../comp/<bundlename>/internal"
// ...

// BundleParams defines the parameters for this bundle.
type BundleParams = internal.BundleParams
```

Components within the bundle can then require `internal.BundleParams` and modify their behavior appropriately:

```go
// --- comp/<bundlename>/foo/foo.go

func newFoo(..., params internal.BundleParams) provides {
    if params.HyperMode { ... }
}
```

See the AutoStart [convention](./conventions.md) for a common BundleParams field.

### Testing

A bundle should have a test file, `bundle_test.go`, to verify the documentation's claim about its dependencies.
This simply uses ValidateApp to check that all dependencies are satisfied when given the full set of required bundles.

```go
func TestBundleDependencies(t *testing.T) {
	require.NoError(t, fx.ValidateApp(
		fx.Supply(core.BundleParams{}),
		core.Bundle,
		fx.Supply(autodiscovery.BundleParams{}),
		autodiscovery.Bundle,
		fx.Supply(BundleParams{}),
		Bundle))
}
```

## Using Other Fx Types

It's fine to provide other, unexported `fx` types in `pkg.Module`, if that is helpful.
Because they are unexported, they will be invisible to users of the component.
