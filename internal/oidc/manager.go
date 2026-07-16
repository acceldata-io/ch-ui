// SPDX-License-Identifier: BUSL-1.1
// Copyright (C) 2024-2026 Caio Ricciuti.
// Part of CH-UI Pro. Licensed under the Business Source License 1.1 (see
// LICENSE.BSL), NOT the Apache-2.0 LICENSE that governs the rest of the repo.

package oidc

import (
	"context"
	"fmt"
	"sync"
)

// Manager holds the currently active OIDC provider (if any) and lets the
// admin SSO settings endpoint swap it at runtime without a server restart.
// A nil provider means SSO is not configured / not active.
type Manager struct {
	mu       sync.RWMutex
	provider *Provider
}

// NewManager returns an empty (inactive) manager.
func NewManager() *Manager {
	return &Manager{}
}

// Provider returns the active provider, or nil when SSO is inactive.
func (m *Manager) Provider() *Provider {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.provider
}

// Active reports whether an initialized provider is available.
func (m *Manager) Active() bool {
	return m.Provider() != nil
}

// Reload runs issuer discovery for the given settings and, on success, swaps
// the active provider. On failure the previously active provider (if any)
// keeps serving, so a bad save never breaks working SSO.
func (m *Manager) Reload(ctx context.Context, s Settings) error {
	if !s.Enabled || !s.Complete() {
		return fmt.Errorf("oidc settings are incomplete")
	}
	p, err := New(ctx, s)
	if err != nil {
		return err
	}
	m.mu.Lock()
	m.provider = p
	m.mu.Unlock()
	return nil
}

// Disable deactivates SSO by dropping the active provider.
func (m *Manager) Disable() {
	m.mu.Lock()
	m.provider = nil
	m.mu.Unlock()
}
