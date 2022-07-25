# What's This?

This repo is an experiment in what a component-based Agent might look like.
It's not intended to actually behave like an Agent -- all non-component implementation is fake.
But it is intended to encompass the current design of the Agent, focusing on the different sorts of interconnections between components required.

This provides a loosely-sketched vision of what we will want to get to with the Agent, encouraging discussion now, during the design phase, but with concrete code to look at.

It is also a chance to try out some "interesting" ideas (such as config reducers) that would be much more difficult to sketch out in a real Agent.
Maybe some of these are bad ideas, maybe they are things we can do later, or maybe we want to do them first thing.

You can see the set of components in [`COMPONENTS.md`](./COMPONENTS.md).

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

A component is defined in a dedicated package under `comp/`, with the following defined in `component.go`:

 * Extensive package-level documentation.
   This should define, as precisely as possible, the behavior of the component, acting as a contract on which users of the component may depend.

   The documentation (both package-level and method-level) should also include:

   * Precise information about which interface methods can be called during the setup phase, and which must only be called after the component is started.
   * Precise information about data ownership of passed values and returned values (does the method modify its arguments?  can the caller modify a returned slice or map?).
   * Precise information about goroutines and blocking (does the method block? is a callback invoked in a dedicated goroutine? what happens if a channel is full?)

 * `pkg.Component` -- the type implemented by the component.
   This can be an empty interface, but is the type by which other components will find this one via `fx`.
   It should have a formulaic doc string like `// Component is the component type.`, deferring documentation to the package docs.
   All interface methods should be exported and thoroughly documented.

 * `pkg.Module` -- an `fx.Option` that can be included in an `fx.App` to make this component available.
   This is sometimes as simple as `var Module = fx.Provide(new)` to inform `fx` about the constructor, but can use `fx.Options` or `fx.Module` as necessary.
   If using `fx.Module`, the first argument should be the root-relative package path, e.g., `"comp/util/log"`.
   It should have a formulaic doc string like `// Module defines the fx options for this component.`

Components should not be nested; that is, no component's Go path should be a prefix of another component's Go path.

### Implementation

The Component interface is implemented by an unexported type with a sensible name such as `launcher` or `provider`.

#### Constructor

The component type has a constructor with an appropriate, unexported name such as `newProvider`.
This is an `fx` constructor, so it can refer to other types and expect them to be automatically supplied:

```golang
func newProvider(log log.Component, config config.Component) Component { ...  }
```

Within the body of the constructor, it may call methods on other components, as long as the component allows calls to those methods during the setup phase.

As an `fx` constructor, it can also take an `fx.Lifetime` argument and set up OnStart and OnStop hooks.

The constructor is passed to `fx.Provide` in the definition of `Module` in `component.go`.

#### Other Fx Types

It's fine to provide other, unexported `fx` types in `pkg.Module`, if that is helpful.
Because they are unexported, they will be invisible to users of the component.

### Testing Support

To support testing, components can optionally provide a mock implementation, with the following in `component.go`.

 * `pkg.Mock` -- the type implemented by the mock version of the component.
   This should embed `pkg.Component`, and provide additional exported methods for manipulating the mock for use by other packages.

 * `pkg.MockModule` -- an `fx.Option` that can be included in a test `App` to get the component's mock implementation.

Here `pkg.MockModule` will typically provide a `newMock` constructor which creates a struct implementing the `pkg.Mock` interface, with no other dependencies.

## Using Components

### Apps (binaries)

Apps represent the final binaries defined by the agent.
Each "flavor" of agent is defined in a different sub-package of `cmd/`.

Apps are formulaic and should not contain any complex logic.
Their job is to parse command-line options, set up an `fx` App, and run it.

### Dependencies

Component dependencies are automatically determined from the arguments to a component constructor.
For example, a component that depends on the log component will have a `logs.Component` in its argument list:

```go
import "github.com/djmitche/dd-agent-comp-experiments/comp/util/log"

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
These components should be defined in an `internal/` package, and included in a `Modules` definition in the parent package.
This single `Modules` can then be included by apps that require that functionality.

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

### Non-Component Code

Code that is not part of a component can be placed under `pkg/`.

This includes
 * ["plain old data"](https://en.wikipedia.org/wiki/Passive_data_structure) types; and
 * Utility types and functions (either in a sub-package of `pkg/util`, or as a top-level `pkg` for more complex implementations),

# Conventions

## Subscriptions

A common mode of interaction between two components is for one to subscribe to notifications from the other.
The `pkg/util/subscriptions` package provides generic support for this.

Subscriptions are usually static for the lifetime of an Agent, and should be made during the setup phase, before components have started.
The documentation for components supporting subscriptions will make this clear.

## Actors

Components often take the form of an [actor](https://en.wikipedia.org/wiki/Actor_model): a dedicated goroutine that processes events in a loop.
This approach requires no concurrency controls, since all activity takes place in one goroutine.
It is also easy to test: start a goroutine, send it some events, and assert on the output.

The `pkg/util/actor` package supports components that use the actor structure, including connecting them to the `fx` life cycle.

## Health Monitoring

Components which can fail, and especially those using the actor model, should register with `comp/health` to monitor their health.
In this context, "failure" is a user-visible problem with the component that can occur after startup.
This may be related to resource exhaustion, user misconfiguration, or an issue in the environment.
Many components can't fail (or at least, we can't yet imagine how they would fail); these do not need to report to the `comp/health` component.

## Plugins

Plugins are things like Launchers, Tailers, Config Providers, Listeners, etc. where there are several implementations that perform the same job in different contexts.

We typically want to include different sets of plugins in different builds, differentiating at build time.
Just including `somepkg.Module` in an `fx.App` is enough to pull in that module's code, causing binary bloat, so these distinctions must be made at build time.
For example, a slimmed-down logs agent for limited systems might only support logging TCP inputs, while a full-fledged logs agent includes support for containers, syslog, files, and so on.

Plugins always "plug in" to some "manager" component (typically named `foomgr`), and should depend on that manager and register themselves with it at startup.
Then, it is up to apps to depend on the necessary plugins.

## Programming Errors

Programming errors, such as calling a method at an inappropriate time, should be handled with `panic(..)` instead of errors.
If an error is returned, it will likely be logged and may not be seen.
A panic, on the other hand, is very noticeable and carries a stack trace that can help the programmer figure out what they've missed.
Try to arrange for such panics to happen consistently, so that such programming errors are quick to find.

# Future Plans

## Component Linting

With good detection of components (already used to generate COMPONENTS.md and CODEOWNERS), we can check that the guidelines are followed.
For example, this check could easily verify that components are not nested, and that every component has a `Component` type and `Module` value.

## Config Reducers

We can use a concept similar to that defined by Redux to simplify the DD configuration used by each component that needs it.
This would involve a "reducer" that extracts data from Viper (or whatever we switch to) and places it in a component-specific struct.
Using struct tags and a utility function would allow for a very regular, greppable arrangement of configuration parameters with very little per-component boilerplate.
This will also ease testing of components: tests can simply provide a filled-in configuration struct, instead of manually setting configuration parameters.

# TODO

 * [DONE] component.yml?
 * [DONE] actor model conventions
 * [DONE] subscription conventions
 * [DONE] guidelines for non-component stuff
 * [DONE] health reporting
 * add Mocks to a component and try them out
 * selecting among multiple implementations of the same component (e.g., Tagger)
 * CLI / subcommands (`agent run`, etc.)
 * status output
 * tlm
 * API
 * use fx.In for complex constructors (maybe with component concrete type?)
