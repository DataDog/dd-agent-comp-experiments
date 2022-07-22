// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package config

// config implements the Component.
//
// XXX: In a real agent, this would use Viper much like pkg/config.
type config struct {
	values map[string]interface{}
}

func newConfig(configFilePath string) (Component, error) {
	// for testing, this is good enough
	return &config{
		values: map[string]interface{}{
			"logs_config.container_collect_all":     true,
			"logs_config.docker_container_use_file": true,
			"logs_config.k8s_container_use_file":    true,
		},
	}, nil
}

// GetInt implements Component#GetInt.
func (c *config) GetInt(key string) int {
	return c.values[key].(int)
}

// GetBool implements Component#GetBool.
func (c *config) GetBool(key string) bool {
	return c.values[key].(bool)
}

// GetString implements Component#GetString.
func (c *config) GetString(key string) string {
	return c.values[key].(string)
}
