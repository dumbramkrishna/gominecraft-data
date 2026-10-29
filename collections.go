package mcdata

// This file defines the indexed collection types and the lazy accessors that
// build them. Each accessor parses its feature file once (via MCData.feature)
// and builds the by-id / by-name indexes, mirroring node-minecraft-data's
// blocksByName / itemsById / ... surface but with typed Go maps.

// Blocks is the indexed block collection.
type Blocks struct {
	Array     []Block
	ByID      map[int]*Block
	ByName    map[string]*Block
	ByStateID map[int]*Block
}

// Items is the indexed item collection.
type Items struct {
	Array  []Item
	ByID   map[int]*Item
	ByName map[string]*Item
}

// Entities is the indexed entity collection.
type Entities struct {
	Array    []Entity
	ByID     map[int]*Entity
	ByName   map[string]*Entity
	MobsByID map[int]*Entity
	ObjsByID map[int]*Entity
}

// Biomes is the indexed biome collection.
type Biomes struct {
	Array  []Biome
	ByID   map[int]*Biome
	ByName map[string]*Biome
}

// Effects is the indexed status-effect collection.
type Effects struct {
	Array  []Effect
	ByID   map[int]*Effect
	ByName map[string]*Effect
}

// Enchantments is the indexed enchantment collection.
type Enchantments struct {
	Array  []Enchantment
	ByID   map[int]*Enchantment
	ByName map[string]*Enchantment
}

// Foods is the indexed food collection.
type Foods struct {
	Array  []Food
	ByID   map[int]*Food
	ByName map[string]*Food
}

// Particles is the indexed particle collection.
type Particles struct {
	Array  []Particle
	ByID   map[int]*Particle
	ByName map[string]*Particle
}

// Instruments is the indexed instrument collection.
type Instruments struct {
	Array []Instrument
	ByID  map[int]*Instrument
}

// Sounds is the indexed sound collection.
type Sounds struct {
	Array  []Sound
	ByID   map[int]*Sound
	ByName map[string]*Sound
}

// Blocks returns the indexed block collection for this version.
func (md *MCData) Blocks() (Blocks, error) {
	return md.blocks.get(func() (Blocks, error) {
		var arr []Block
		if _, err := md.feature("blocks", &arr); err != nil {
			return Blocks{}, err
		}
		b := Blocks{
			Array:     arr,
			ByID:      make(map[int]*Block, len(arr)),
			ByName:    make(map[string]*Block, len(arr)),
			ByStateID: map[int]*Block{},
		}
		for i := range arr {
			blk := &b.Array[i]
			// Fill state-id bounds for versions that omit them (flattening era
			// uses id<<4 .. +15 with defaultState = minStateId).
			if blk.MaxStateID == 0 && blk.MinStateID == 0 && blk.DefaultState == 0 {
				blk.MinStateID = blk.ID << 4
				blk.MaxStateID = blk.MinStateID + 15
				blk.DefaultState = blk.MinStateID
			}
			b.ByID[blk.ID] = blk
			b.ByName[blk.Name] = blk
			for s := blk.MinStateID; s <= blk.MaxStateID; s++ {
				b.ByStateID[s] = blk
			}
		}
		return b, nil
	})
}

// Items returns the indexed item collection for this version.
func (md *MCData) Items() (Items, error) {
	return md.items.get(func() (Items, error) {
		var arr []Item
		if _, err := md.feature("items", &arr); err != nil {
			return Items{}, err
		}
		it := Items{Array: arr, ByID: make(map[int]*Item, len(arr)), ByName: make(map[string]*Item, len(arr))}
		for i := range arr {
			it.ByID[arr[i].ID] = &it.Array[i]
			it.ByName[arr[i].Name] = &it.Array[i]
		}
		return it, nil
	})
}

// Entities returns the indexed entity collection for this version.
func (md *MCData) Entities() (Entities, error) {
	return md.entities.get(func() (Entities, error) {
		var arr []Entity
		if _, err := md.feature("entities", &arr); err != nil {
			return Entities{}, err
		}
		e := Entities{
			Array:    arr,
			ByID:     make(map[int]*Entity, len(arr)),
			ByName:   make(map[string]*Entity, len(arr)),
			MobsByID: map[int]*Entity{},
			ObjsByID: map[int]*Entity{},
		}
		for i := range arr {
			ent := &e.Array[i]
			e.ByID[ent.ID] = ent
			e.ByName[ent.Name] = ent
			switch ent.Type {
			case "mob":
				e.MobsByID[ent.ID] = ent
			case "object":
				e.ObjsByID[ent.ID] = ent
			}
		}
		return e, nil
	})
}

// Biomes returns the indexed biome collection for this version.
func (md *MCData) Biomes() (Biomes, error) {
	return md.biomes.get(func() (Biomes, error) {
		var arr []Biome
		if _, err := md.feature("biomes", &arr); err != nil {
			return Biomes{}, err
		}
		b := Biomes{Array: arr, ByID: make(map[int]*Biome, len(arr)), ByName: make(map[string]*Biome, len(arr))}
		for i := range arr {
			b.ByID[arr[i].ID] = &b.Array[i]
			b.ByName[arr[i].Name] = &b.Array[i]
		}
		return b, nil
	})
}

// Effects returns the indexed status-effect collection for this version.
func (md *MCData) Effects() (Effects, error) {
	return md.effects.get(func() (Effects, error) {
		var arr []Effect
		if _, err := md.feature("effects", &arr); err != nil {
			return Effects{}, err
		}
		e := Effects{Array: arr, ByID: make(map[int]*Effect, len(arr)), ByName: make(map[string]*Effect, len(arr))}
		for i := range arr {
			e.ByID[arr[i].ID] = &e.Array[i]
			e.ByName[arr[i].Name] = &e.Array[i]
		}
		return e, nil
	})
}

// Enchantments returns the indexed enchantment collection for this version.
func (md *MCData) Enchantments() (Enchantments, error) {
	return md.enchants.get(func() (Enchantments, error) {
		var arr []Enchantment
		if _, err := md.feature("enchantments", &arr); err != nil {
			return Enchantments{}, err
		}
		e := Enchantments{Array: arr, ByID: make(map[int]*Enchantment, len(arr)), ByName: make(map[string]*Enchantment, len(arr))}
		for i := range arr {
			e.ByID[arr[i].ID] = &e.Array[i]
			e.ByName[arr[i].Name] = &e.Array[i]
		}
		return e, nil
	})
}

// Foods returns the indexed food collection for this version.
func (md *MCData) Foods() (Foods, error) {
	return md.foods.get(func() (Foods, error) {
		var arr []Food
		if _, err := md.feature("foods", &arr); err != nil {
			return Foods{}, err
		}
		f := Foods{Array: arr, ByID: make(map[int]*Food, len(arr)), ByName: make(map[string]*Food, len(arr))}
		for i := range arr {
			f.ByID[arr[i].ID] = &f.Array[i]
			f.ByName[arr[i].Name] = &f.Array[i]
		}
		return f, nil
	})
}

// Particles returns the indexed particle collection for this version.
func (md *MCData) Particles() (Particles, error) {
	return md.particles.get(func() (Particles, error) {
		var arr []Particle
		if _, err := md.feature("particles", &arr); err != nil {
			return Particles{}, err
		}
		p := Particles{Array: arr, ByID: make(map[int]*Particle, len(arr)), ByName: make(map[string]*Particle, len(arr))}
		for i := range arr {
			p.ByID[arr[i].ID] = &p.Array[i]
			p.ByName[arr[i].Name] = &p.Array[i]
		}
		return p, nil
	})
}

// Instruments returns the indexed instrument collection for this version.
func (md *MCData) Instruments() (Instruments, error) {
	return md.instruments.get(func() (Instruments, error) {
		var arr []Instrument
		if _, err := md.feature("instruments", &arr); err != nil {
			return Instruments{}, err
		}
		in := Instruments{Array: arr, ByID: make(map[int]*Instrument, len(arr))}
		for i := range arr {
			in.ByID[arr[i].ID] = &in.Array[i]
		}
		return in, nil
	})
}

// Sounds returns the indexed sound collection for this version.
func (md *MCData) Sounds() (Sounds, error) {
	return md.sounds.get(func() (Sounds, error) {
		var arr []Sound
		if _, err := md.feature("sounds", &arr); err != nil {
			return Sounds{}, err
		}
		s := Sounds{Array: arr, ByID: make(map[int]*Sound, len(arr)), ByName: make(map[string]*Sound, len(arr))}
		for i := range arr {
			s.ByID[arr[i].ID] = &s.Array[i]
			s.ByName[arr[i].Name] = &s.Array[i]
		}
		return s, nil
	})
}

// --- Convenience single-item lookups (ignore load errors; return ok=false). ---

// BlockByName returns the block with the given name.
func (md *MCData) BlockByName(name string) (*Block, bool) {
	b, err := md.Blocks()
	if err != nil {
		return nil, false
	}
	v, ok := b.ByName[name]
	return v, ok
}

// BlockByID returns the block with the given numeric id.
func (md *MCData) BlockByID(id int) (*Block, bool) {
	b, err := md.Blocks()
	if err != nil {
		return nil, false
	}
	v, ok := b.ByID[id]
	return v, ok
}

// BlockByStateID returns the block owning the given block state id.
func (md *MCData) BlockByStateID(stateID int) (*Block, bool) {
	b, err := md.Blocks()
	if err != nil {
		return nil, false
	}
	v, ok := b.ByStateID[stateID]
	return v, ok
}

// ItemByName returns the item with the given name.
func (md *MCData) ItemByName(name string) (*Item, bool) {
	it, err := md.Items()
	if err != nil {
		return nil, false
	}
	v, ok := it.ByName[name]
	return v, ok
}

// ItemByID returns the item with the given numeric id.
func (md *MCData) ItemByID(id int) (*Item, bool) {
	it, err := md.Items()
	if err != nil {
		return nil, false
	}
	v, ok := it.ByID[id]
	return v, ok
}
