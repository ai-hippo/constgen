package generator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ai-hippo/constgen/internal/model"
)

// PhpGenerator emits PHP arrays of constants.
type PhpGenerator struct{}

func (g PhpGenerator) FileName(spec model.Spec, profile string) string {
    return fmt.Sprintf(
        "%s.php",
        spec.Name,
    )
}

func (g PhpGenerator) Generate(spec model.Spec, profile string, lang string) string {
    var b strings.Builder

    constants, ok := spec.Profiles[profile]
    if !ok {
        return ""
    }

    b.WriteString("<?php\n\n")

    // class name equals file base (without extension)
    className := strings.TrimSuffix(g.FileName(spec, profile), ".php")
    b.WriteString(fmt.Sprintf("class %s {\n", className))

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
            b.WriteString(fmt.Sprintf("    public const %s = %s;\n", key, fmt.Sprintf("%q", v)))
        case bool:
            b.WriteString(fmt.Sprintf("    public const %s = %t;\n", key, v))
        default:
            b.WriteString(fmt.Sprintf("    public const %s = %v;\n", key, v))
        }
    }

    // functions
    if raw, ok := constants["functions"]; ok {
        if fm, ok2 := raw.(map[string]interface{}); ok2 {
            for fname, val := range fm {
                // collect constant names
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

                    // optional override from target
                    sig := ""
                    if t, ok := spec.Targets[lang]; ok {
                        sig = t.FunctionSignatures[fname]
                    }

                    if sig != "" {
                        if strings.Contains(sig, "{") {
                            // user provided full function body — emit as-is and skip default body
                            b.WriteString(fmt.Sprintf("    %s\n", sig))
                            continue
                        }
                        b.WriteString(fmt.Sprintf("    %s {\n", sig))
                    } else {
                        b.WriteString(fmt.Sprintf("    public static function %s($value) {\n", fname))
                    }

                    // assume string comparisons by default
                    b.WriteString("        $v = strtoupper($value);\n")
                    for _, c := range consts {
                        b.WriteString(fmt.Sprintf("        if ($v === self::%s) { return true; }\n", c))
                    }
                    b.WriteString("        return false;\n    }\n")
            }
        }
    }

    b.WriteString("}\n")

    return b.String()
}
