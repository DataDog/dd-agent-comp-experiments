// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package testing

import (
	"os"
	"strings"
)

// InTestBinary is true if the running binary is a Go test.
var InTestBinary bool

func init() {
	// 'go test' generates binaries with a `.test` suffix and runs
	// them, so we can detect that this is a test binary by looking
	// at the filename.
	InTestBinary = strings.HasSuffix(os.Args[0], ".test")
}
