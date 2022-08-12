// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package status implements the `agent status` command.
package status

import (
	"fmt"

	"github.com/DataDog/dd-agent-comp-experiments/cmd/agent/root"
	"github.com/DataDog/dd-agent-comp-experiments/cmd/common"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/ipc/ipcclient"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/status"
	"github.com/DataDog/dd-agent-comp-experiments/pkg/util/fxapps"
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

type cmdArgs struct {
	section string
}

func command(_ *cobra.Command, args []string) error {
	var cmdArgs cmdArgs
	if len(args) > 0 {
		cmdArgs.section = args[0]
	}

	return fxapps.OneShot(statusCmd,
		fx.Supply(cmdArgs),
		common.SharedOptions(root.ConfFilePath, true),
	)
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

func statusCmd(ipcclient ipcclient.Component, status status.Component, cmdArgs cmdArgs) error {
	statusStr, err := getStatusRemote(ipcclient, cmdArgs.section)
	if err != nil {
		return err
	}

	fmt.Printf("%s", statusStr)
	return nil
}
