package handlers

import (
	"pets/internal/service"
	"pets/pkg/logger"
)

// Handlers is a server handlers struct
type Handlers struct {
	srv service.IService
}

// NewHandlers return new Handlers instance
func NewHandlers(srv service.IService) *Handlers {
	h := &Handlers{}

	h.srv = srv

	logger.Log().WithField("layer", "Handlers").Infof("handlers created")

	return h
}
