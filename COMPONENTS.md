# Agent Components

This file lists all components defined in this repository, with their package summary.
Click the links for more documentation.

## [comp/config](https://pkg.go.dev/github.com/djmitche/dd-agent-comp-experiments/comp/config)

*Datadog Team*: agent-shared-components

Package config implements a component to handle agent configuration.  This
component wraps Viper.

## [comp/health](https://pkg.go.dev/github.com/djmitche/dd-agent-comp-experiments/comp/health)

*Datadog Team*: agent-shared-components

Package health implements a component that monitors the health of other
components.

## [comp/logs/agent](https://pkg.go.dev/github.com/djmitche/dd-agent-comp-experiments/comp/logs/agent)

*Datadog Team*: agent-metrics-logs

Package agent implements a component representing the logs agent.  This
component coordinates activity related to gathering logs and forwarding them
to the Datadog intake.

## [comp/logs/internal/sourcemgr](https://pkg.go.dev/github.com/djmitche/dd-agent-comp-experiments/comp/logs/internal/sourcemgr)

*Datadog Team*: agent-metrics-logs

Package sourcemgr implements a component managing logs-agent sources (type
LogSource).  It receives additions and removals of sources from other
components, and it informs subscribers of these additions and removals.

## [comp/logs/launchers/file](https://pkg.go.dev/github.com/djmitche/dd-agent-comp-experiments/comp/logs/launchers/file)

*Datadog Team*: agent-metrics-logs

Package file implements a launcher that responds to file sources by starting
tailers for the indicated files.

## [comp/logs/launchers/launchermgr](https://pkg.go.dev/github.com/djmitche/dd-agent-comp-experiments/comp/logs/launchers/launchermgr)

*Datadog Team*: agent-metrics-logs

Package launchermgr implements a component managing logs-agent launchers.  It collects
the set of loaded launchers during start-up, and allows enumeration and retrieval
as necessary.

## [comp/util/log](https://pkg.go.dev/github.com/djmitche/dd-agent-comp-experiments/comp/util/log)

*Datadog Team*: agent-shared-components

Package log implements a component to handle logging internal to the agent.
