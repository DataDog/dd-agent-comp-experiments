# Conventions

This section contains a set of conventions for handling common situations.
Each subsection can be read independently.

## Registrations

Components generally need to talk to one another!
In simple cases, that occurs by method calls.
But in many cases, a single component needs to communicate with a number of other components that all share some characteristics.
For example, the `comp/core/health` component monitors the health of many other components, and `comp/workloadmeta /scheduler` provides workload events to an arbitrary number of subscribers.

The convention in the Agent codebase is to use [value groups](./fx.md#value-groups) to accomplish this.
The _collecting_ component requires a slice of some _registration type_, and the _providing_ components provide that registration type.
Consider a simple case of a server component to which endpoints can be attached.
The server is the collecting component, requiring a slice of type `[]*Endpoint`, where `*Endpoint` is the registration type.
Providing components provide values of type `*Endpoint`.

The collecting component should define the registration type and a constructor for it.

```go
// --- server/component.go ---

// Endpoint is provided by components that wish to register an endpoint
// with this server.  If a nil *Endpoint is provided, it will be ignored.
type Endpoint struct { .. }

// NewEndpoint creates a new Endpoint.
func NewEndpoint(route string, handler func()) *Endpoint { .. }
```

Its implementation then requires a slice of the registration type, using `group:"true"`:

```go
// --- server/server.go ---
type dependencies struct {
    fx.In

    Endpoints []*Endpoint `group:"true"`
}

func newServer(deps dependencies) Component {
    // ..
    for _, e := range deps.Endpoints {
        if e == nil {
            continue
        }
        // ..
    }
    // ..
}
```

It's good practice to ignore nil values, as that allows providing components to skip
the registration if desired.

Finally, the providing component (in this case, `foo`) includes a registration in its output:

```go
// --- foo/foo.go ---
type provides struct {
    fx.Out

    Component
    ServerRegistration *server.Registration `group:"true"`
}

func newFoo(deps dependencies) provides {
    // ..
    return provides{
        Component: foo,
        ServerRegistration: server.NewRegistration("/things/foo", foo.handler),
    }
}
```

This technique has some caveats to be aware of:

 * The providing components are instantiated before the collecting component.
 * Fx treats value groups as the collecting component depending on all of the providing components.
   This means that the providing components cannot depend on the collecting component.
 * Fx will instantiate _all_ providing components before the collecting component.
   This may lead to components being instantiated in unexpected circumstances.
   The AutoStart convention is meant to partially address this issue.
 * Omitting the `group:"true"` in either place it appears above will lead to Fx silently ignoring the registration.

## Subscriptions

Subscriptions are a common form of registration, and have support in the `pkg/util/subscriptions` package.

To implement a subscription, the collecting component defines a message type, a Subscription type, and a Subscribe function:

```go
// --- eventpub/component.go ---
type Event struct { .. }

// Subscription is provided by components that wish to subscribe to events
// from this component.  A nil *Subscription will be ignored.
type Subscription = subscriptions.Subscription[Event]

// Subscribe creates a new Subscription.
func Subscribe() Subscription {
    return subscriptions.NewSubscription[Event]()
}

// --- eventpub/eventpub.go ---
type eventpub struct {
    // .. (*eventpub implements Component)
    subscriptions *subscriptions.SubscriptionPoint[Event]
}

type dependencies struct {
    fx.In

    // ...
    Subscriptions []Subscription `group:"true"`
}

func newSourceMgr(deps dependencies) Component {
    ep := &eventpub{
        subscriptions: subscriptions.NewSubscriptionPoint[Event](deps.Subscriptions),
    }
    // ...
}

func (ep *eventpub) publishEvent(evt Event) {
    // notify all subscribers
    ep.subscriptions.Notify(evt)
}
```

See the `pkg/util/subscriptions` documentation for more details.

## Actors

Components often take the form of an [actor](https://en.wikipedia.org/wiki/Actor_model): a dedicated goroutine that processes events in a loop.
This approach requires no concurrency controls, since all activity takes place in one goroutine.
It is also easy to test: start a goroutine, send it some events, and assert on the output.

The `pkg/util/actor` package supports components that use the actor structure, including connecting them to the `fx` life cycle.

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

