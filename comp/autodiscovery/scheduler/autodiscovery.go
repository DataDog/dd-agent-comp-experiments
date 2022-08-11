// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package scheduler

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/DataDog/dd-agent-comp-experiments/comp/autodiscovery/internal"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/health"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/log"
	"github.com/DataDog/dd-agent-comp-experiments/pkg/util/actor"
	"github.com/DataDog/dd-agent-comp-experiments/pkg/util/subscriptions"
	"go.uber.org/fx"
)

type autoDiscovery struct {
	// Mutex covers all fields in this type.
	sync.Mutex

	// log is the log component
	log log.Component

	// configChangeTx connects to receivers of messages about config additions/removals
	configChangeTx subscriptions.Transmitter[ConfigChange]

	// actor manages the goroutine "monitoring" for container/pod changes
	actor actor.Actor

	health *health.Handle
}

type dependencies struct {
	fx.In

	Lc     fx.Lifecycle
	Log    log.Component
	Pub    subscriptions.Publisher[ConfigChange]
	Params internal.BundleParams
}

func newAD(deps dependencies) (Component, health.Registration) {
	healthReg := health.NewRegistration(componentName)
	ad := &autoDiscovery{
		log:            deps.Log,
		configChangeTx: deps.Pub.Transmitter(),
		health:         healthReg.Handle,
	}
	if deps.Params.ShouldStart() {
		ad.actor.HookLifecycle(deps.Lc, ad.run)
	}
	return ad, healthReg
}

func (ad *autoDiscovery) run(ctx context.Context) {
	monitor, stopMonitor := ad.health.LivenessMonitor(time.Second)
	scheduled := []*Config{}
	tkr := time.NewTicker(time.Second)
	for {
		select {
		case <-tkr.C:
			if len(scheduled) == 0 || rand.Intn(2) == 0 {
				cfg := &Config{Name: fmt.Sprintf("cfg-%d", rand.Int63())}
				scheduled = append(scheduled, cfg)
				ad.log.Debug("scheduling", cfg.Name)
				ad.configChangeTx.Notify(ConfigChange{IsScheduled: true, Config: cfg})
			} else {
				i := rand.Intn(len(scheduled))
				cfg := scheduled[i]
				scheduled = append(scheduled[:i], scheduled[i+1:]...)
				ad.log.Debug("unscheduling", cfg.Name)
				ad.configChangeTx.Notify(ConfigChange{IsScheduled: false, Config: cfg})
			}
		case <-monitor:
		case <-ctx.Done():
			stopMonitor()
			return
		}
	}
}
