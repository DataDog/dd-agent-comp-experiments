// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package agent

import (
	configPkg "github.com/djmitche/dd-agent-comp-experiments/comp/config"
)

type config struct {
	dockerContainerUseFile bool
	k8sContainerUseFile    bool
	containerCollectAll    bool
}

// newConfig generates a *config from a configPkg.Component.

// TODO: this is very similar to a Redux "reducer": mapping a large data
// structure to just the data of interest.  With a bit of reflection we
// could build general support for this, using struct tags:
//
//    dockerContainerUseFile bool `config:"logs_config.docker_container_use_file"`
//
// Strict adherence would _vastly_ simplify determining what components use
// which configuration, and how it's used.
//
// The Redux model is also a good one for config update: re-apply reducers, and
// those can generate further events when the reduced fields change.
func newConfig(cfg configPkg.Component) *config {
	return &config{
		dockerContainerUseFile: cfg.GetBool("logs_config.docker_container_use_file"),
		k8sContainerUseFile:    cfg.GetBool("logs_config.k8s_container_use_file"),
		containerCollectAll:    cfg.GetBool("logs_config.container_collect_all"),
	}
}
