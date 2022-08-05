# Open Questions

## Remote Components

[The RFC](https://github.com/DataDog/architecture/blob/master/rfcs/agent-component-architecture/rfc.md#remote-component-considerations) mentions remote components, but without much detail.
In general, this has been a "stretch" goal -- we are mostly building an in-process, Go API between components and not a cross-process API with all of the complexity that brings.

There are at least two areas where we may want to consider fledgling "remote component" support, though:

 * Several components publish some information on the internal IPC API (5001) that is accessed by the GUI and CLI.
   Should the `agent status` command start a remote "version" of the `status` component that has the same Go API and somehow automatically links to the running agent?
   Or should we do something simpler, treating these as client and server sides of a (non-component) API?

 * In the existing Agent, the tagger can either operate locally or use a remote tagger, with the same interface.
   So there is at least some prior art to supporting remote versions of components.

# TODO

 * Issues for people to hack on
 * How do we use MockModule in the face of Bundles?
 * Put IPCAPI client stuff in cmd/ instead of the same component
 * OneShot should use a Hook, rather than run things before startup
 * Registration types in registering components should be fx.Out and have `group:"true"` defined
 * Put LivenessMonitior in actor pkg
 * put all options in a variable for each app, and validate in a unit test (@ogaca-dd)
   * Maybe a common.RunCobraCommand?
 * > As the pattern to use fxtest.New + fx.Populate + defer app.RequireStart().RequireStop() will be common in unit tests, what do you think about providing a function for hiding this complexity?
   (@ogaca-dd)
 * doc that flare cb will be called concurrently
 * make IfConfigured the default AutoStart value
 * 'actor' field is unnecessary in actor components
 * more tests
   * some kind of external test (`./test/`?)
 * more components
   * tagger **soon - open questions about local/remote**
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
