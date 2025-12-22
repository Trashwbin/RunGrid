//go:build !windows

package main

import (
	"errors"

	"rungrid/backend/domain"
)

func (a *App) GetCursorAnchorPosition(width, height int) (domain.Point, error) {
	return domain.Point{}, errors.New("cursor position not supported")
}
