package impl

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/realentity/realentity-node/internal/services"
)

// MathRequest represents the payload for math operations
type MathRequest struct {
	Operation string    `json:"operation"` // "add", "subtract", "multiply", "divide", "sqrt", "power"
	Numbers   []float64 `json:"numbers"`   // Numbers to operate on
}

// MathResponse represents the response from math operations
type MathResponse struct {
	Operation string    `json:"operation"`
	Numbers   []float64 `json:"numbers"`
	Result    float64   `json:"result"`
}

// CreateMathService creates a mathematical operations service
func CreateMathService(nodeID string) *services.Service {
	return &services.Service{
		Name:        "math",
		Description: "Mathematical operations service",
		Version:     "1.0.0",
		Metadata: map[string]string{
			"category":   "math",
			"operations": "add,subtract,multiply,divide,sqrt,power",
			"cost":       "free",
		},
		Handler: func(payload []byte) ([]byte, error) {
			var req MathRequest
			if err := json.Unmarshal(payload, &req); err != nil {
				return nil, fmt.Errorf("invalid math request: %v", err)
			}

			log.Printf("math service executing: operation='%s', numbers=%v", req.Operation, req.Numbers)

			if len(req.Numbers) == 0 {
				return nil, fmt.Errorf("no numbers provided")
			}

			var result float64
			switch strings.ToLower(req.Operation) {
			case "add":
				result = 0
				for _, num := range req.Numbers {
					result += num
				}
			case "subtract":
				if len(req.Numbers) < 2 {
					return nil, fmt.Errorf("subtract requires at least 2 numbers")
				}
				result = req.Numbers[0]
				for i := 1; i < len(req.Numbers); i++ {
					result -= req.Numbers[i]
				}
			case "multiply":
				result = 1
				for _, num := range req.Numbers {
					result *= num
				}
			case "divide":
				if len(req.Numbers) < 2 {
					return nil, fmt.Errorf("divide requires at least 2 numbers")
				}
				result = req.Numbers[0]
				for i := 1; i < len(req.Numbers); i++ {
					if req.Numbers[i] == 0 {
						return nil, fmt.Errorf("division by zero")
					}
					result /= req.Numbers[i]
				}
			case "sqrt":
				if len(req.Numbers) != 1 {
					return nil, fmt.Errorf("sqrt requires exactly 1 number")
				}
				if req.Numbers[0] < 0 {
					return nil, fmt.Errorf("sqrt of negative number")
				}
				result = math.Sqrt(req.Numbers[0])
			case "power":
				if len(req.Numbers) != 2 {
					return nil, fmt.Errorf("power requires exactly 2 numbers (base, exponent)")
				}
				result = math.Pow(req.Numbers[0], req.Numbers[1])
			default:
				return nil, fmt.Errorf("unsupported operation: %s", req.Operation)
			}

			response := MathResponse{
				Operation: req.Operation,
				Numbers:   req.Numbers,
				Result:    result,
			}

			resultBytes, err := json.Marshal(response)
			if err != nil {
				return nil, err
			}

			resultStr := strconv.FormatFloat(result, 'f', -1, 64)
			log.Printf("math service completed: result='%s'", resultStr)
			return resultBytes, nil
		},
	}
}

// Register this service with the global registry
func init() {
	services.GlobalServiceRegistry.RegisterServiceFactory("math", CreateMathService)
}
