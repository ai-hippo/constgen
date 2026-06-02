package generator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ai-hippo/constgen/internal/model"
)

// ReactGenerator emits a JS module suitable for React projects.
type ReactGenerator struct{}

func (g ReactGenerator) FileName(spec model.Spec, profile string) string {
    return fmt.Sprintf(
        "%s.js",
        spec.Name,
    )
}

func (g ReactGenerator) Generate(spec model.Spec, profile string, lang string) string {
    var b strings.Builder

    constants, ok := spec.Profiles[profile]
    if !ok {
        return ""
    }

    name := strings.TrimSuffix(g.FileName(spec, profile), ".js")
    b.WriteString(fmt.Sprintf("export const %s = {\n", name))

    // deterministic ordering
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
        case string:
            b.WriteString(fmt.Sprintf("  \"%s\": \"%s\",\n", key, v))
        case bool:
            b.WriteString(fmt.Sprintf("  \"%s\": %t,\n", key, v))
        default:
            b.WriteString(fmt.Sprintf("  \"%s\": %v,\n", key, v))
        }
    }

    // functions
    if raw, ok := constants["functions"]; ok {
        if fm, ok2 := raw.(map[string]interface{}); ok2 {
            for fname, val := range fm {
                var constsList []string
                switch arr := val.(type) {
                case []interface{}:
                    for _, item := range arr {
                        if s, ok := item.(string); ok {
                            constsList = append(constsList, s)
                        }
                    }
                case string:
                    constsList = append(constsList, arr)
                }

                sig := ""
                if t, ok := spec.Targets[lang]; ok {
                    sig = t.FunctionSignatures[fname]
                }

                if sig != "" {
                    if strings.Contains(sig, "{") {
                        b.WriteString(fmt.Sprintf("  %s\n", sig))
                    } else {
                        b.WriteString(fmt.Sprintf("  %s: %s,\n", fname, sig))
                    }
                } else {
                    b.WriteString(fmt.Sprintf("  \"%s\": (value) => {\n", fname))
                    b.WriteString("    const v = String(value).toUpperCase();\n")
                    for _, c := range constsList {
                        b.WriteString(fmt.Sprintf("    if (v === %s) return true;\n", c))
                    }
                    b.WriteString("    return false;\n  },\n")
                }
            }
        }
    }

    b.WriteString("};\n")

    return b.String()
}
