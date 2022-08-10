# Agent Components

This file lists all components defined in this repository, with their package summary.
Click the links for more documentation.

## [comp/autodiscovery](https://pkg.go.dev/github.com/DataDog/dd-agent-comp-experiments/comp/autodiscovery) (Component Bundle)

*Datadog Team*: container-integrations

Package autodiscovery implements autodiscovery: collecting integration configuration
and monitoring for updates.

### [comp/autodiscovery/scheduler](https://pkg.go.dev/github.com/DataDog/dd-agent-comp-experiments/comp/autodiscovery/scheduler)

Package scheduler broadcasts changes to discovered configuration
configuration to its subscribers.  Subscriptions are created by providing a
Subscription value in value-group "autodiscovery".

## [comp/core](https://pkg.go.dev/github.com/DataDog/dd-agent-comp-experiments/comp/core) (Component Bundle)

*Datadog Team*: agent-shared-components

Package core implements the "core" bundle, providing services common to all
agent flavors and binaries.

### [comp/core/config](https://pkg.go.dev/github.com/DataDog/dd-agent-comp-experiments/comp/core/config)

Package config implements a component to handle agent configuration.  This
component wraps Viper.

### [comp/core/flare](https://pkg.go.dev/github.com/DataDog/dd-agent-comp-experiments/comp/core/flare)

Package flare implements a component creates flares for submission to support.

### [comp/core/health](https://pkg.go.dev/github.com/DataDog/dd-agent-comp-experiments/comp/core/health)

Package health implements a component that monitors the health of other
components.

### [comp/core/ipc/ipcclient](https://pkg.go.dev/github.com/DataDog/dd-agent-comp-experiments/comp/core/ipc/ipcclient)

Package ipcclient implements a component to access the IPC server remotely.

### [comp/core/ipc/ipcserver](https://pkg.go.dev/github.com/DataDog/dd-agent-comp-experiments/comp/core/ipc/ipcserver)

Package ipcserver implements a component to manage the IPC API server and act
as a client.

### [comp/core/log](https://pkg.go.dev/github.com/DataDog/dd-agent-comp-experiments/comp/core/log)

Package log implements a component to handle logging internal to the agent.

### [comp/core/status](https://pkg.go.dev/github.com/DataDog/dd-agent-comp-experiments/comp/core/status)

Package status implements the functionality behind `agent status`.

## [comp/logs](https://pkg.go.dev/github.com/DataDog/dd-agent-comp-experiments/comp/logs) (Component Bundle)

*Datadog Team*: agent-metrics-logs

Package logs collects the packages related to the logs agent.

### [comp/logs/agent](https://pkg.go.dev/github.com/DataDog/dd-agent-comp-experiments/comp/logs/agent)

Package agent implements a component representing the logs agent.  This
component coordinates activity related to gathering logs and forwarding them
to the Datadog intake.

### [comp/logs/internal/sourcemgr](https://pkg.go.dev/github.com/DataDog/dd-agent-comp-experiments/comp/logs/internal/sourcemgr)

Package sourcemgr implements a component managing logs-agent sources (type
LogSource).  It receives additions and removals of sources from other
components, and it informs subscribers of these additions and removals.

### [comp/logs/launchers/file](https://pkg.go.dev/github.com/DataDog/dd-agent-comp-experiments/comp/logs/launchers/file)

Package file implements a launcher that responds to file sources by starting
tailers for the indicated files.

### [comp/logs/launchers/launchermgr](https://pkg.go.dev/github.com/DataDog/dd-agent-comp-experiments/comp/logs/launchers/launchermgr)

Package launchermgr implements a component managing logs-agent launchers.  It collects
the set of loaded launchers during start-up, and allows enumeration and retrieval
as necessary.

## [comp/trace](https://pkg.go.dev/github.com/DataDog/dd-agent-comp-experiments/comp/trace) (Component Bundle)

*Datadog Team*: trace-agent

Package logs collects the packages related to the logs agent.

### [comp/trace/agent](https://pkg.go.dev/github.com/DataDog/dd-agent-comp-experiments/comp/trace/agent)

Package agent implements a component representing the trace agent.

### [comp/trace/internal/httpreceiver](https://pkg.go.dev/github.com/DataDog/dd-agent-comp-experiments/comp/trace/internal/httpreceiver)

Package httpreceiver listens for incoming spans via HTTP and submits them to
the APM agent pipeline.

### [comp/trace/internal/processor](https://pkg.go.dev/github.com/DataDog/dd-agent-comp-experiments/comp/trace/internal/processor)

Package processor handles processing spans for the trace agent.

### [comp/trace/internal/tracewriter](https://pkg.go.dev/github.com/DataDog/dd-agent-comp-experiments/comp/trace/internal/tracewriter)

Package tracewriter buffers traces and APM events, flushing them to the
Datadog API.
