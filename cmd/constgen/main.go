package main

import (
	"flag"
	"log"

	"github.com/ai-hippo/constgen/internal/runner"
)

func main() {

	// CLI command mode
	command := flag.String("cmd", "generate", "command: generate")
	inputDir := flag.String("input", "yaml", "input directory for YAML specs")
	outputDir := flag.String("out", "generated", "output directory")

	flag.Parse()

	switch *command {

	case "generate":
		r := runner.New(*inputDir, *outputDir)
		if err := r.Run(); err != nil {
			log.Fatalf("generation failed: %v", err)
		}

	default:
		log.Fatalf("unknown command: %s", *command)
	}
}