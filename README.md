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

# Documentation

The [documentation](./doc) covers how to build and use components, including a description of Fx, the framework used to coordinate component interactions.
It will become a part of the Agent's developer documentation, so bears close attention when reviewing this repository.
