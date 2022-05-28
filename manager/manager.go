// Package manager implements various reset "managers" which handle Minecraft
// instances and their changing states.
package manager

import (
	"resetti/cfg"
	"resetti/mc"
	"resetti/x11"

	obs "github.com/woofdoggo/go-obs"
)

type Manager interface {
	SetConfig(cfg.Config)
	SetDeps(*x11.Client, *obs.Client)

	Restart([]mc.Instance) error
	Start([]mc.Instance, chan error) error
	Stop()
}
