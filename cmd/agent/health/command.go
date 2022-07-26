// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package health implements the `agent health` command.
package health

import (
	"fmt"

	"github.com/DataDog/dd-agent-comp-experiments/cmd/agent/root"
	"github.com/DataDog/dd-agent-comp-experiments/cmd/common"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/health"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/ipc/ipcclient"
	"github.com/DataDog/dd-agent-comp-experiments/pkg/util/fxapps"
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:   "health",
		Short: "Get the health of Agent's components",
		RunE:  command,
	}
)

func command(_ *cobra.Command, args []string) error {
	return fxapps.OneShot(healthCmd,
		common.SharedOptions(root.ConfFilePath, true),
	)
}

func getHealthRemote(ipcclient ipcclient.Component) (map[string]health.ComponentHealth, error) {
	var content map[string]health.ComponentHealth
	err := ipcclient.GetJSON("/agent/health", &content)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func healthCmd(ipcclient ipcclient.Component, health health.Component) error {
	resp, err := getHealthRemote(ipcclient)
	if err != nil {
		return err
	}

	for component, h := range resp {
		fmt.Printf("%s: ", component)
		if h.Healthy {
			fmt.Printf("OK\n")
		} else {
			fmt.Printf("UNHEALTHY (%s)\n", h.Message)
		}
	}

	return nil
}
