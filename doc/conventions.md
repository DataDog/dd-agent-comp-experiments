# Conventions

## Package paths

 * `cmd/<appname>/main.go` -- entrypoint for each app
 * `pkg/<pkgnname>/...` utility types and functions, ["plain old data"](https://en.wikipedia.org/wiki/Passive_data_structure) types
 * `pkg/util/<utilname>/...` utility packages that are not in the form of a component
 * `comp/<bundlename>` bundles
 * `comp/<bundlename>/...` components (all components are in bundles)

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
The component should plugin to `comp/core/ipc/ipcserver` to provide the result of the first method over the IPC API.

This arrangement locates both the client and server sides of the IPC API in one module.
The command implementation (under `cmd/`) then simply calls `GetStatusRemote` and formats the result for display.

## Health Monitoring

Components which can fail, and especially those using the actor model, should register with `comp/core/health` to monitor their health.
In this context, "failure" is a user-visible problem with the component that can occur after startup.
This may be related to resource exhaustion, user misconfiguration, or an issue in the environment.
Many components can't fail (or at least, we can't yet imagine how they would fail); these do not need to report to the `comp/core/health` component.

## Plugins

Plugins are things like Launchers, Tailers, Config Providers, Listeners, etc. where there are several implementations that perform the same job in different contexts.

We typically want to include different sets of plugins in different builds, differentiating at build time.
Just including `somepkg.Module` in an `fx.App` is enough to pull in that module's code, causing binary bloat, so these distinctions must be made at build time.
For example, a slimmed-down logs agent for limited systems might only support logging TCP inputs, while a full-fledged logs agent includes support for containers, syslog, files, and so on.

Plugins always "plug in" to some "manager" component (typically named `foomgr`), and should depend on that manager and register themselves with it at startup.
Then, it is up to apps to depend on the necessary plugins.

This is accomplished using Fx's "value groups", where the plugins are all members of the same value group.
The manager component will define the name and type (typically `Registration`) of the group.
The plugins then use an `fx.Out` struct to register:

```go
type provides struct {
    fx.Out
    Component
    FooReg *foomgr.Registration `group:"foo"`
} //                 IMPORTANT! ^^^^^^^^^^^^^
  // Without this tag, build will succeed but the registration will be ignored

func newBar(deps dependencies) provides {
    comp := bar {
        fooReg: foomgr.NewRegistration(..),
        // ...
    }
    return provides {
        Component: comp,
        FooReg: comp.fooReg,
    }
}
```

Here the `provides` struct indicates that the `newBar` constructor returns multiple values, including the component itself and a foomgr.Registration in the group "foo".

The foomgr component will then gather all of the defined Registrations with

```go
type dependencies struct {
    fx.In
    Registrations []*Registration
    // ..
}
```

