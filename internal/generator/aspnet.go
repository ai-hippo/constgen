package generator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ai-hippo/constgen/internal/model"
)

// AspnetGenerator emits a C# class suitable for ASP.NET projects.
type AspnetGenerator struct{}

func (g AspnetGenerator) FileName(spec model.Spec, profile string) string {
    return fmt.Sprintf(
        "%s.cs",
        spec.Name,
    )
}

func (g AspnetGenerator) Generate(spec model.Spec, profile string, lang string) string {
    var b strings.Builder

    constants, ok := spec.Profiles[profile]
    if !ok {
        return ""
    }

    className := strings.TrimSuffix(g.FileName(spec, profile), ".cs")
    b.WriteString(fmt.Sprintf("namespace Generated {\n    public static class %s {\n", className))

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
            b.WriteString(fmt.Sprintf("        public const string %s = \"%s\";\n", key, v))
        case int:
            b.WriteString(fmt.Sprintf("        public const int %s = %d;\n", key, v))
        case int64:
            b.WriteString(fmt.Sprintf("        public const long %s = %dL;\n", key, v))
        case float64:
            b.WriteString(fmt.Sprintf("        public const double %s = %v;\n", key, v))
        case bool:
            b.WriteString(fmt.Sprintf("        public const bool %s = %t;\n", key, v))
        default:
            b.WriteString(fmt.Sprintf("        public const string %s = \"%v\";\n", key, v))
        }
    }

    // functions
    if raw, ok := constants["functions"]; ok {
        if fm, ok2 := raw.(map[string]interface{}); ok2 {
            for fname, val := range fm {
                var consts []string
                switch arr := val.(type) {
                case []interface{}:
                    for _, item := range arr {
                        if s, ok := item.(string); ok {
                            consts = append(consts, s)
                        }
                    }
                case string:
                    consts = append(consts, arr)
                }

                // determine if string-based
                isString := false
                for _, c := range consts {
                    if v, ok := constants[c]; ok {
                        if _, ok2 := v.(string); ok2 {
                            isString = true
                            break
                        }
                    }
                }

                if isString {
                    sig := ""
                    if t, ok := spec.Targets[lang]; ok {
                        sig = t.FunctionSignatures[fname]
                    }
                    if sig != "" {
                        if strings.Contains(sig, "{") {
                            // full function provided — emit and skip default body
                            b.WriteString(fmt.Sprintf("        %s\n", sig))
                            continue
                        }
                        b.WriteString(fmt.Sprintf("        %s {\n", sig))
                    } else {
                        b.WriteString(fmt.Sprintf("        public static bool %s(string value) {\n", fname))
                    }
                    b.WriteString("            if (value == null) return false;\n")
                    for _, c := range consts {
                        b.WriteString(fmt.Sprintf("            if (string.Equals(value, %s, System.StringComparison.OrdinalIgnoreCase)) return true;\n", c))
                    }
                    b.WriteString("            return false;\n        }\n")
                } else {
                    sig := ""
                    if t, ok := spec.Targets[lang]; ok {
                        sig = t.FunctionSignatures[fname]
                    }
                    if sig != "" {
                        if strings.Contains(sig, "{") {
                            b.WriteString(fmt.Sprintf("        %s\n", sig))
                        } else {
                            b.WriteString(fmt.Sprintf("        %s {\n", sig))
                        }
                    } else {
                        b.WriteString(fmt.Sprintf("        public static bool %s(int value) {\n", fname))
                    }
                    for _, c := range consts {
                        b.WriteString(fmt.Sprintf("            if (value == %s) return true;\n", c))
                    }
                    b.WriteString("            return false;\n        }\n")
                }
            }
        }
    }

    b.WriteString("    }\n}\n")

    return b.String()
}
