//go:build !windows

package hotkey

import (
	"context"

	"rungrid/backend/domain"
)

type Manager struct{}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) Start(_ context.Context) {}

func (m *Manager) Stop() {}

func (m *Manager) Apply(_ []domain.HotkeyBinding) ([]domain.HotkeyIssue, error) {
	return nil, ErrUnsupported
}
