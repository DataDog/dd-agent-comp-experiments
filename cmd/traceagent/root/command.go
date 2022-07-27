// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package root defines the root 'agent' command.
package root

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// ConfFilePath holds the path to the folder containing the configuration
	// file, to allow overrides from the command line
	ConfFilePath string
)

func MakeCommand(subcommands ...*cobra.Command) *cobra.Command {
	// AgentCmd is the root command
	agentCmd := &cobra.Command{
		Use:          fmt.Sprintf("%s [command]", os.Args[0]),
		Short:        "Datadog Agent at your service, just for tracing.",
		SilenceUsage: true,
	}

	agentCmd.PersistentFlags().StringVarP(&ConfFilePath, "cfgpath", "c", "", "path to directory containing datadog.yaml")

	for _, sub := range subcommands {
		agentCmd.AddCommand(sub)
	}

	return agentCmd
}
