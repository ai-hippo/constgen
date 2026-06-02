package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ai-hippo/constgen/internal/model"
)

// Engine orchestrates generation across languages and profiles.
// It delegates ALL value resolution to ProfileResolver inside generators.
type Engine struct {
	Spec *model.Spec
}

// Generator interface stays unchanged.
// Each generator is responsible only for rendering output.
type Generator interface {
	Generate(spec model.Spec, profile string, lang string) string
	FileName(spec model.Spec, profile string) string
}

// registry returns all supported language generators.
func registry() map[string]Generator {
	return map[string]Generator{
		"go":      GoGenerator{},
		"php":     PhpGenerator{},
		"java":    JavaGenerator{},
		"ruby":    RubyGenerator{},
		"aspnet":  AspnetGenerator{},
		"angular": AngularGenerator{},
		"react":   ReactGenerator{},
		"vue":     VueGenerator{},
		"python":  PythonGenerator{},
		"kotlin":  KotlinGenerator{},
	}
}

// Generate runs full pipeline:
// language → profile → output file
func (e *Engine) Generate() error {

	base := "generated"

	// compute registry once
	reg := registry()

	for lang, target := range e.Spec.Targets {

		if !target.Enabled {
			continue
		}

		gen, ok := reg[lang]
		if !ok {
			return fmt.Errorf("no generator found for language: %s", lang)
		}

		for _, profile := range target.Profiles {

			// IMPORTANT: spec is already inside engine → DO NOT pass spec again
			output := gen.Generate(*e.Spec, profile, lang)

			outDir := filepath.Join(base, lang, profile)

			if err := os.MkdirAll(outDir, 0755); err != nil {
				return err
			}

			fileName := gen.FileName(*e.Spec, profile)
			path := filepath.Join(outDir, fileName)

			if err := os.WriteFile(path, []byte(output), 0644); err != nil {
				return err
			}

			fmt.Println("generated:", path)
		}
	}

	return nil
}
