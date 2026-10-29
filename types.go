package mcdata

import "encoding/json"

// This file defines the typed, JSON-mappable structs for each minecraft-data
// collection. Fields that vary in shape across Minecraft versions (for example
// block "drops", which may be plain IDs in some versions and structured objects
// in others) are kept as json.RawMessage so callers can decode them against the
// exact version's schema without us losing data.

// Block describes a single block type.
type Block struct {
	ID           int             `json:"id"`
	Name         string          `json:"name"`
	DisplayName  string          `json:"displayName"`
	Hardness     *float64        `json:"hardness,omitempty"`
	Resistance   float64         `json:"resistance"`
	StackSize    int             `json:"stackSize"`
	Diggable     bool            `json:"diggable"`
	Material     string          `json:"material,omitempty"`
	Transparent  bool            `json:"transparent"`
	EmitLight    int             `json:"emitLight"`
	FilterLight  int             `json:"filterLight"`
	DefaultState int             `json:"defaultState"`
	MinStateID   int             `json:"minStateId"`
	MaxStateID   int             `json:"maxStateId"`
	States       []BlockState    `json:"states,omitempty"`
	HarvestTools map[string]bool `json:"harvestTools,omitempty"`
	Drops        json.RawMessage `json:"drops,omitempty"`
	BoundingBox  string          `json:"boundingBox,omitempty"`
}

// BlockState is one property in a block's state definition.
type BlockState struct {
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	NumValues int      `json:"num_values"`
	Values    []string `json:"values,omitempty"`
}

// Item describes a single item type.
type Item struct {
	ID                int             `json:"id"`
	Name              string          `json:"name"`
	DisplayName       string          `json:"displayName"`
	StackSize         int             `json:"stackSize"`
	MaxDurability     *int            `json:"maxDurability,omitempty"`
	EnchantCategories []string        `json:"enchantCategories,omitempty"`
	RepairWith        []string        `json:"repairWith,omitempty"`
	Variations        json.RawMessage `json:"variations,omitempty"`
}

// Entity describes a single entity type.
type Entity struct {
	ID           int      `json:"id"`
	InternalID   *int     `json:"internalId,omitempty"`
	Name         string   `json:"name"`
	DisplayName  string   `json:"displayName"`
	Width        float64  `json:"width"`
	Height       float64  `json:"height"`
	Type         string   `json:"type"`
	Category     string   `json:"category,omitempty"`
	MetadataKeys []string `json:"metadataKeys,omitempty"`
}

// Biome describes a single biome.
type Biome struct {
	ID               int      `json:"id"`
	Name             string   `json:"name"`
	Category         string   `json:"category,omitempty"`
	Temperature      float64  `json:"temperature"`
	HasPrecipitation *bool    `json:"has_precipitation,omitempty"`
	Precipitation    string   `json:"precipitation,omitempty"`
	Dimension        string   `json:"dimension,omitempty"`
	DisplayName      string   `json:"displayName"`
	Color            int      `json:"color"`
	Rainfall         *float64 `json:"rainfall,omitempty"`
}

// Effect describes a status effect.
type Effect struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Type        string `json:"type"`
}

// EnchantmentCost is the linear cost model (a*level + b) for an enchantment.
type EnchantmentCost struct {
	A int `json:"a"`
	B int `json:"b"`
}

// Enchantment describes an enchantment.
type Enchantment struct {
	ID           int             `json:"id"`
	Name         string          `json:"name"`
	DisplayName  string          `json:"displayName"`
	MaxLevel     int             `json:"maxLevel"`
	MinCost      EnchantmentCost `json:"minCost"`
	MaxCost      EnchantmentCost `json:"maxCost"`
	TreasureOnly bool            `json:"treasureOnly"`
	Curse        bool            `json:"curse"`
	Exclude      []string        `json:"exclude,omitempty"`
	Category     string          `json:"category,omitempty"`
	Weight       int             `json:"weight"`
	Tradeable    bool            `json:"tradeable"`
	Discoverable bool            `json:"discoverable"`
}

// Food describes a consumable food item.
type Food struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	DisplayName      string  `json:"displayName"`
	StackSize        int     `json:"stackSize"`
	FoodPoints       float64 `json:"foodPoints"`
	Saturation       float64 `json:"saturation"`
	EffectiveQuality float64 `json:"effectiveQuality"`
	SaturationRatio  float64 `json:"saturationRatio"`
}

// Window describes a GUI window / inventory.
type Window struct {
	ID     string          `json:"id"`
	Name   string          `json:"name"`
	Slots  json.RawMessage `json:"slots,omitempty"`
	Opened json.RawMessage `json:"openedWith,omitempty"`
}

// Particle describes a particle type.
type Particle struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Instrument describes a note-block instrument.
type Instrument struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// MapIcon describes a map marker icon.
type MapIcon struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	AppearsOnMz string `json:"appearsOnItemFrame,omitempty"`
	Visible     *bool  `json:"visibleInItemFrame,omitempty"`
}

// Sound describes a sound effect.
type Sound struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Attribute describes an entity attribute.
type Attribute struct {
	Name     string          `json:"name"`
	Resource string          `json:"resource"`
	Min      float64         `json:"min"`
	Max      float64         `json:"max"`
	Default  float64         `json:"default"`
	Type     json.RawMessage `json:"type,omitempty"`
}

// BlockLoot describes the loot table entry for a block.
type BlockLoot struct {
	Block string          `json:"block"`
	Drops json.RawMessage `json:"drops"`
}

// EntityLoot describes the loot table entry for an entity.
type EntityLoot struct {
	Entity string          `json:"entity"`
	Drops  json.RawMessage `json:"drops"`
}
