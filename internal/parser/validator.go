package parser

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/ai-hippo/constgen/internal/model"
)

var constantRegex = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// profile names:
var profileRegex = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

var functionNameRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// ValidateSpec validates YAML BEFORE generation begins.
func ValidateSpec(spec *model.Spec) error {

	// -----------------------------------
	// name validation
	// -----------------------------------

	if spec.Name == "" {
		return errors.New("name cannot be empty")
	}

	// IMPORTANT:
	// name is intentionally flexible.
	//
	// Example:
	// BooleanCodes
	// Fruits
	// CardTypes
	//
	// We DO NOT enforce uppercase rules here.


	// -----------------------------------
	// profiles validation
	// -----------------------------------

	if len(spec.Profiles) == 0 {
		return errors.New("at least one profile is required")
	}

	for profileName, constants := range spec.Profiles {

		// profile naming rule
		if !profileRegex.MatchString(profileName) {
			return fmt.Errorf(
				"invalid profile name: %s",
				profileName,
			)
		}

		// empty profile allowed
		//
		// Example:
		//
		// integer:
		//
		// means:
		// "profile exists but unused"
		//
		if len(constants) == 0 {
			continue
		}

		for key, value := range constants {

			// reserved key: functions
			if key == "functions" {
				fm, ok := value.(map[string]interface{})
				if !ok {
					return fmt.Errorf("invalid 'functions' definition in profile '%s'", profileName)
				}

				for fname, fval := range fm {
					if !functionNameRegex.MatchString(fname) {
						return fmt.Errorf("invalid function name '%s' in profile '%s'", fname, profileName)
					}

					switch arr := fval.(type) {
					case []interface{}:
						for _, item := range arr {
							s, ok := item.(string)
							if !ok || !constantRegex.MatchString(s) {
								return fmt.Errorf("invalid function constant reference in profile '%s' for function '%s'", profileName, fname)
							}
						}
					case string:
						if !constantRegex.MatchString(arr) {
							return fmt.Errorf("invalid function constant reference in profile '%s' for function '%s'", profileName, fname)
						}
					default:
						return fmt.Errorf("invalid function definition for '%s' in profile '%s'", fname, profileName)
					}
				}

				continue
			}

			// validate constant key format
			if !constantRegex.MatchString(key) {
				return fmt.Errorf(
					"invalid constant key: %s",
					key,
				)
			}

			// validate supported value types
			switch value.(type) {

			case string:
			case int:
			case int64:
			case float64:
			case bool:

			default:
				return fmt.Errorf(
					"unsupported value type for key '%s' in profile '%s'",
					key,
					profileName,
				)
			}
		}
	}


	// -----------------------------------
	// targets validation
	// -----------------------------------

	if len(spec.Targets) == 0 {
		return errors.New("at least one target is required")
	}

	for lang, target := range spec.Targets {

		if !target.Enabled {
			continue
		}

		if len(target.Profiles) == 0 {
			return fmt.Errorf(
				"target '%s' must define at least one profile",
				lang,
			)
		}

		for _, profile := range target.Profiles {

			if _, ok := spec.Profiles[profile]; !ok {
				return fmt.Errorf(
					"target '%s' references unknown profile '%s'",
					lang,
					profile,
				)
			}
		}
	}

	// validate optional function_signatures in targets
	for lang, target := range spec.Targets {
		if target.FunctionSignatures == nil {
			continue
		}

		for fname, sig := range target.FunctionSignatures {
			if !functionNameRegex.MatchString(fname) {
				return fmt.Errorf("invalid function name '%s' in target '%s' function_signatures", fname, lang)
			}
			if sig == "" {
				return fmt.Errorf("function signature for '%s' in target '%s' must be a non-empty string", fname, lang)
			}
		}
	}

	return nil
}