//go:build auth
// +build auth

package http

import (
	"github.com/jessicatarra/greenlight/internal/config"
	"github.com/jessicatarra/greenlight/ms/auth/internal/domain/mocks"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

func TestNewService(t *testing.T) {

	mockAppl := &mocks.Appl{}
	mockCfg := config.Config{}
	mockLogger := slog.Logger{}

	service := NewService(mockAppl, mockCfg, &mockLogger)

	assert.NotNil(t, service)
}
