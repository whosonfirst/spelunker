package funcs

import (
	"testing"

	"github.com/whosonfirst/spelunker/v2/http"
)

func TestReplaceAll(t *testing.T) {

	uris := http.DefaultURIs()

	v := http.ReplaceAll(uris.Descendants, "{id}", int64(136251273))

	if v != "/id/136251273/descendants" {
		t.Fatalf("Failed replacement")
	}
}

func TestURIForIdSimple(t *testing.T) {

	uris := http.DefaultURIs()

	v := http.URIForIdSimple(uris.Descendants, int64(136251273))

	if v != "/id/136251273/descendants" {
		t.Fatalf("Failed to derive URI for ID")
	}
}

func TestIsAPlacetype(t *testing.T) {

	tests := map[string]string{
		"country":       "a country",
		"airport":       "an airport",
		"ocean":         "an ocean",
		"region":        "a region",
		"island":        "an island",
		"neighbourhood": "a neighbourhood",
		"custom":        "a custom placetype",
		"empire":        "an empire",
	}

	for pt, expected := range tests {

		is_pt := IsAPlacetype(pt)

		if is_pt != expected {
			t.Fatalf("Failed to derive correct 'is a placetype' string for %s. Expected '%s', got '%s'", pt, expected, is_pt)
		}
	}
}
