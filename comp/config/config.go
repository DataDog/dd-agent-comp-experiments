// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package config

import (
	"strings"

	"github.com/spf13/viper"
)

// config implements the Component.
type config struct {
	viper *viper.Viper
}

func newConfig(params ModuleParams) (Component, error) {
	v := viper.New()
	v.SetConfigName("datadog")
	v.SetEnvPrefix("DD_")
	v.SetConfigType("yaml")
	if params.ConfFilePath != "" {
		v.AddConfigPath(params.ConfFilePath)
		if strings.HasSuffix(params.ConfFilePath, ".yaml") {
			v.SetConfigFile(params.ConfFilePath)
		}
	}
	v.AddConfigPath("/etc/datadog-agent")

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	// for testing, this is good enough
	return &config{
		viper: v,
	}, nil
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
