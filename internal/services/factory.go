package services

import (
	"fmt"
	"log"
)

// ServiceFactory is a function that creates a service instance
type ServiceFactory func(nodeID string) *Service

// ServiceRegistry holds all available service factories
type ServiceRegistry struct {
	factories map[string]ServiceFactory
}

// NewServiceRegistry creates a new service registry
func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		factories: make(map[string]ServiceFactory),
	}
}

// RegisterServiceFactory registers a service factory
func (sr *ServiceRegistry) RegisterServiceFactory(name string, factory ServiceFactory) {
	sr.factories[name] = factory
	log.Printf("Service factory registered: %s", name)
}

// CreateService creates a service instance by name
func (sr *ServiceRegistry) CreateService(name string, nodeID string) (*Service, error) {
	factory, exists := sr.factories[name]
	if !exists {
		return nil, fmt.Errorf("service factory not found: %s", name)
	}

	return factory(nodeID), nil
}

// ListAvailableServices returns all available service names
func (sr *ServiceRegistry) ListAvailableServices() []string {
	names := make([]string, 0, len(sr.factories))
	for name := range sr.factories {
		names = append(names, name)
	}
	return names
}

// CreateAllServices creates instances of all registered services
func (sr *ServiceRegistry) CreateAllServices(nodeID string) ([]*Service, error) {
	services := make([]*Service, 0, len(sr.factories))

	for name, factory := range sr.factories {
		service := factory(nodeID)
		services = append(services, service)
		log.Printf("Created service: %s", name)
	}

	return services, nil
}

// Global service registry instance
var GlobalServiceRegistry = NewServiceRegistry()
