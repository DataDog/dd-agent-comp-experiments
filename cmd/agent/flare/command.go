// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package flare implements the `agent flare` command.
package flare

import (
	"fmt"

	"github.com/djmitche/dd-agent-comp-experiments/cmd/agent/root"
	"github.com/djmitche/dd-agent-comp-experiments/cmd/common"
	"github.com/djmitche/dd-agent-comp-experiments/comp/core/flare"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var (
	Cmd = &cobra.Command{
		Use:   "flare",
		Short: "Get a flare from the Agent",
		RunE:  command,
	}
)

func command(_ *cobra.Command, args []string) error {
	app := fx.New(
		common.SharedOptions(root.ConfFilePath, true),
		common.OneShot(flareCmd),
	)
	return common.RunApp(app)
}

func flareCmd(flare flare.Component) error {
	archiveFile, err := flare.CreateFlareRemote()
	if err != nil {
		fmt.Printf("Could not contact agent: %s\n", err)
		fmt.Printf("Proceeding with local flare.\n")
		archiveFile, err = flare.CreateFlare()
	}
	if err != nil {
		return err
	}

	fmt.Printf("Generated flare file %s\n", archiveFile)
	return nil
}
