# Overview of Components

The Agent is structured as a collection of components working together.
Depending on how the binary is built, and how it is invoked, different components may be instantiated.
The behavior of the components depends on the Agent configuration.

Components are structured in a dependency graph.
For example, the comp/logs/agent component depends on the comp/core/config component to access Agent configuration.
At startup, a few top-level components are requested, and [Fx](./fx.md) automatically instantiates all of the required components.

## What is a Component?

Any well-defined portion of the codebase, with a clearly documented API surface, _can_ be a component.
As an aid to thinking about this question, consider four "levels" where it might apply:

1. Meta: large-scale parts of the Agent that use many other components. Example: DogStatsD or Logs-Agent.
2. Service: something that can be used at several locations (for example by different applications). Example: Forwarder.
3. Internal: something that is used to implement a service or meta component, but doesn't make sense outside the component. Examples: DogStatsD's TimeSampler, or a workloadmeta Listener.
4. Implementation: a type that is used to implement internal components. Example: Forwarder's DiskUsageLimit.

In general, meta and service-level functionality should always be implemented as components.
Implementation-level functionality should not.
Internal functionality is left to the descretio of the implementing team: it's fine for a meta or service component to be implemented as one large, complex component, if that makes the most sense for the team.

## Bundles

The previous section suggests there is a large and growing number of components, and listing those components out repeatedly could grow tiresome and cause bugs.
Component bundles provide a way to manipulate multiple components, usually at the meta or service level, as a single unit.
For example, while Logs-Agent is internally composed of many components, those can be addressed as a unit with `comp/logs.Bundle`.

## Apps

TODO

## Build-Time and Runtime Dependencies

TODO
