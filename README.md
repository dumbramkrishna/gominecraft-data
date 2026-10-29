# gominecraft-data

A fast, type-safe Go accessor for the [PrismarineJS/minecraft-data](https://github.com/PrismarineJS/minecraft-data)
dataset — a Go reimagining of [node-minecraft-data](https://github.com/PrismarineJS/node-minecraft-data).

All Minecraft game data (blocks, items, entities, biomes, effects, recipes,
protocol, …) for every supported Java **and** Bedrock version is embedded in the
binary at build time and exposed through typed Go structs and indexed lookups.

```go
import mcdata "github.com/dumbramkrishna/gominecraft-data"

md, err := mcdata.Load("1.21.4")
if err != nil {
    log.Fatal(err)
}

stone, _ := md.BlockByName("stone")
fmt.Println(stone.ID, stone.DisplayName) // 1 Stone

sword, _ := md.ItemByName("diamond_sword")
fmt.Println(sword.StackSize)             // 1

// Bedrock uses the "bedrock_" (or legacy "pe_") prefix:
be, _ := mcdata.Load("bedrock_1.21.0")
```

## Why it's faster than the original

| | node-minecraft-data | gominecraft-data |
|---|---|---|
| Data access | runtime `require` of JSON | typed structs, parsed once, cached |
| Warm lookup | object property access | **~16 ns, 0 allocations** |
| Shipping | ~hundreds of MB of JSON | **~30 MB** gzip-embedded, no external files |
| Type safety | dynamic | compile-time typed structs |
| Dependencies | — | **zero** (stdlib only) |

The 499 MB raw dataset is gzip-compressed into ~30 MB of embedded assets. A
version's feature file is decompressed and parsed into typed structs the first
time you touch it, then cached — so you pay a one-time ~10 ms per feature and
every subsequent lookup is a map hit.

## API

`Load(version)` returns an `*MCData` for the resolved version. It accepts a bare
version (`"1.21.4"`), a major version (`"1.21"`), or an edition-prefixed string
(`"bedrock_1.21.0"`, `"pe_1.16"`). Results are cached.

Indexed collections (parsed lazily, then cached):

- `Blocks()` → `Array`, `ByID`, `ByName`, `ByStateID`
- `Items()`, `Entities()` (+ `MobsByID`, `ObjsByID`), `Biomes()`, `Effects()`,
  `Enchantments()`, `Foods()`, `Particles()`, `Instruments()`, `Sounds()`

Convenience lookups: `BlockByName`, `BlockByID`, `BlockByStateID`, `ItemByName`,
`ItemByID`.

Version comparison on `md.Version`: `GT`, `GTE`, `LT`, `LTE`, `EQ`, `Compare` —
all ordered by the client `DataVersion`.

## Updating the dataset

The raw data is vendored as a git submodule and regenerated into embedded
assets by the code generator:

```sh
git submodule update --remote minecraft-data
go run ./cmd/gmd-gen   # rewrites internal/mcdb/assets + paths_gen.go
go test ./...
```

## License

MIT. Embeds data from PrismarineJS/minecraft-data under its own license.
Minecraft is a trademark of Mojang Studios; this project is not affiliated with
Mojang or Microsoft.
