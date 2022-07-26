// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package health implements the `agent health` command.
package health

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/djmitche/dd-agent-comp-experiments/cmd/agent/root"
	"github.com/djmitche/dd-agent-comp-experiments/cmd/common"
	"github.com/djmitche/dd-agent-comp-experiments/comp/config"
	"github.com/djmitche/dd-agent-comp-experiments/comp/logs/launchers/file"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var (
	Cmd = &cobra.Command{
		Use:   "health",
		Short: "Get the health of Agent's components",
		RunE:  command,
	}
)

func logsAgentPluginOptions() fx.Option {
	return fx.Options(
		// this list would be different for other agent flavors
		file.Module,
		fx.Invoke(func(file.Component) {}),
	)
}

func command(_ *cobra.Command, args []string) error {
	app := fx.New(
		common.SharedOptions(root.ConfFilePath, true),
		common.OneShot(health),
	)
	return common.RunApp(app)
}

func health(config config.Component) error {
	port := config.GetInt("cmd_port")
	url := fmt.Sprintf("http://127.0.0.1:%d/agent/health", port)
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	health := map[string]map[string]interface{}{}
	err = json.Unmarshal(body, &health)
	if err != nil {
		return err
	}

	for component, h := range health {
		fmt.Printf("%s: ", component)
		if h["Healthy"].(bool) {
			fmt.Printf("OK\n")
		} else {
			fmt.Printf("UNHEALTHY (%s)\n", h["Message"].(string))
		}
	}

	return nil
}
