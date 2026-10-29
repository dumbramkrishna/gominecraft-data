// Package mcdata is a fast, type-safe Go accessor for the PrismarineJS
// minecraft-data dataset. It is a Go reimagining of node-minecraft-data.
//
// All game data is embedded at build time (gzip-compressed) and parsed into
// typed Go structs on first access, then cached, so repeated lookups are free.
//
//	md, err := mcdata.Load("1.21.4")
//	if err != nil { ... }
//	stone, ok := md.BlockByName("stone")
//
// Bedrock versions use the "bedrock_" (or legacy "pe_") prefix, e.g.
// mcdata.Load("bedrock_1.21.0").
package mcdata

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/dumbramkrishna/gominecraft-data/internal/mcdb"
)

// MCData is the typed dataset for one resolved Minecraft version. It is safe for
// concurrent use; each collection is parsed and indexed exactly once on first
// access and cached thereafter. Obtain one with Load.
type MCData struct {
	// Version is the resolved version this dataset describes.
	Version Version

	blocks      lazy[Blocks]
	items       lazy[Items]
	entities    lazy[Entities]
	biomes      lazy[Biomes]
	effects     lazy[Effects]
	enchants    lazy[Enchantments]
	foods       lazy[Foods]
	particles   lazy[Particles]
	instruments lazy[Instruments]
	sounds      lazy[Sounds]
}

var (
	cacheMu sync.Mutex
	cache   = map[string]*MCData{}
)

// Load returns the dataset for the given version string. The version may be a
// bare Minecraft version ("1.21.4"), a major version ("1.21"), or edition-
// prefixed ("bedrock_1.21.0", "pe_1.16"). Results are cached, so calling Load
// repeatedly for the same resolved version returns the same *MCData.
func Load(version string) (*MCData, error) {
	v, err := resolveVersion(version)
	if err != nil {
		return nil, err
	}
	key := v.Edition + "/" + v.DataFolder

	cacheMu.Lock()
	defer cacheMu.Unlock()
	if md, ok := cache[key]; ok {
		return md, nil
	}
	md := &MCData{Version: v}
	cache[key] = md
	return md, nil
}

// MustLoad is like Load but panics on error. Useful for package-level vars with
// a known-good version.
func MustLoad(version string) *MCData {
	md, err := Load(version)
	if err != nil {
		panic(err)
	}
	return md
}

// feature reads and JSON-decodes a feature file into dst. ok is false when the
// feature is simply not present for this version (dst is left untouched).
func (md *MCData) feature(name string, dst any) (ok bool, err error) {
	raw, ok, err := mcdb.ReadFeature(md.Version.Edition, md.Version.DataFolder, name)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	if err := json.Unmarshal(raw, dst); err != nil {
		return false, fmt.Errorf("mcdata: parse %s for %s: %w", name, md.Version, err)
	}
	return true, nil
}

// lazy memoizes a single computed value together with its error.
type lazy[T any] struct {
	once sync.Once
	val  T
	err  error
}

func (l *lazy[T]) get(load func() (T, error)) (T, error) {
	l.once.Do(func() { l.val, l.err = load() })
	return l.val, l.err
}
