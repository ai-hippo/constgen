package parser

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/ai-hippo/constgen/internal/model"
)

// ParseYAML loads and validates a YAML spec.
func ParseYAML(path string) (*model.Spec, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var spec model.Spec

	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, err
	}

	if err := ValidateSpec(&spec); err != nil {
		return nil, err
	}

	return &spec, nil
}