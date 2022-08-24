// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package flare implements the `agent flare` command.
package flare

import (
	"errors"
	"fmt"

	"github.com/DataDog/dd-agent-comp-experiments/cmd/agent/root"
	"github.com/DataDog/dd-agent-comp-experiments/cmd/common"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/flare"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/ipc/ipcclient"
	"github.com/DataDog/dd-agent-comp-experiments/pkg/util/fxapps"
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:   "flare",
		Short: "Get a flare from the Agent",
		RunE:  command,
	}
)

func command(_ *cobra.Command, args []string) error {
	return fxapps.OneShot(flareCmd,
		common.SharedOptions(root.ConfFilePath, true),
	)
}

func getFlareRemote(ipcclient ipcclient.Component) (string, error) {
	var content map[string]string
	err := ipcclient.GetJSON("/agent/flare", &content)
	if err != nil {
		return "", err
	}
	if msg, found := content["error"]; found {
		return "", fmt.Errorf("Error from Agent: %s", msg)
	}

	if filename, found := content["filename"]; found {
		return filename, nil
	}
	return "", errors.New("No filename received from Agent")
}

func flareCmd(ipcclient ipcclient.Component, flare flare.Component) error {
	archiveFile, err := getFlareRemote(ipcclient)
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
