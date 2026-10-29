package mcdata

import "testing"

func TestLoadResolvesVersion(t *testing.T) {
	md, err := Load("1.21.4")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if md.Version.Edition != "pc" {
		t.Errorf("edition = %q, want pc", md.Version.Edition)
	}
	if md.Version.DataFolder == "" {
		t.Error("empty data folder")
	}
}

func TestBlockLookups(t *testing.T) {
	md := MustLoad("1.21.4")
	stone, ok := md.BlockByName("stone")
	if !ok {
		t.Fatal("stone not found")
	}
	if stone.DisplayName != "Stone" {
		t.Errorf("stone.DisplayName = %q", stone.DisplayName)
	}
	again, ok := md.BlockByID(stone.ID)
	if !ok || again.Name != "stone" {
		t.Errorf("round-trip by id failed: %+v", again)
	}
	if _, ok := md.BlockByStateID(stone.DefaultState); !ok {
		t.Errorf("stone default state %d not indexed", stone.DefaultState)
	}
}

func TestItemsAndBiomes(t *testing.T) {
	md := MustLoad("1.21.4")
	if _, ok := md.ItemByName("diamond_sword"); !ok {
		t.Error("diamond_sword missing")
	}
	bi, err := md.Biomes()
	if err != nil {
		t.Fatalf("Biomes: %v", err)
	}
	if len(bi.Array) == 0 {
		t.Error("no biomes")
	}
}

func TestCache(t *testing.T) {
	a := MustLoad("1.21.4")
	b := MustLoad("1.21.4")
	if a != b {
		t.Error("Load did not return cached instance")
	}
}

func TestVersionCompare(t *testing.T) {
	md := MustLoad("1.21.4")
	older, err := md.Version.GT("1.8")
	if err != nil {
		t.Fatalf("GT: %v", err)
	}
	if !older {
		t.Error("1.21.4 should be newer than 1.8")
	}
}

func TestBedrock(t *testing.T) {
	md, err := Load("bedrock_1.21.0")
	if err != nil {
		t.Skipf("bedrock 1.21.0 unavailable: %v", err)
	}
	if md.Version.Edition != "bedrock" {
		t.Errorf("edition = %q, want bedrock", md.Version.Edition)
	}
}
