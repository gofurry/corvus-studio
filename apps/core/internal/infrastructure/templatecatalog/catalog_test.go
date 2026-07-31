package templatecatalog_test

import (
	"errors"
	"testing"

	checklistdomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/checklist"
	templatedomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/template"
	"github.com/gofurry/corvus-studio/apps/core/internal/infrastructure/templatecatalog"
)

func TestSteamComingSoonTemplateIsStableAndComplete(t *testing.T) {
	catalog, err := templatecatalog.New()
	if err != nil {
		t.Fatalf("new catalog: %v", err)
	}
	definition, err := catalog.Get("steam-coming-soon", "1.0.0")
	if err != nil {
		t.Fatalf("get template: %v", err)
	}
	if definition.SchemaVersion != 1 || definition.ReviewedAt != "2026-07-31" || len(definition.Items) != 12 {
		t.Fatalf("template metadata = %#v", definition)
	}
	counts := map[checklistdomain.Source]int{}
	for index, item := range definition.Items {
		if item.SortOrder != int64(index+1) {
			t.Fatalf("item %q order = %d, want %d", item.Key, item.SortOrder, index+1)
		}
		counts[item.Source]++
	}
	if counts[checklistdomain.SourcePlatformTemplate] != 8 || counts[checklistdomain.SourceCorvusTemplate] != 4 {
		t.Fatalf("source counts = %#v", counts)
	}

	definition.Items[0].Title = "mutated"
	fresh, err := catalog.Latest("steam-coming-soon")
	if err != nil {
		t.Fatal(err)
	}
	if fresh.Items[0].Title == "mutated" {
		t.Fatal("catalog returned shared mutable template data")
	}
}

func TestCatalogReturnsNotFound(t *testing.T) {
	catalog, err := templatecatalog.New()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := catalog.Get("steam-coming-soon", "9.9.9"); !errors.Is(err, templatedomain.ErrNotFound) {
		t.Fatalf("get unknown version error = %v", err)
	}
}
