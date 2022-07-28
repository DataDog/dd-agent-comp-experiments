// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package status implements the `agent status` command.
package status

import (
	"fmt"

	"github.com/djmitche/dd-agent-comp-experiments/cmd/agent/root"
	"github.com/djmitche/dd-agent-comp-experiments/cmd/common"
	"github.com/djmitche/dd-agent-comp-experiments/comp/status"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var (
	Cmd = &cobra.Command{
		Use:   "status",
		Short: "Get the running Agent's status",
		RunE:  command,
	}
)

func command(_ *cobra.Command, args []string) error {
	app := fx.New(
		common.SharedOptions(root.ConfFilePath, true),
		common.OneShot(statusCmd),
	)
	return common.RunApp(app)
}

func statusCmd(status status.Component) error {
	statusStr, err := status.GetStatusRemote("")
	if err != nil {
		return err
	}

	fmt.Printf("%s", statusStr)
	return nil
}
