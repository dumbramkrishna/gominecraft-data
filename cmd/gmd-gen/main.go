// Command gmd-gen builds the embedded data assets and routing table for
// gominecraft-data from the minecraft-data submodule.
//
// For every (edition, version, feature) entry in minecraft-data's
// data/dataPaths.json it finds the backing file, gzip-compresses it into
// internal/mcdb/assets/, and records its path in a generated Go routing table
// (internal/mcdb/paths_gen.go). The per-edition common metadata files
// (protocolVersions.json, versions.json, features.json, legacy.json) are
// embedded the same way.
//
// Run from the repository root:
//
//	go run ./cmd/gmd-gen
package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

func main() {
	dataDir := flag.String("data", filepath.Join("minecraft-data", "data"), "path to minecraft-data/data")
	outDir := flag.String("out", filepath.Join("internal", "mcdb"), "output package directory")
	flag.Parse()

	if err := run(*dataDir, *outDir); err != nil {
		fmt.Fprintln(os.Stderr, "gmd-gen:", err)
		os.Exit(1)
	}
}

// routing is editions -> versions -> features -> embedded asset path.
type routing map[string]map[string]map[string]string

func run(dataDir, outDir string) error {
	dp, err := os.ReadFile(filepath.Join(dataDir, "dataPaths.json"))
	if err != nil {
		return fmt.Errorf("read dataPaths.json: %w", err)
	}
	var paths map[string]map[string]map[string]string
	if err := json.Unmarshal(dp, &paths); err != nil {
		return fmt.Errorf("parse dataPaths.json: %w", err)
	}

	assetsDir := filepath.Join(outDir, "assets")
	if err := os.RemoveAll(assetsDir); err != nil {
		return fmt.Errorf("clean assets: %w", err)
	}

	route := routing{}
	var total, embedded int

	for _, edition := range sortedKeys(paths) {
		route[edition] = map[string]map[string]string{}
		for _, version := range sortedKeys(paths[edition]) {
			route[edition][version] = map[string]string{}
			for _, feature := range sortedKeys(paths[edition][version]) {
				total++
				srcDir := paths[edition][version][feature] // e.g. "pc/1.21.4"
				file, ok := findFeatureFile(dataDir, srcDir, feature)
				if !ok {
					// A declared feature with no backing file: skip quietly;
					// the loader reports it as "not present for this version".
					continue
				}
				assetRel := filepath.ToSlash(filepath.Join("assets", srcDir, filepath.Base(file)+".gz"))
				dst := filepath.Join(outDir, assetRel)
				if err := gzipCopy(file, dst); err != nil {
					return fmt.Errorf("gzip %s: %w", file, err)
				}
				route[edition][version][feature] = assetRel
				embedded++
			}
		}
	}

	// Embed per-edition common metadata too.
	for _, edition := range sortedKeys(paths) {
		commonDir := filepath.Join(dataDir, edition, "common")
		entries, err := os.ReadDir(commonDir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
				continue
			}
			src := filepath.Join(commonDir, e.Name())
			assetRel := filepath.ToSlash(filepath.Join("assets", edition, "common", e.Name()+".gz"))
			if err := gzipCopy(src, filepath.Join(outDir, assetRel)); err != nil {
				return fmt.Errorf("gzip common %s: %w", src, err)
			}
			feature := "common/" + strings.TrimSuffix(e.Name(), ".json")
			if route[edition]["common"] == nil {
				route[edition]["common"] = map[string]string{}
			}
			route[edition]["common"][feature] = assetRel
		}
	}

	if err := writePathsFile(outDir, route); err != nil {
		return err
	}

	fmt.Printf("gmd-gen: embedded %d/%d feature files across %d editions\n", embedded, total, len(paths))
	return nil
}

// findFeatureFile locates the file backing a feature inside data/<srcDir>.
// Features normally map to "<feature>.json" but some (e.g. proto) use another
// extension, so we fall back to the first file named "<feature>.*".
func findFeatureFile(dataDir, srcDir, feature string) (string, bool) {
	dir := filepath.Join(dataDir, srcDir)
	if p := filepath.Join(dir, feature+".json"); fileExists(p) {
		return p, true
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", false
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasPrefix(name, feature+".") {
			return filepath.Join(dir, name), true
		}
	}
	return "", false
}

func gzipCopy(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	var buf bytes.Buffer
	zw, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return err
	}
	if _, err := zw.Write(data); err != nil {
		return err
	}
	if err := zw.Close(); err != nil {
		return err
	}
	return os.WriteFile(dst, buf.Bytes(), 0o644)
}

func fileExists(p string) bool {
	info, err := os.Stat(p)
	return err == nil && !info.IsDir()
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

var pathsTmpl = template.Must(template.New("paths").Parse(`// Code generated by cmd/gmd-gen. DO NOT EDIT.

package mcdb

import "embed"

// AssetsFS holds every gzip-compressed data file embedded at build time.
//
//go:embed assets
var AssetsFS embed.FS

// Paths maps edition -> version -> feature -> the embedded asset path (a
// gzip-compressed file inside AssetsFS). The "common" pseudo-version holds
// per-edition metadata (keys prefixed "common/").
var Paths = map[string]map[string]map[string]string{
{{- range $edition := .Editions}}
	{{printf "%q" $edition}}: {
	{{- range $version := index $.Order $edition}}
		{{printf "%q" $version}}: {
		{{- range $f := index (index $.Features $edition) $version}}
			{{printf "%q" $f.Feature}}: {{printf "%q" $f.Path}},
		{{- end}}
		},
	{{- end}}
	},
{{- end}}
}
`))

type featurePath struct {
	Feature string
	Path    string
}

func writePathsFile(outDir string, route routing) error {
	editions := sortedKeys(route)
	order := map[string][]string{}
	features := map[string]map[string][]featurePath{}
	for _, ed := range editions {
		order[ed] = sortedKeys(route[ed])
		features[ed] = map[string][]featurePath{}
		for _, ver := range order[ed] {
			fs := make([]featurePath, 0, len(route[ed][ver]))
			for _, f := range sortedKeys(route[ed][ver]) {
				fs = append(fs, featurePath{Feature: f, Path: route[ed][ver][f]})
			}
			features[ed][ver] = fs
		}
	}

	var buf bytes.Buffer
	err := pathsTmpl.Execute(&buf, map[string]any{
		"Editions": editions,
		"Order":    order,
		"Features": features,
	})
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(outDir, "paths_gen.go"), buf.Bytes(), 0o644)
}
