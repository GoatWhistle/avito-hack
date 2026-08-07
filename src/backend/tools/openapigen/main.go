//go:build tools

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	openAPIVersion = "3.1.0"
	specFileMode   = 0o600
)

var errMergedSpecNotMapping = errors.New("merged spec is not a mapping")

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "usage: openapigen <swagger.yaml> <overlay.yaml> <out.yaml>")
		os.Exit(1)
	}

	if err := run(os.Args[1], os.Args[2], os.Args[3]); err != nil {
		fmt.Fprintln(os.Stderr, "openapigen:", err)
		os.Exit(1)
	}
}

func run(generatedPath, overlayPath, outPath string) error {
	generated, err := readYAML(generatedPath)
	if err != nil {
		return fmt.Errorf("read generated spec: %w", err)
	}

	overlay, err := readYAML(overlayPath)
	if err != nil {
		return fmt.Errorf("read overlay: %w", err)
	}

	spec, ok := merge(normalize(generated), overlay).(map[string]any)
	if !ok {
		return errMergedSpecNotMapping
	}

	spec["openapi"] = openAPIVersion

	delete(spec, "externalDocs")

	out, err := yaml.Marshal(spec)
	if err != nil {
		return fmt.Errorf("marshal spec: %w", err)
	}

	header := "# Generated from Go annotations by `make api-spec` — do not edit by hand.\n" +
		"# Sources: swaggo annotations in src/backend and docs/openapi.overlay.yaml.\n"

	return os.WriteFile(outPath, append([]byte(header), out...), specFileMode)
}

func readYAML(path string) (map[string]any, error) {
	raw, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	var doc map[string]any
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		return nil, err
	}

	return doc, nil
}

func normalize(spec map[string]any) map[string]any {
	renameSchemas(spec)
	walk(spec, fixRequestBody)
	walk(spec, fixEmptySecurity)
	walk(spec, fixMultipart)

	return spec
}
