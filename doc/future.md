# Future Plans

## Component Linting

With good detection of components (already used to generate COMPONENTS.md and CODEOWNERS), we can check that the guidelines are followed.
For example, this check could easily verify that components are not nested, and that every component has a `Component` type and `Module` value.

## Config Reducers

We can use a concept similar to that defined by Redux to simplify the DD configuration used by each component that needs it.
This would involve a "reducer" that extracts data from Viper (or whatever we switch to) and places it in a component-specific struct.
Using struct tags and a utility function would allow for a very regular, greppable arrangement of configuration parameters with very little per-component boilerplate.
This will also ease testing of components: tests can simply provide a filled-in configuration struct, instead of manually setting configuration parameters.

