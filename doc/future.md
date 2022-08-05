# Future Plans

## Component Linting

With good detection of components (already used to generate COMPONENTS.md and CODEOWNERS), we can check that the guidelines are followed.
For example, this check could easily verify that components are not nested, and that every component has a `Component` type and `Module` value.

## Config Reducers

We can use a concept similar to that defined by Redux to simplify the DD configuration used by each component that needs it.
This would involve a "reducer" that extracts data from Viper (or whatever we switch to) and places it in a component-specific struct.
Using struct tags and a utility function would allow for a very regular, greppable arrangement of configuration parameters with very little per-component boilerplate.
This will also ease testing of components: tests can simply provide a filled-in configuration struct, instead of manually setting configuration parameters.

## Component Reconfiguration and Restart

Since we have per-component health monitoring, it may be useful to be able to react automatically to unhealthy cmoponents, perhaps by restarting them.
This would require a more sophisticated lifecycle implementation than that provided by Fx, but `fx.Lifecyle`'s deign is a good place to start.

We may also want to support dynamic reconfiguration of the Agent.
This would require

* A way to determine what components are impacted by a configuration change (perhaps via config reducers)
* A way to adopt new configuration within a component
  * For trivial cases this may entail simply changing a field in the component's struct, or performing some internal adjustment such as starting or stopping workers.
  * Some components may be able to restart in-place without losing data
  * Some components may need to perform a complex dance to continue processing data with a new instance without interruption.

In general, this will be _very_ difficult to test thorougly.
