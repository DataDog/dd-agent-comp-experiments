// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package config

import (
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// config implements the Component.
type config struct {
	viper *viper.Viper
}

type dependencies struct {
	fx.In

	Params ModuleParams `optional:"true"`
}

func newConfig(deps dependencies) (Component, error) {
	v := viper.New()
	v.SetConfigName("datadog")
	v.SetEnvPrefix("DD_")
	v.SetConfigType("yaml")
	if deps.Params.ConfFilePath != "" {
		v.AddConfigPath(deps.Params.ConfFilePath)
		if strings.HasSuffix(deps.Params.ConfFilePath, ".yaml") {
			v.SetConfigFile(deps.Params.ConfFilePath)
		}
	}
	v.AddConfigPath("/etc/datadog-agent")

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	c := &config{
		viper: v,
	}

	return c, nil
}

// GetInt implements Component#GetInt.
func (c *config) GetInt(key string) int {
	return c.viper.GetInt(key)
}

// GetBool implements Component#GetBool.
func (c *config) GetBool(key string) bool {
	return c.viper.GetBool(key)
}

// GetString implements Component#GetString.
func (c *config) GetString(key string) string {
	return c.viper.GetString(key)
}

// WriteConfig implements Component#WriteConfig.
func (c *config) WriteConfig(filename string) error {
	return c.viper.SafeWriteConfigAs(filename)
}
