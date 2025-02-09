package ctl

import (
	"fmt"

	"github.com/woofdoggo/resetti/internal/cfg"
	"github.com/woofdoggo/resetti/internal/mc"
	"github.com/woofdoggo/resetti/internal/obs"
	"github.com/woofdoggo/resetti/internal/x11"
)

// Multi implements a traditional Multi-instance interface, where the user
// plays and resets one instance at a time.
type Multi struct {
	host *Controller
	conf *cfg.Profile
	obs  *obs.Client
	x    *x11.Client

	instances []mc.InstanceInfo
	states    []mc.State
	active    int
}

// Setup implements Frontend.
func (m *Multi) Setup(deps frontendDependencies) error {
	m.host = deps.host
	m.conf = deps.conf
	m.obs = deps.obs
	m.x = deps.x

	m.active = 0
	m.instances = make([]mc.InstanceInfo, len(deps.instances))
	m.states = make([]mc.State, len(deps.states))
	copy(m.instances, deps.instances)
	copy(m.states, deps.states)

	m.host.FocusInstance(0)
	return nil
}

// Input implements Frontend.
func (m *Multi) Input(input Input) {
	actions := m.conf.Keybinds[input.Bind]
	if input.Held {
		return
	}
	for _, action := range actions.IngameActions {
		switch action.Type {
		case cfg.ActionIngameFocus:
			m.host.FocusInstance(m.active)
		case cfg.ActionIngameReset:
			if m.x.GetActiveWindow() != m.instances[m.active].Wid {
				continue
			}
			next := (m.active + 1) % len(m.states)
			current := m.active
			if m.host.ResetInstance(current) {
				m.host.PlayInstance(next)
				m.active = next
				m.updateObs()
				m.host.RunHook(HookReset)
			}
		case cfg.ActionIngameRes:
			m.host.ToggleResolution(m.active)
		}
	}
}

// ProcessEvent implements Frontend.
func (m *Multi) ProcessEvent(x11.Event) {
	// Do nothing.
}

// Update implements Frontend.
func (m *Multi) Update(update mc.Update) {
	m.states[update.Id] = update.State
}

// updateObs changes which instance is visible on the OBS scene.
func (m *Multi) updateObs() {
	if !m.conf.Obs.Enabled {
		return
	}
	m.obs.BatchAsync(obs.SerialRealtime, func(b *obs.Batch) {
		for i := 1; i <= len(m.states); i += 1 {
			b.SetItemVisibility("Instance", fmt.Sprintf("MC %d", i), i-1 == m.active)
		}
	})
}
