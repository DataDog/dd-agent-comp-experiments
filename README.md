# What's This?

This repo is an experiment in what a component-based Agent might look like.
It's not intended to actually behave like an Agent -- all non-component implementation is fake.
But it is intended to encompass the current design of the Agent, focusing on the different sorts of interconnections between components required.

This provides a loosely-sketched vision of what we will want to get to with the Agent, encouraging discussion now, during the design phase, but with concrete code to look at.

It is also a chance to try out some "interesting" ideas (such as config reducers) that would be much more difficult to sketch out in a real Agent.
Maybe some of these are bad ideas, maybe they are things we can do later, or maybe we want to do them first thing.

You can see the set of components in [`COMPONENTS.md`](./COMPONENTS.md).

# Building

To build, run `inv build`.
This will build all binaries in the root of the repository.

To regenerate files in the repo, run `inv generate`.

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
   To improve logging, use `fx.Module(comppath, ..)`, where `comppath` is the root-relative package path, e.g., `"comp/util/log"`.
   It should have a formulaic doc string like `// Module defines the fx options for this component.`

Components should not be nested; that is, no component's Go path should be a prefix of another component's Go path.

### Implementation

The completed `component.go` looks like this:

```go
// Package foo ... (detailed doc comment for the component)
package config

// team: some-team-name

// Component is the component type.
type Component interface {
	// Foo is ... (detailed doc comment)
	Foo(key string) string
}

// Module defines the fx options for this component.
var Module = fx.Module(
    "comp/foo", // (package path of the component)
    fx.Provide(newFoo),
)
```

The Component interface is implemented in another file by an unexported type with a sensible name such as `launcher` or `provider`.

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

func newFoo(deps dependencies) Component { ...  }
```

The constructor `newFoo` is an `fx` constructor, so it can refer to other types and expect them to be automatically supplied.
For very simple constructors, listing the dependencies inline is OK, but most will want to use the `dependencies` pattern shown above.
As an `fx` constructor, it can also take an `fx.Lifetime` argument and set up OnStart or OnStop hooks.

The constructor can return either `Component` if it is infallible, or `(Component, error)` if it could fail.
Failure will crash the agent with a suitable message.

Within the body of the constructor, it may call methods on other components, as long as that component allows calls to the methods during the setup phase.

### Parameterized Components

Some components require parameters before they are instantiated.
For example, `comp/config` requires the path to the configuration file so that it can be ready to answer config requests as soon as it is instantiated.
Other components may wish to provide different implementations depending on these parameters; for example, `comp/health` need not monitor anything if not in a running agent.

To support this, components can define a `pkg.ModuleParams` type and expect that it be supplied by the app.

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
    Params ModuleParams `optional:"true"`
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
type Mock interface {
    // Component methods are included in Mock.
    Component

    // AddedFoos returns the foos added by AddFoo calls on the mock implementation.
    AddedFoos() []Foo
}

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
    fx.Supply(foo.ModuleParams{
        MaxFoos: 13,
    }),
    ...
)
```

### Component Dependencies

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

# Conventions

## Package paths

 * `cmd/<appname>/main.go` -- entrypoint for each app
 * `pkg/<pkgnname>/...` utility types and functions, ["plain old data"](https://en.wikipedia.org/wiki/Passive_data_structure) types
 * `pkg/util/<utilname>/...` utility packages that are not in the form of a component
 * `comp/...` components
 * `comp/<bundlename>/...` component bundles

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

## IPC API Commands

Several commands, such as `agent status` or `agent config`, call the running Agent's IPC API and format the result.
Components implementing this pattern should generally have two similar methods, such as `GetStatus` and `GetStatusRemote`.
The first method gathers the data locally, and the second requests the same data via the IPC API.
The component should register with `comp/ipcapi` to provide the result of the first method over the IPC API.

This arrangement locates both the client and server sides of the IPC API in one module.
The command implementation (under `cmd/`) then simply calls `GetStatusRemote` and formats the result for display.

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

# Open Questions

## Subprocesses

We want to support running some "things" (we should have a term for this!) as subprocesses, as is currently done for trace-agent, system-probe, process-agent, and security-agent.
Should these be different binaries (as they are now), or the same binary with different arguments?

# TODO

 * [DONE] component.yml?
 * [DONE] actor model conventions
 * [DONE] subscription conventions
 * [DONE] guidelines for non-component stuff
 * [DONE] health reporting
 * [DONE] add Mocks to a component and try them out
 * [DONE] selecting among multiple implementations of the same component (e.g., Tagger)
 * [DONE] CLI / subcommands (`agent run`, etc.)
 * [DONE] use fx.In for complex constructors (maybe with component concrete type?)
 * [DONE] API
 * status output
 * tlm / expvars
 * startup for bundles? the hidden Invoke() of a constructor that does nothing is weird.
 * more tests
   * some kind of external test (`./test/`?)
 * more components
   * tagger
   * wlm
   * AD + plugins
   * Check runners
   * DSD parts
   * Remaining Logs-Agent parts
   * Some other agents??
   * Serializer
   * Forwarder
   * [DONE] Flares
 * Subprocesses?
