package errors

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type (
	// Params is used to replace placeholders in an error template with the corresponding values
	Params map[string]interface{}

	errorTemplate struct {
		Message string `yaml:"message"`
	}
)

var templates map[string]errorTemplate

// LoadErrors load a error_templates.yaml file containing error templates
func LoadErrors(file string) error {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	templates = map[string]errorTemplate{}
	return yaml.Unmarshal(bytes, &templates)
}

// NewAPIError creates a new APIError with the given HTTP status code, error code, and parameters
func NewAPIError(status int, code string, params Params) *APIError {
	err := &APIError{
		Status:    status,
		ErrorCode: code,
		Message:   code,
	}

	// Check that a template exist for this kind or HTTP error
	// If it does, values are going to be change according to the template
	if template, ok := templates[code]; ok {
		message := template.Message

		for key, value := range params {
			message = strings.Replace(message, "{"+key+"}", fmt.Sprint(value), -1)
		}
		err.Message = message
	}

	return err
}
