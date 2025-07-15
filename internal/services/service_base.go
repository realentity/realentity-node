package services

import (
	"fmt"
	"sync"
)

// ServiceProvider defines the interface that all services must implement
type ServiceProvider interface {
	// Basic service identification
	GetName() string
	GetVersion() string
	GetDescription() string

	// Service lifecycle
	Start() error
	Stop() error
	IsEnabled() bool

	// Service execution
	Execute(request ServiceRequest) (*ServiceResponse, error)

	// Configuration
	GetConfig() map[string]interface{}
	SetConfig(config map[string]interface{}) error

	// Health check
	HealthCheck() error
}

// BaseService provides a common implementation for services
type BaseService struct {
	name        string
	version     string
	description string
	enabled     bool
	config      map[string]interface{}
	mutex       sync.RWMutex
}

// NewBaseService creates a new base service
func NewBaseService(name, version, description string) *BaseService {
	return &BaseService{
		name:        name,
		version:     version,
		description: description,
		enabled:     false,
		config:      make(map[string]interface{}),
	}
}

// GetName returns the service name
func (bs *BaseService) GetName() string {
	return bs.name
}

// GetVersion returns the service version
func (bs *BaseService) GetVersion() string {
	return bs.version
}

// GetDescription returns the service description
func (bs *BaseService) GetDescription() string {
	return bs.description
}

// Start enables the service
func (bs *BaseService) Start() error {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()

	if bs.enabled {
		return fmt.Errorf("service %s is already running", bs.name)
	}

	bs.enabled = true
	return nil
}

// Stop disables the service
func (bs *BaseService) Stop() error {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()

	if !bs.enabled {
		return fmt.Errorf("service %s is not running", bs.name)
	}

	bs.enabled = false
	return nil
}

// IsEnabled returns whether the service is running
func (bs *BaseService) IsEnabled() bool {
	bs.mutex.RLock()
	defer bs.mutex.RUnlock()
	return bs.enabled
}

// GetConfig returns the service configuration
func (bs *BaseService) GetConfig() map[string]interface{} {
	bs.mutex.RLock()
	defer bs.mutex.RUnlock()

	// Return a copy to avoid concurrent modification
	config := make(map[string]interface{})
	for k, v := range bs.config {
		config[k] = v
	}
	return config
}

// SetConfig updates the service configuration
func (bs *BaseService) SetConfig(config map[string]interface{}) error {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()

	// Clear existing config and set new values
	bs.config = make(map[string]interface{})
	for k, v := range config {
		bs.config[k] = v
	}

	return nil
}

// HealthCheck performs a basic health check
func (bs *BaseService) HealthCheck() error {
	bs.mutex.RLock()
	defer bs.mutex.RUnlock()

	if !bs.enabled {
		return fmt.Errorf("service %s is not enabled", bs.name)
	}

	return nil
}

// Execute must be implemented by concrete services
func (bs *BaseService) Execute(request ServiceRequest) (*ServiceResponse, error) {
	return nil, fmt.Errorf("execute method not implemented for service %s", bs.name)
}
