package templatecatalog

import (
	"bytes"
	"cmp"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"slices"

	templatedomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/template"
)

//go:embed templates/*.json
var templateFiles embed.FS

type Catalog struct {
	definitions map[string]templatedomain.Definition
	latest      map[string]string
}

func New() (*Catalog, error) {
	entries, err := templateFiles.ReadDir("templates")
	if err != nil {
		return nil, fmt.Errorf("read embedded release templates: %w", err)
	}
	catalog := &Catalog{
		definitions: make(map[string]templatedomain.Definition, len(entries)),
		latest:      make(map[string]string, len(entries)),
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		content, err := templateFiles.ReadFile("templates/" + entry.Name())
		if err != nil {
			return nil, fmt.Errorf("read release template %q: %w", entry.Name(), err)
		}
		definition, err := decode(content)
		if err != nil {
			return nil, fmt.Errorf("decode release template %q: %w", entry.Name(), err)
		}
		identity := key(definition.Key, definition.Version)
		if _, exists := catalog.definitions[identity]; exists {
			return nil, fmt.Errorf("duplicate release template %q", identity)
		}
		catalog.definitions[identity] = definition
		if current, exists := catalog.latest[definition.Key]; !exists || definition.Version > current {
			catalog.latest[definition.Key] = definition.Version
		}
	}
	if len(catalog.definitions) == 0 {
		return nil, errors.New("no embedded release templates")
	}
	return catalog, nil
}

func (catalog *Catalog) Get(templateKey string, version string) (templatedomain.Definition, error) {
	if catalog == nil {
		return templatedomain.Definition{}, errors.New("release template catalog is not initialized")
	}
	definition, exists := catalog.definitions[key(templateKey, version)]
	if !exists {
		return templatedomain.Definition{}, templatedomain.ErrNotFound
	}
	return clone(definition), nil
}

func (catalog *Catalog) Latest(templateKey string) (templatedomain.Definition, error) {
	if catalog == nil {
		return templatedomain.Definition{}, errors.New("release template catalog is not initialized")
	}
	version, exists := catalog.latest[templateKey]
	if !exists {
		return templatedomain.Definition{}, templatedomain.ErrNotFound
	}
	return catalog.Get(templateKey, version)
}

func decode(content []byte) (templatedomain.Definition, error) {
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	var definition templatedomain.Definition
	if err := decoder.Decode(&definition); err != nil {
		return templatedomain.Definition{}, err
	}
	var trailing json.RawMessage
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return templatedomain.Definition{}, errors.New("template must contain exactly one JSON object")
	}
	if err := definition.Validate(); err != nil {
		return templatedomain.Definition{}, err
	}
	slices.SortStableFunc(definition.Items, func(left, right templatedomain.Item) int {
		return cmp.Compare(left.SortOrder, right.SortOrder)
	})
	return definition, nil
}

func key(templateKey string, version string) string { return templateKey + "@" + version }

func clone(definition templatedomain.Definition) templatedomain.Definition {
	definition.Items = slices.Clone(definition.Items)
	return definition
}
