# Conventions

This section contains a set of conventions for handling common situations.
Each subsection can be read independently.

## Registrations

Components generally need to talk to one another!
In simple cases, that occurs by method calls.
But in many cases, a single component needs to communicate with a number of other components that all share some characteristics.
For example, the `comp/core/health` component monitors the health of many other components, and `comp/workloadmeta/scheduler` provides workload events to an arbitrary number of subscribers.

The convention in the Agent codebase is to use [value groups](./fx.md#value-groups) to accomplish this.
The _collecting_ component requires a slice of some _collected type_, and the _providing_ components provide values of that type.
Consider a simple case of a server component to which endpoints can be attached.
The server is the collecting component, requiring a slice of type `[]*endpoint`, where `*endpoint` is the collected type.
Providing components provide values of type `*endpoint`.

The convention is to "wrap" the collected type in a struct type which embeds `fx.Out` and has tag `group:"pkgname"`, where `pkgname` is the short package name (Fx requires a group name, and this is as good as any).
This helps providing components avoid the common mistake of omitting the tag.
The collected type can be exported, such as if it has useful methods for providing components to call, but can also be unexported.

The collecting component should define the registration type and a constructor for it.

```go
// --- server/component.go ---

// ...
// Server endpoints are provided by other components, by providing a server.Registration
// instance.
//  ...
package server

type Registration struct {
    fx.Out

    Endpoint endpoint `group:"server"`
}

// NewRegistration creates a new Registration instance for the given endpoint.
func NewRegistration(route string, handler func()) Registration { .. }
```

Its implementation then requires a slice of the collected type, again using `group:"server"`:

```go
// --- server/server.go ---

// endpoint defines an endpoint on this server.
type endpoint struct { .. }

type dependencies struct {
    fx.In

    Registrations []endpoint `group:"server"`
}

func newServer(deps dependencies) Component {
    // ..
    for _, e := range deps.Registrations {
        if e.handler == nil {
            continue
        }
        // ..
    }
    // ..
}
```

It's good practice to ignore zero values, as that allows providing components to skip
the registration if desired.

Finally, the providing component (in this case, `foo`) includes a registration in its output:

```go
// --- foo/foo.go ---
func newFoo(deps dependencies) (Component, server.Registration) {
    // ..
    return foo, server.NewRegistration("/things/foo", foo.handler)
}
```

This technique has some caveats to be aware of:

 * The providing components are instantiated before the collecting component.
 * Fx treats value groups as the collecting component depending on all of the providing components.
   This means that the providing components cannot depend on the collecting component.
 * Fx will instantiate _all_ providing components before the collecting component.
   This may lead to components being instantiated in unexpected circumstances.
   The AutoStart convention is meant to partially address this issue.

## Subscriptions

Subscriptions are a common form of registration, and have support in the `pkg/util/subscriptions` package.

To implement a subscription, the collecting component defines a message type.
Providing components provide `subscriptions.Subscription[coll.Message]`, from which they can obtain a `subscriptions.Receiver[coll.Message]`.
Collecting components require `subcriptions.Publisher[coll.Messag]`, from which they can obtain a `subscriptions.Transmitter[coll.Message]`.

```go
// --- announcer/component.go ---

// ...
// To subscribe to these announcements, provide a subscriptions.Subscription[announcer.Announcement].
// ...
package announcer
```

```go
// --- announcer/announcer.go ---

func newAnnouncer(pub subscriptions.Publisher[Anouncement]) Component {
    return &announcer{announcementTx: pub.Transmitter()}  // (get a Transmitter from the Publisher)
}

// .. later send messages with 
    ann.eventTx.Notify(a)
```

```go
// --- listener/listener.go ---

func newListener() (Component, subscriptions.Subscription[announcer.Announcement]) {
    sub := subscriptions.NewSubscription[Event]()
    return &listener{eventRx: sub.Receiver}, sub
}

// .. later receive messages with
    a := <- l.eventRx.Chan()
```

If a receiving component decides it does not want to subscribe after all (such
as, if it is not started), it can return the zero value,
`subscriptions.Subscription[Event]{}`, from its constructor.

See the `pkg/util/subscriptions` documentation for more details.

## Actors

Components often take the form of an [actor](https://en.wikipedia.org/wiki/Actor_model): a dedicated goroutine that processes events in a loop.
This approach requires no concurrency controls, since all activity takes place in one goroutine.
It is also easy to test: start a goroutine, send it some events, and assert on the output.

The `pkg/util/actor` package supports components that use the actor structure, including connecting them to the `fx` life cycle and automatic "liveness monitoring".

A component structured as an actor typically looks like this:

```go
func newThing(lc fx.Lifecycle) (Component, health.Registration) {
    healthReg := health.NewRegistration(componentName)
    t := &thing{..}
    actor := actor.New()
    actor.HookLifecycle(lc, t.run)
    actor.MonitorLiveness(healthReg.Handle, time.Second)
    return thing, healthReg
}

func (t *thing) run(ctx context.Context, alive <-chan struct{}) {
    for {
        select {
            // .. receive from some component specific channels
            case <-alive:
            case <-ctx.Done():
                return
        }
    }
}
```

## Component Auto-Startup

It's easy for a component to be instantiated unexpectedly, if it is an indirect dependency of another component that is needed in a particular app.
Value groups can also lead to unexpected instantiation.

To avoid surprises, components do not automatically start just because they were instantiated.
Instead, BundleParams for most bundles include an `AutoStart` field which dictates whether the component should or should not start up.
"One-shot" apps set this value to `Never`, disabling any incidentally-required components from actually starting.

The type is defined in `pkg/util/startup`, and can have the values `Always`, `Never`, or `IfConfigured`.
To support the `IfConfigured` value, bundles which take configuration can define a `ShouldStart` method on `BundleParams`:

```go
// ShouldStart determines whether the bundle should start, based on
// configuration.
func (p BundleParams) ShouldStart(config config.Component) bool {
    return p.AutoStart.ShouldStart(config.GetBool("foo_agent.enabled"))
}
```

Each component in the bundle that has any active functionality then consults this function in its constructor, for example:

```go
func newFoo(deps dependencies) provides {
    f := &foo { .. }
    if deps.Params.ShouldStart(deps.Config) {
        f.actor.HookLifecycle(deps.Lc, l.run)
        f.subscription = eventpub.Subscribe()
    }
    return provides{
        Component:    f,
        Subscription: f.subscription,
    }
}
```

Note that l.subscription is nil if the component does not start, meaning that the component will not receive messages, which might cause the event publisher to block.

## IPC API Commands

Several commands, such as `agent status` or `agent config`, call the running Agent's IPC API and format the result.
Components implementing the data used by these commands generally have a method to get the data, such as `GetStatus`, and also publish this information over the IPC API.

The subcommand (e.g., `agent status`) implements the client side of this transaction, using `comp/core/ipc/ippclient` to fetch, format, and display the data.

## Health Monitoring

Components which can fail, and especially those using the actor model, should register with `comp/core/health` to monitor their health.
In this context, "failure" is a user-visible problem with the component that can occur after startup.
This may be related to resource exhaustion, user misconfiguration, or an issue in the environment.
Many components can't fail (or at least, we can't yet imagine how they would fail); these do not need to report to the `comp/core/health` component.

## Binary and App Common Support

(This support needs more development)

Most apps can include `cmd/common.SharedOptions(..)` in their `fx.App` to provide common component bundles and their parameters.
This takes a `oneShot` argument which distinguishes one-shot apps like `agent flare` from long-running apps like `trace-agent run`.

One-shot apps can use `common.OneShot` to run a function and shut down the app when it completes.
Long-running apps can use `common.RunApp`, which takes care to return an error on failure, rather than calling `os.Exit` as Fx's `app.Run` does.

## Non-Component Code

Code that does not directly implement a component, but provides supporting functionality or ["plain old data"](https://en.wikipedia.org/wiki/Passive_data_structure) structures, should be in `pkg/`.
This might be a top-level `pkg/` directory or, for utilities, under `pkg/util/`.

