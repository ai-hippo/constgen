package generator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ai-hippo/constgen/internal/model"
)

// KotlinGenerator emits a Kotlin object with constants and helper functions.
type KotlinGenerator struct{}

func (g KotlinGenerator) FileName(spec model.Spec, profile string) string {
    return fmt.Sprintf("%s.kt", spec.Name)
}

func (g KotlinGenerator) Generate(spec model.Spec, profile string, lang string) string {
    var b strings.Builder

    constants, ok := spec.Profiles[profile]
    if !ok {
        return ""
    }

    name := strings.TrimSuffix(g.FileName(spec, profile), ".kt")
    b.WriteString(fmt.Sprintf("object %s {\n", name))

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
            b.WriteString(fmt.Sprintf("    const val %s = \"%s\"\n", key, v))
        case int:
            b.WriteString(fmt.Sprintf("    val %s = %d\n", key, v))
        case int64:
            b.WriteString(fmt.Sprintf("    val %s = %dL\n", key, v))
        case float64:
            b.WriteString(fmt.Sprintf("    val %s = %v\n", key, v))
        case bool:
            b.WriteString(fmt.Sprintf("    val %s = %t\n", key, v))
        default:
            b.WriteString(fmt.Sprintf("    val %s = \"%v\"\n", key, v))
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

                if sig != "" && strings.Contains(sig, "fun ") {
                    b.WriteString(fmt.Sprintf("    %s\n", sig))
                    continue
                }

                // default: string-based helper
                b.WriteString(fmt.Sprintf("    fun %s(value: String): Boolean {\n", fname))
                b.WriteString("        val v = value.toUpperCase()\n")
                for _, c := range constsList {
                    b.WriteString(fmt.Sprintf("        if (v == %s) return true\n", c))
                }
                b.WriteString("        return false\n    }\n")
            }
        }
    }

    b.WriteString("}\n")

    return b.String()
}
