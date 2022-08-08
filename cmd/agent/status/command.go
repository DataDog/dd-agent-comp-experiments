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
	"github.com/djmitche/dd-agent-comp-experiments/comp/core/ipc/ipcclient"
	"github.com/djmitche/dd-agent-comp-experiments/comp/core/status"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var (
	Cmd = &cobra.Command{
		Use:   "status [section]",
		Short: "Get the running Agent's status, optionally showing only a single section",
		RunE:  command,
		Args:  cobra.MaximumNArgs(1),
	}
)

func command(_ *cobra.Command, args []string) error {
	var section string
	if len(args) > 0 {
		section = args[0]
	}

	app := fx.New(
		common.SharedOptions(root.ConfFilePath, true),
		common.OneShot(func(ipcclient ipcclient.Component, status status.Component) error {
			return statusCmd(ipcclient, status, section)
		}),
	)
	return common.RunApp(app)
}

func getStatusRemote(ipcclient ipcclient.Component, section string) (string, error) {
	var content map[string]string
	path := "/agent/status"
	if section != "" {
		path = fmt.Sprintf("%s?section=%s", path, section)
	}

	err := ipcclient.GetJSON(path, &content)
	if err != nil {
		return "", err
	}

	return content["status"], nil
}

func statusCmd(ipcclient ipcclient.Component, status status.Component, section string) error {
	statusStr, err := getStatusRemote(ipcclient, section)
	if err != nil {
		return err
	}

	fmt.Printf("%s", statusStr)
	return nil
}
