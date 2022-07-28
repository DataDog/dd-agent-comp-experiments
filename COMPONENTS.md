# Agent Components

This file lists all components defined in this repository, with their package summary.
Click the links for more documentation.

## [comp/config](https://pkg.go.dev/github.com/djmitche/dd-agent-comp-experiments/comp/config@v0.0.2)

*Datadog Team*: agent-shared-components

Package config implements a component to handle agent configuration.  This
component wraps Viper.

## [comp/flare](https://pkg.go.dev/github.com/djmitche/dd-agent-comp-experiments/comp/flare@v0.0.2)

*Datadog Team*: agent-shared-components

Package flare implements a component creates flares for submission to support.

## [comp/health](https://pkg.go.dev/github.com/djmitche/dd-agent-comp-experiments/comp/health@v0.0.2)

*Datadog Team*: agent-shared-components

Package health implements a component that monitors the health of other
components.

## [comp/ipcapi](https://pkg.go.dev/github.com/djmitche/dd-agent-comp-experiments/comp/ipcapi@v0.0.2)

*Datadog Team*: agent-shared-components

Package ipcapi implements a component to manage the IPC API server and act
as a client.

## [comp/logs/agent](https://pkg.go.dev/github.com/djmitche/dd-agent-comp-experiments/comp/logs/agent@v0.0.2)

*Datadog Team*: agent-metrics-logs

Package agent implements a component representing the logs agent.  This
component coordinates activity related to gathering logs and forwarding them
to the Datadog intake.

## [comp/logs/internal/sourcemgr](https://pkg.go.dev/github.com/djmitche/dd-agent-comp-experiments/comp/logs/internal/sourcemgr@v0.0.2)

*Datadog Team*: agent-metrics-logs

Package sourcemgr implements a component managing logs-agent sources (type
LogSource).  It receives additions and removals of sources from other
components, and it informs subscribers of these additions and removals.

## [comp/logs/launchers/file](https://pkg.go.dev/github.com/djmitche/dd-agent-comp-experiments/comp/logs/launchers/file@v0.0.2)

*Datadog Team*: agent-metrics-logs

Package file implements a launcher that responds to file sources by starting
tailers for the indicated files.

## [comp/logs/launchers/launchermgr](https://pkg.go.dev/github.com/djmitche/dd-agent-comp-experiments/comp/logs/launchers/launchermgr@v0.0.2)

*Datadog Team*: agent-metrics-logs

Package launchermgr implements a component managing logs-agent launchers.  It collects
the set of loaded launchers during start-up, and allows enumeration and retrieval
as necessary.

## [comp/trace/agent](https://pkg.go.dev/github.com/djmitche/dd-agent-comp-experiments/comp/trace/agent@v0.0.2)

*Datadog Team*: trace-agent

Package agent implements a component representing the trace agent.

## [comp/trace/internal/httpreceiver](https://pkg.go.dev/github.com/djmitche/dd-agent-comp-experiments/comp/trace/internal/httpreceiver@v0.0.2)

*Datadog Team*: trace-agent

Package httpreceiver listens for incoming spans via HTTP and submits them to
the APM agent pipeline.

## [comp/trace/internal/processor](https://pkg.go.dev/github.com/djmitche/dd-agent-comp-experiments/comp/trace/internal/processor@v0.0.2)

*Datadog Team*: trace-agent

Package processor handles processing spans for the trace agent.

## [comp/trace/internal/tracewriter](https://pkg.go.dev/github.com/djmitche/dd-agent-comp-experiments/comp/trace/internal/tracewriter@v0.0.2)

*Datadog Team*: trace-agent

Package tracewriter buffers traces and APM events, flushing them to the
Datadog API.

## [comp/util/log](https://pkg.go.dev/github.com/djmitche/dd-agent-comp-experiments/comp/util/log@v0.0.2)

*Datadog Team*: agent-shared-components

Package log implements a component to handle logging internal to the agent.
