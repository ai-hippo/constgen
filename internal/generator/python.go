package generator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ai-hippo/constgen/internal/model"
)

// PythonGenerator emits a Python class with constants and static helpers.
type PythonGenerator struct{}

func (g PythonGenerator) FileName(spec model.Spec, profile string) string {
    return fmt.Sprintf("%s.py", spec.Name)
}

func (g PythonGenerator) Generate(spec model.Spec, profile string, lang string) string {
    var b strings.Builder

    constants, ok := spec.Profiles[profile]
    if !ok {
        return ""
    }

    className := strings.TrimSuffix(g.FileName(spec, profile), ".py")
    b.WriteString(fmt.Sprintf("class %s:\n", className))

    // deterministic ordering
    var keys []string
    for k := range constants {
        if k == "functions" {
            continue
        }
        keys = append(keys, k)
    }
    sort.Strings(keys)

    if len(keys) == 0 {
        b.WriteString("    pass\n")
    }

    for _, key := range keys {
        value := constants[key]
        switch v := value.(type) {
        case string:
            b.WriteString(fmt.Sprintf("    %s = %q\n", key, v))
        case bool:
            if v {
                b.WriteString(fmt.Sprintf("    %s = True\n", key))
            } else {
                b.WriteString(fmt.Sprintf("    %s = False\n", key))
            }
        default:
            b.WriteString(fmt.Sprintf("    %s = %v\n", key, v))
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

                if sig != "" && strings.Contains(sig, "def ") {
                    // user provided full function def — emit as-is
                    b.WriteString(fmt.Sprintf("%s\n", sig))
                    continue
                }

                // default: staticmethod with string comparison
                b.WriteString(fmt.Sprintf("    @staticmethod\n    def %s(value):\n", fname))
                b.WriteString("        if value is None:\n            return False\n")
                b.WriteString("        v = str(value).upper()\n")
                for _, c := range constsList {
                    b.WriteString(fmt.Sprintf("        if v == %s:\n            return True\n", c))
                }
                b.WriteString("        return False\n\n")
            }
        }
    }

    return b.String()
}
