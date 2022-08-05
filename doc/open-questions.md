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
