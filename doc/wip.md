# Open Questions

## Subprocesses

We want to support running some "things" (we should have a term for this!) as subprocesses, as is currently done for trace-agent, system-probe, process-agent, and security-agent.
Should these be different binaries (as they are now), or the same binary with different arguments?

It may be that this is unrelated to components -- need more info.

## Remote Components

[The RFC](https://github.com/DataDog/architecture/blob/master/rfcs/agent-component-architecture/rfc.md#remote-component-considerations) mentions remote components, but without much detail.
In general, this has been a "stretch" goal -- we are mostly building an in-process, Go API between components and not a cross-process API with all of the complexity that brings.

There are at least two areas where we may want to consider fledgling "remote component" support, though:

 * Several components publish some information on the internal IPC API (5001) that is accessed by the GUI and CLI.
   Should the `agent status` command start a remote "version" of the `status` component that has the same Go API and somehow automatically links to the running agent?
   Or should we do something simpler, treating these as client and server sides of a (non-component) API?

 * In the existing Agent, the tagger can either operate locally or use a remote tagger, with the same interface.
   So there is at least some prior art to supporting remote versions of components.

## Enabling / Disabling and Component Bundles

Consider configs like `logs_enabled` or `apm_config.enabled` -- when these are false, it means that several components that are already instantiated should not actually start.
One options here is to thread a bunch of `Enable` methods through these components.
Then an app would call
```go
fx.Invoke(func(config config.Component, agent logsagent.Component) {
    if config.GetBool("logs_enabled") {
        logsagent.Enable()
    }
})
```

and the logs agent's Enable method would call the Enable method on its internal components (schedulers, etc. -- anything that does something active).
This approach is verbose -- writing `Enable` methods everywhere, adding an `enable bool` field, and then checking that field all over the place.
It also complicates things like health monitoring, status, and ipcserver -- those all require registration during startup, _before_ it's known whether the component is enabled.
So we'll need a way for all of those to handle disabled components.

See [#2](https://github.com/djmitche/dd-agent-comp-experiments/pull/2) for an attempt at solving this (which fails because the Health registration comes after start).

# TODO

 * Docs
     * bundles (and remove module params)
       * not optional
     * doc kinds of components, whether they should use Fx: https://github.com/djmitche/dd-agent-comp-experiments/pull/1#discussion_r936350828
     * nil subscriptions
 * OneShot should use a Hook, rather than run things before startup
 * put all options in a variable for each app, and validate in a unit test (@ogaca-dd)
   * Maybe a common.RunCobraCommand?
 * What happens when a subscriber doesn't consume because it's not started?
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
 * tlm / expvars
 * Subprocesses?
