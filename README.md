# What's This?

This repo is an experiment in what a component-based Agent might look like.
It's not intended to actually behave like an Agent -- all non-component implementation is fake.
But it is intended to encompass the current design of the Agent, focusing on the different sorts of interconnections between components required.

This provides a loosely-sketched vision of what we will want to get to with the Agent, encouraging discussion now, during the design phase, but with concrete code to look at.

It is also a chance to try out some "interesting" ideas (such as config reducers) that would be much more difficult to sketch out in a real Agent.
Maybe some of these are bad ideas, maybe they are things we can do later, or maybe we want to do them first thing.

# Component Guidelines

What follows are draft guidelines for writing components.
This repository follows these guidelines.

If these guidelines are adopted, then the expectation is that all components would follow them.
This reduces cognitive load when using components -- no need to remember one component's peculiarities.
It also allows Agent-wide changes, where we make the same formulaic change to each component.
If a situation arises that contradicts the guidelines, then we can update the guidelines (and change all affected components).

## Prerequisites

You should be familiar with [`fx`](https://pkg.go.dev/go.uber.org/fx), the tool used for dependency injection.

## Component Definition

A component is defined in a dedicated package, with the following defined in `component.go`:

 * Extensive package-level documentation.
   This should define, as precisely as possible, the behavior of the component, acting as a contract on which users of the component may depend.
   The documentation should also include:

   * Precise information about which interface methods can be called during the setup phase, and which must only be called after the component is started.

 * `pkg.Component` -- the type implemented by the component.
   This can be an empty interface, but is the type by which other components will find this one via `fx`.
   It should have a formulaic doc string like `// Component is the component type.`, deferring documentation to the package docs.
   All interface methods should be exported and thoroughly documented.

 * `pkg.Module` -- an `fx.Option` that can be included in an `fx.App` to make this component available.
   This is sometimes as simple as `var Module = fx.Provide(new)` to inform `fx` about the constructor, but can use `fx.Options` or `fx.Module` as necessary.
   It should have a formulaic doc string like `// Module defines the fx options for this component.`

Any other exported types relevant to the component should also be included in `component.go`.
This ensures that the source file itself is a useful reference, in addition to Godoc-generated documentation.

### Implementation

The component implementation begins with a constructor, `new`.
This is unexported because other packages will not call it directly, but via `fx` requirements in other constructors or `fx.Invoke` calls.
This is an `fx` constructor, so it can refer to other types and expect them to be automatically supplied:

```golang
func new(log logpkg.Component, config configpkg.Component) Component { ...  }
```

As an `fx` constructor, it can also take an `fx.Lifetime` argument and set up OnStart and OnStop hooks.

### Testing Support

To support testing, components can optionally provide a mock implementation, with the following in `component.go`.

 * `pkg.Mock` -- the type implemented by the mock version of the component.
   This should embed `pkg.Component`, and provide additional exported methods for manipulating the mock for use by other packages.

 * `pkg.MockModule` -- an `fx.Option` that can be included in a test `App` to get the component's mock implementation.

Here `pkg.MockModule` will typically provide a `newMock` constructor which creates a struct implementing the `pkg.Mock` interface, with no other dependencies.

## Using Components

### Apps (binaries)

### Dependencies

### Testing

```go
func TestMyComponent(t *testing.T) {
    var comp Component
    var other otherpkg.Component
    app := fxtest.New(t,
        Module,
        otherpkg.MockModule, // use the mock version of otherpkg
        fx.Populate(&comp),
        fx.Populate(&other),
    )
    defer app.RequireStart.RequireStop()

    other.(otherpkg.Mock).SetSomeValue(10)                      // Arrange
    comp.DoTheThing()                                           // Act
    require.Equal(t, 20, other.(otherpkg.Mock).GetSomeResult()) // Assert
}
```

# Conventions

## Config Reducers

## Subscriptions

## Plugins

Plugins are things like Launchers, Tailers, Config Providers, Listeners, etc. where there are several implementations that perform the same job in different contexts.

We typically want to include different sets of plugins in different builds, differentiating at build time.
For example, a slimmed-down logs agent for limited systems might only support logging TCP inputs, while a full-fledged logs agent includes support for containers, syslog, files, and so on.

Plugins always "plug in" to some collection, and should depend on that collection and register themselves with that collection at startup.
Then, it is up to apps to depend on the necessary plugins.
TODO: ^^ this might be weird-looking, since the app never _does_ anything.  Maybe a "register" method?

# TODO

 * component.yml?
     * Team ownership
     * COMPONENTS.md
     * Component linting
 * should _all_ components be linked in apps, or can a component that has some "dedicated" components include those in its `fx.Module`?
   * maybe everything lists its deps, but then how do you inject testing deps?
 * add Mocks to a compoent and try them out
 * nesting is allowed, right?
 * everything under comp/ or pkg/ or whatever?
