package generator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ai-hippo/constgen/internal/model"
)

// RubyGenerator emits a Ruby module with constants.
type RubyGenerator struct{}

func (g RubyGenerator) FileName(spec model.Spec, profile string) string {
    return fmt.Sprintf(
        "%s.rb",
        spec.Name,
    )
}

func (g RubyGenerator) Generate(spec model.Spec, profile string, lang string) string {
    var b strings.Builder

    constants, ok := spec.Profiles[profile]
    if !ok {
        return ""
    }
    // class/module name equals file base (without extension)
    name := strings.TrimSuffix(g.FileName(spec, profile), ".rb")
    b.WriteString(fmt.Sprintf("module %s\n", name))

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
            b.WriteString(fmt.Sprintf("  %s = %q\n", key, v))
        default:
            b.WriteString(fmt.Sprintf("  %s = %v\n", key, v))
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

                sig := ""
                if t, ok := spec.Targets[lang]; ok {
                    sig = t.FunctionSignatures[fname]
                }
                if sig != "" {
                    if strings.Contains(sig, "{") || strings.Contains(sig, "end") {
                        // full method provided — emit and skip default body
                        b.WriteString(fmt.Sprintf("  %s\n", sig))
                        continue
                    }
                    b.WriteString(fmt.Sprintf("  %s\n", sig))
                } else {
                    b.WriteString(fmt.Sprintf("  def self.%s(value)\n", fname))
                }
                b.WriteString("    return false if value.nil?\n")
                b.WriteString("    v = value.to_s.upcase\n")
                for _, c := range consts {
                    b.WriteString(fmt.Sprintf("    return true if v == %s\n", c))
                }
                b.WriteString("    false\n  end\n")
            }
        }
    }

    b.WriteString("end\n")

    return b.String()
}
