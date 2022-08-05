// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package status

// Registration is provided by other components in order to register sections
// for status reporting.
type registration struct {
	// section is the name of the status section
	section string

	// order determines the order of the sections
	order int

	// cb generates the content of the section
	cb func() string
}

// byOrder supports sorting sections by order.
type byOrder []registration

func (a byOrder) Len() int           { return len(a) }
func (a byOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byOrder) Less(i, j int) bool { return a[i].order < a[j].order }
