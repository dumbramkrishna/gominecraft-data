// Package mcdb holds the embedded, gzip-compressed minecraft-data assets and
// the generated routing table that locates them. It is an internal detail of
// gominecraft-data; the public API lives in the root mcdata package.
package mcdb

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

// ReadFeature returns the decompressed bytes of a feature file for an edition
// and version, or ok=false if that feature is not present for the version.
func ReadFeature(edition, version, feature string) (data []byte, ok bool, err error) {
	versions, ok := Paths[edition]
	if !ok {
		return nil, false, nil
	}
	features, ok := versions[version]
	if !ok {
		return nil, false, nil
	}
	path, ok := features[feature]
	if !ok {
		return nil, false, nil
	}
	data, err = readAsset(path)
	return data, err == nil, err
}

// ReadCommon returns the decompressed bytes of a per-edition common metadata
// file (e.g. "protocolVersions"), or ok=false if absent.
func ReadCommon(edition, name string) (data []byte, ok bool, err error) {
	return ReadFeature(edition, "common", "common/"+name)
}

func readAsset(path string) ([]byte, error) {
	raw, err := AssetsFS.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("asset %s: %w", path, err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("gunzip %s: %w", path, err)
	}
	defer zr.Close()
	out, err := io.ReadAll(zr)
	if err != nil {
		return nil, fmt.Errorf("decompress %s: %w", path, err)
	}
	return out, nil
}

// Editions returns the available edition keys ("pc", "bedrock").
func Editions() []string {
	eds := make([]string, 0, len(Paths))
	for e := range Paths {
		eds = append(eds, e)
	}
	return eds
}
