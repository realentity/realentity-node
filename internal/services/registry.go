package services

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Service represents a service that this node can provide
type Service struct {
	Name        string            `json:"name"`        // e.g., "ai.infer", "resize.image"
	Description string            `json:"description"` // Human-readable description
	Version     string            `json:"version"`     // Service version
	Metadata    map[string]string `json:"metadata"`    // Additional service info
	Handler     ServiceHandler    `json:"-"`           // Function to execute the service
}

// ServiceHandler defines the interface for service execution
type ServiceHandler func(payload []byte) ([]byte, error)

// ServiceRequest represents an incoming service request
type ServiceRequest struct {
	Service   string          `json:"service"`   // Service name to execute
	Payload   json.RawMessage `json:"payload"`   // Service-specific data
	RequestID string          `json:"requestId"` // Unique request identifier
}

// ServiceResponse represents a service execution response
type ServiceResponse struct {
	RequestID string          `json:"requestId"`
	Success   bool            `json:"success"`
	Result    json.RawMessage `json:"result,omitempty"`
	Error     string          `json:"error,omitempty"`
}

// Registry manages local services for this node
type Registry struct {
	services map[string]*Service
	mutex    sync.RWMutex
}

// NewRegistry creates a new service registry
func NewRegistry() *Registry {
	return &Registry{
		services: make(map[string]*Service),
	}
}

// RegisterService adds a service to the registry
func (r *Registry) RegisterService(service *Service) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if service.Name == "" {
		return fmt.Errorf("service name cannot be empty")
	}

	if service.Handler == nil {
		return fmt.Errorf("service handler cannot be nil")
	}

	r.services[service.Name] = service
	return nil
}

// GetService retrieves a service by name
func (r *Registry) GetService(name string) (*Service, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	service, exists := r.services[name]
	return service, exists
}

// GetAllServices returns all registered services
func (r *Registry) GetAllServices() map[string]*Service {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	services := make(map[string]*Service)
	for name, service := range r.services {
		services[name] = service
	}
	return services
}

// ExecuteService runs a service with the given payload
func (r *Registry) ExecuteService(request *ServiceRequest) *ServiceResponse {
	service, exists := r.GetService(request.Service)
	if !exists {
		return &ServiceResponse{
			RequestID: request.RequestID,
			Success:   false,
			Error:     fmt.Sprintf("service '%s' not found", request.Service),
		}
	}

	result, err := service.Handler(request.Payload)
	if err != nil {
		return &ServiceResponse{
			RequestID: request.RequestID,
			Success:   false,
			Error:     err.Error(),
		}
	}

	return &ServiceResponse{
		RequestID: request.RequestID,
		Success:   true,
		Result:    json.RawMessage(result),
	}
}

// ListServices returns a list of service names
func (r *Registry) ListServices() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	names := make([]string, 0, len(r.services))
	for name := range r.services {
		names = append(names, name)
	}
	return names
}

// Global registry instance
var GlobalRegistry = NewRegistry()
