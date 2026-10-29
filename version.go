package mcdata

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/dumbramkrishna/gominecraft-data/internal/mcdb"
)

// Version identifies a resolved Minecraft version within an edition and carries
// the metadata needed to order versions against one another. Comparisons use
// the integer DataVersion, the same ordering the Minecraft client uses.
type Version struct {
	// Edition is "pc" (Java) or "bedrock".
	Edition string
	// MinecraftVersion is the human version string, e.g. "1.21.4".
	MinecraftVersion string
	// MajorVersion is the version family, e.g. "1.21".
	MajorVersion string
	// DataFolder is the data-bearing version key used to load assets.
	DataFolder string
	// Version is the protocol version number.
	Version int
	// DataVersion is the client "data version"; the canonical ordering key.
	DataVersion int
	// UsesNetty is true for post-Netty (modern) protocol versions.
	UsesNetty bool
	// ReleaseType is "release", "snapshot", etc. when known.
	ReleaseType string
}

// protoEntry mirrors one record in minecraft-data's protocolVersions.json.
type protoEntry struct {
	MinecraftVersion string `json:"minecraftVersion"`
	Version          int    `json:"version"`
	DataVersion      int    `json:"dataVersion"`
	UsesNetty        bool   `json:"usesNetty"`
	MajorVersion     string `json:"majorVersion"`
	ReleaseType      string `json:"releaseType"`
}

// registry holds the resolved version metadata for one edition.
type registry struct {
	byMinecraft map[string]protoEntry
	byMajor     map[string]protoEntry // newest entry wins
	dataFolders map[string]struct{}   // version keys that have embedded data
	folderOrder []string              // data folders sorted by DataVersion asc
	folderMeta  map[string]protoEntry
}

var (
	registriesMu sync.Mutex
	registries   = map[string]*registry{}
)

func getRegistry(edition string) (*registry, error) {
	registriesMu.Lock()
	defer registriesMu.Unlock()
	if r, ok := registries[edition]; ok {
		return r, nil
	}
	r, err := buildRegistry(edition)
	if err != nil {
		return nil, err
	}
	registries[edition] = r
	return r, nil
}

func buildRegistry(edition string) (*registry, error) {
	raw, ok, err := mcdb.ReadCommon(edition, "protocolVersions")
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("mcdata: unknown edition %q", edition)
	}
	var entries []protoEntry
	if err := json.Unmarshal(raw, &entries); err != nil {
		return nil, fmt.Errorf("mcdata: parse protocolVersions for %q: %w", edition, err)
	}

	r := &registry{
		byMinecraft: make(map[string]protoEntry, len(entries)),
		byMajor:     make(map[string]protoEntry, len(entries)),
		dataFolders: map[string]struct{}{},
		folderMeta:  map[string]protoEntry{},
	}
	// protocolVersions.json is newest-first. Iterate reversed so newer entries
	// overwrite older ones and win for a given major version.
	for i := len(entries) - 1; i >= 0; i-- {
		e := entries[i]
		r.byMinecraft[e.MinecraftVersion] = e
		if e.MajorVersion != "" {
			r.byMajor[e.MajorVersion] = e
		}
	}

	for ver := range mcdb.Paths[edition] {
		if ver == "common" {
			continue
		}
		r.dataFolders[ver] = struct{}{}
		if meta, ok := r.folderMetaFor(ver); ok {
			r.folderMeta[ver] = meta
		}
	}

	for ver := range r.dataFolders {
		r.folderOrder = append(r.folderOrder, ver)
	}
	sort.Slice(r.folderOrder, func(i, j int) bool {
		return r.folderMeta[r.folderOrder[i]].DataVersion < r.folderMeta[r.folderOrder[j]].DataVersion
	})
	return r, nil
}

// folderMetaFor finds the proto metadata for a data folder key, trying an exact
// minecraft-version match then a major-version match.
func (r *registry) folderMetaFor(ver string) (protoEntry, bool) {
	if e, ok := r.byMinecraft[ver]; ok {
		return e, true
	}
	if e, ok := r.byMajor[ver]; ok {
		return e, true
	}
	return protoEntry{}, false
}

// resolveVersion turns a user-supplied version string into a Version bound to a
// data-bearing folder, or an error if nothing matches.
func resolveVersion(input string) (Version, error) {
	edition, ver := splitEdition(input)
	r, err := getRegistry(edition)
	if err != nil {
		return Version{}, err
	}

	// 1. Direct data folder hit.
	if _, ok := r.dataFolders[ver]; ok {
		return r.versionFor(edition, ver), nil
	}
	// 2. Known minecraft version -> its major version's data folder.
	if e, ok := r.byMinecraft[ver]; ok {
		if folder, ok := r.folderForDataVersion(e.DataVersion); ok {
			return r.versionFor(edition, folder), nil
		}
	}
	// 3. Known major version with data.
	if e, ok := r.byMajor[ver]; ok {
		if folder, ok := r.folderForDataVersion(e.DataVersion); ok {
			return r.versionFor(edition, folder), nil
		}
	}
	return Version{}, fmt.Errorf("mcdata: no data for version %q (edition %q)", ver, edition)
}

// folderForDataVersion returns the data folder whose DataVersion is the largest
// value <= target (the newest schema that applies to that version).
func (r *registry) folderForDataVersion(target int) (string, bool) {
	best := ""
	bestDV := -1 << 62 // sentinel below any real DataVersion
	for _, folder := range r.folderOrder {
		dv := r.folderMeta[folder].DataVersion
		if dv <= target && dv > bestDV {
			bestDV = dv
			best = folder
		}
	}
	if best == "" {
		return "", false
	}
	return best, true
}

func (r *registry) versionFor(edition, folder string) Version {
	meta := r.folderMeta[folder]
	return Version{
		Edition:          edition,
		MinecraftVersion: pick(meta.MinecraftVersion, folder),
		MajorVersion:     pick(meta.MajorVersion, folder),
		DataFolder:       folder,
		Version:          meta.Version,
		DataVersion:      meta.DataVersion,
		UsesNetty:        meta.UsesNetty,
		ReleaseType:      meta.ReleaseType,
	}
}

// splitEdition parses an optional edition prefix: "pc_1.19", "bedrock_1.19",
// "pe_1.19" (alias for bedrock), or a bare version (defaults to pc).
func splitEdition(input string) (edition, version string) {
	input = strings.TrimSpace(input)
	if strings.HasPrefix(input, "pe_") {
		return "bedrock", strings.TrimPrefix(input, "pe_")
	}
	for _, ed := range []string{"pc_", "bedrock_"} {
		if strings.HasPrefix(input, ed) {
			return strings.TrimSuffix(ed, "_"), strings.TrimPrefix(input, ed)
		}
	}
	return "pc", input
}

func pick(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

// Comparison operators. Each compares by DataVersion.

// Compare returns -1, 0, or 1 as v is older than, equal to, or newer than other.
func (v Version) Compare(other Version) int {
	switch {
	case v.DataVersion < other.DataVersion:
		return -1
	case v.DataVersion > other.DataVersion:
		return 1
	default:
		return 0
	}
}

// GTE reports whether v is at least as new as the named version.
func (v Version) GTE(other string) (bool, error) {
	return v.cmp(other, func(c int) bool { return c >= 0 })
}

// GT reports whether v is newer than the named version.
func (v Version) GT(other string) (bool, error) {
	return v.cmp(other, func(c int) bool { return c > 0 })
}

// LTE reports whether v is at most as new as the named version.
func (v Version) LTE(other string) (bool, error) {
	return v.cmp(other, func(c int) bool { return c <= 0 })
}

// LT reports whether v is older than the named version.
func (v Version) LT(other string) (bool, error) {
	return v.cmp(other, func(c int) bool { return c < 0 })
}

// EQ reports whether v resolves to the same data version as the named version.
func (v Version) EQ(other string) (bool, error) {
	return v.cmp(other, func(c int) bool { return c == 0 })
}

func (v Version) cmp(other string, pred func(int) bool) (bool, error) {
	o, err := resolveVersion(editionQualified(v.Edition, other))
	if err != nil {
		return false, err
	}
	return pred(v.Compare(o)), nil
}

func editionQualified(edition, version string) string {
	if strings.Contains(version, "_") {
		return version
	}
	return edition + "_" + version
}

// String renders the version as "edition_minecraftVersion".
func (v Version) String() string {
	return v.Edition + "_" + v.MinecraftVersion
}
