// Command example demonstrates basic use of gominecraft-data.
package main

import (
	"fmt"
	"log"

	mcdata "github.com/dumbramkrishna/gominecraft-data"
)

func main() {
	md, err := mcdata.Load("1.21.4")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Loaded %s (protocol %d, dataVersion %d)\n",
		md.Version, md.Version.Version, md.Version.DataVersion)

	if stone, ok := md.BlockByName("stone"); ok {
		fmt.Printf("stone: id=%d display=%q material=%q\n",
			stone.ID, stone.DisplayName, stone.Material)
	}

	blocks, _ := md.Blocks()
	items, _ := md.Items()
	biomes, _ := md.Biomes()
	fmt.Printf("%d blocks, %d items, %d biomes\n",
		len(blocks.Array), len(items.Array), len(biomes.Array))

	newer, _ := md.Version.GT("1.16")
	fmt.Printf("1.21.4 newer than 1.16? %v\n", newer)
}
