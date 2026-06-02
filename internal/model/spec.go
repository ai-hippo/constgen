package model

// Spec is the validated AST used by generators.
// At this stage, YAML has already been validated strictly.
type Spec struct {
	Name string `yaml:"name"`
	Profiles map[string]map[string]any `yaml:"profiles"`
	Targets map[string]TargetSpec `yaml:"targets"`
}

// TargetSpec defines which profiles a language should generate.
type TargetSpec struct {
	Enabled  bool `yaml:"enabled"`
	Profiles []string `yaml:"profiles"`
	// Optional language-specific function signature overrides.
	FunctionSignatures map[string]string `yaml:"function_signatures,omitempty"`
}