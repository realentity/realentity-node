package services

import (
	"fmt"
	"log"
	"sync"
)

// ServiceManager handles the lifecycle and management of services
type ServiceManager struct {
	services map[string]ServiceProvider
	mutex    sync.RWMutex
}

// NewServiceManager creates a new service manager
func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		services: make(map[string]ServiceProvider),
	}
}

// RegisterService registers a service with the manager
func (sm *ServiceManager) RegisterService(service ServiceProvider) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	name := service.GetName()
	if _, exists := sm.services[name]; exists {
		return fmt.Errorf("service %s already registered", name)
	}

	sm.services[name] = service
	log.Printf("Service registered: %s (version %s)", name, service.GetVersion())
	return nil
}

// UnregisterService removes a service from the manager
func (sm *ServiceManager) UnregisterService(name string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	service, exists := sm.services[name]
	if !exists {
		return fmt.Errorf("service %s not found", name)
	}

	// Stop the service if it's running
	if service.IsEnabled() {
		if err := service.Stop(); err != nil {
			log.Printf("Error stopping service %s: %v", name, err)
		}
	}

	delete(sm.services, name)
	log.Printf("Service unregistered: %s", name)
	return nil
}

// GetService retrieves a service by name
func (sm *ServiceManager) GetService(name string) (ServiceProvider, error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	service, exists := sm.services[name]
	if !exists {
		return nil, fmt.Errorf("service %s not found", name)
	}

	return service, nil
}

// ListServices returns a list of all registered service names
func (sm *ServiceManager) ListServices() []string {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	names := make([]string, 0, len(sm.services))
	for name := range sm.services {
		names = append(names, name)
	}

	return names
}

// StartService starts a specific service
func (sm *ServiceManager) StartService(name string) error {
	sm.mutex.RLock()
	service, exists := sm.services[name]
	sm.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("service %s not found", name)
	}

	return service.Start()
}

// StopService stops a specific service
func (sm *ServiceManager) StopService(name string) error {
	sm.mutex.RLock()
	service, exists := sm.services[name]
	sm.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("service %s not found", name)
	}

	return service.Stop()
}

// StartAllServices starts all registered services
func (sm *ServiceManager) StartAllServices() error {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	var errors []string
	for name, service := range sm.services {
		if !service.IsEnabled() {
			if err := service.Start(); err != nil {
				errors = append(errors, fmt.Sprintf("%s: %v", name, err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to start some services: %v", errors)
	}

	return nil
}

// StopAllServices stops all running services
func (sm *ServiceManager) StopAllServices() error {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	var errors []string
	for name, service := range sm.services {
		if service.IsEnabled() {
			if err := service.Stop(); err != nil {
				errors = append(errors, fmt.Sprintf("%s: %v", name, err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to stop some services: %v", errors)
	}

	return nil
}

// GetServiceInfo returns information about all services
func (sm *ServiceManager) GetServiceInfo() map[string]interface{} {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	info := make(map[string]interface{})
	for name, service := range sm.services {
		info[name] = map[string]interface{}{
			"enabled":     service.IsEnabled(),
			"version":     service.GetVersion(),
			"description": service.GetDescription(),
			"config":      service.GetConfig(),
		}
	}

	return info
}

// GetServiceStatus returns status information for a specific service
func (sm *ServiceManager) GetServiceStatus(name string) (map[string]interface{}, error) {
	sm.mutex.RLock()
	service, exists := sm.services[name]
	sm.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("service %s not found", name)
	}

	return map[string]interface{}{
		"name":        service.GetName(),
		"enabled":     service.IsEnabled(),
		"version":     service.GetVersion(),
		"description": service.GetDescription(),
		"config":      service.GetConfig(),
	}, nil
}

// ExecuteService executes a service with the given request
func (sm *ServiceManager) ExecuteService(serviceName string, request ServiceRequest) (*ServiceResponse, error) {
	service, err := sm.GetService(serviceName)
	if err != nil {
		return nil, err
	}

	if !service.IsEnabled() {
		return nil, fmt.Errorf("service %s is not enabled", serviceName)
	}

	return service.Execute(request)
}

// Global service manager instance
var GlobalServiceManager = NewServiceManager()
