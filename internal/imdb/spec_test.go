package imdb

import (
	"testing"
	"time"
)

func TestToTraktItemUsesRFC3339Timestamps(t *testing.T) {
	rating := 8.0
	item := Item{
		ID: "tt1234567", Kind: itemTypeMovie, Rating: &rating,
		Created: time.Date(2026, 10, 27, 12, 30, 0, 0, time.FixedZone("local", 3*60*60)),
	}
	got := item.ToTraktItem().Movie
	if got.RatedAt == nil || *got.RatedAt != "2026-10-27T09:30:00Z" {
		t.Fatalf("unexpected rated_at: %v", got.RatedAt)
	}
	if got.WatchedAt == nil || *got.WatchedAt != *got.RatedAt {
		t.Fatalf("watched_at must use the same UTC timestamp: %v", got.WatchedAt)
	}
	if got.Rating == nil || *got.Rating != rating || got.IDMeta.IMDb != item.ID {
		t.Fatal("rating or IMDb identifier changed")
	}
}
