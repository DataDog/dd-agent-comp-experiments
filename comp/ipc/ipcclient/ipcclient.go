// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package ipcclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/djmitche/dd-agent-comp-experiments/comp/config"
	"go.uber.org/fx"
)

type client struct {
	// port is the port on which the server is running.
	port int
}

type dependencies struct {
	fx.In

	Config config.Component
}

func newClient(deps dependencies) Component {
	a := &client{
		port: deps.Config.GetInt("cmd_port"),
	}
	return a
}

// GetJSON implements Component#GetJSON.
func (a *client) GetJSON(path string, v any) error {
	url := fmt.Sprintf("http://127.0.0.1:%d%s", a.port, path)
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error contacting Agent: %s", err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Error contacting Agent: %s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, v)

	if err != nil {
		return fmt.Errorf("Error decoding Agent response: %s", err)
	}

	return nil

}
