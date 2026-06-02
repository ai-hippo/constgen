package runner

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ai-hippo/constgen/internal/generator"
	"github.com/ai-hippo/constgen/internal/parser"
)

// Runner scans input directory and processes YAML specs.
type Runner struct {
	InputDir  string
	OutputDir string
}

// New creates runner instance.
func New(input, output string) *Runner {

	return &Runner{
		InputDir:  input,
		OutputDir: output,
	}
}

// Run executes full generation lifecycle.
func (r *Runner) Run() error {

	// -----------------------------------
	// CLEAN GENERATED DIRECTORY FIRST
	// -----------------------------------

	//
	// IMPORTANT:
	//
	// We ALWAYS generate from a clean state.
	//
	// This prevents:
	// - stale files
	// - orphaned outputs
	// - deleted profile leftovers
	// - renamed spec artifacts
	//

	if err := os.RemoveAll(r.OutputDir); err != nil {
		return err
	}

	// recreate base output directory
	if err := os.MkdirAll(r.OutputDir, 0755); err != nil {
		return err
	}

	// -----------------------------------
	// PROCESS YAML FILES
	// -----------------------------------

	files, err := os.ReadDir(r.InputDir)
	if err != nil {
		return err
	}

	for _, file := range files {

		// skip directories + non-yaml files
		if file.IsDir() ||
			!strings.HasSuffix(file.Name(), ".yaml") {
			continue
		}

		fullPath := filepath.Join(
			r.InputDir,
			file.Name(),
		)

		log.Println("processing:", fullPath)

		// parse yaml
		spec, err := parser.ParseYAML(fullPath)
		if err != nil {
			return err
		}

		// validate yaml
		if err := parser.ValidateSpec(spec); err != nil {
			return err
		}

		// generate outputs
		engine := &generator.Engine{
			Spec: spec,
		}

		if err := engine.Generate(); err != nil {
			return err
		}
	}

	return nil
}