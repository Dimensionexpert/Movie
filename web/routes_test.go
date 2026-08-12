package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Dimensionexpert/movieLibrary/internal/cache"
)

func TestStaticAssetsAreServedFromEmbeddedFiles(t *testing.T) {
	mux := NewMux(nil, cache.NewCache(1))
	req := httptest.NewRequest(http.MethodGet, "/style.css", nil)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("GET /style.css status = %d, want %d", res.Code, http.StatusOK)
	}
	if contentType := res.Header().Get("Content-Type"); !strings.HasPrefix(contentType, "text/css") {
		t.Errorf("Content-Type = %q, want text/css", contentType)
	}
	if !strings.Contains(res.Body.String(), ".page-shell") {
		t.Error("embedded stylesheet did not contain expected CSS")
	}
}
