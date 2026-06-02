package generator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ai-hippo/constgen/internal/model"
)

// GoGenerator emits Go constants with proper typing.
type GoGenerator struct{}

func (g GoGenerator) FileName(spec model.Spec, profile string) string {
	return fmt.Sprintf(
		"%s.go",
		spec.Name,
	)
}

// Generate renders Go constants and helper functions.
func (g GoGenerator) Generate(spec model.Spec, profile string, lang string) string {

	var b strings.Builder

	// fetch profile directly
	constants, ok := spec.Profiles[profile]
	if !ok {
		return ""
	}

	b.WriteString("package generated\n\n")
	b.WriteString("const (\n")

	// collect keys deterministically and skip reserved 'functions'
	var keys []string
	for k := range constants {
		if k == "functions" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		value := constants[key]

		switch v := value.(type) {

		// -----------------------------------
		// string constants
		// -----------------------------------

		case string:

			b.WriteString(
				fmt.Sprintf(
					"    %s = \"%s\"\n",
					key,
					v,
				),
			)

		// -----------------------------------
		// integer constants
		// -----------------------------------

		case int:

			b.WriteString(
				fmt.Sprintf(
					"    %s = %d\n",
					key,
					v,
				),
			)

		case int64:

			b.WriteString(
				fmt.Sprintf(
					"    %s = %d\n",
					key,
					v,
				),
			)

		// -----------------------------------
		// float constants
		// -----------------------------------

		case float64:

			b.WriteString(
				fmt.Sprintf(
					"    %s = %v\n",
					key,
					v,
				),
			)

		// -----------------------------------
		// boolean constants
		// -----------------------------------

		case bool:

			b.WriteString(
				fmt.Sprintf(
					"    %s = %t\n",
					key,
					v,
				),
			)

		// -----------------------------------
		// fallback
		// -----------------------------------

		default:

			b.WriteString(
				fmt.Sprintf(
					"    %s = \"%v\"\n",
					key,
					v,
				),
			)
		}
	}

	b.WriteString(")\n\n")

	// generate helper functions if requested
	funcs := map[string][]string{}
	if raw, ok := constants["functions"]; ok {
		if fm, ok2 := raw.(map[string]interface{}); ok2 {
			for fname, val := range fm {
				switch arr := val.(type) {
				case []interface{}:
					for _, item := range arr {
						if s, ok := item.(string); ok {
							funcs[fname] = append(funcs[fname], s)
						}
					}
				case string:
					funcs[fname] = append(funcs[fname], arr)
				}
			}
		}
	}

	if len(funcs) > 0 {
		for fname, consts := range funcs {
			// assume string profile if first referenced constant is a string
			var isString bool
			for _, c := range consts {
				if v, ok := constants[c]; ok {
					if _, ok2 := v.(string); ok2 {
						isString = true
					}
					break
				}
			}

			// allow per-target signature override
			sig := ""
			if t, ok := spec.Targets[lang]; ok {
				sig = t.FunctionSignatures[fname]
			}

			if isString {
				if sig != "" {
					if strings.Contains(sig, "{") {
						// full function provided — emit as-is
						b.WriteString(fmt.Sprintf("%s\n", sig))
						continue
					}
					b.WriteString(fmt.Sprintf("%s {\n", sig))
				} else {
					b.WriteString(fmt.Sprintf("func %s(value string) bool {\n", fname))
				}

				b.WriteString("    v := strings.ToUpper(value)\n")
				for _, c := range consts {
					b.WriteString(fmt.Sprintf("    if v == %s { return true }\n", c))
				}
				b.WriteString("    return false\n}\n")
			} else {
				if sig != "" {
					if strings.Contains(sig, "{") {
						b.WriteString(fmt.Sprintf("%s\n", sig))
						continue
					}
					b.WriteString(fmt.Sprintf("%s {\n", sig))
				} else {
					b.WriteString(fmt.Sprintf("func %s(value interface{}) bool {\n", fname))
				}

				for _, c := range consts {
					b.WriteString(fmt.Sprintf("    if value == %s { return true }\n", c))
				}
				b.WriteString("    return false\n}\n")
			}
		}
	}

	return b.String()
}