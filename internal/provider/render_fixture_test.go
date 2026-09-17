package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

type renderFixture struct {
	*httptest.Server
	mu                       sync.Mutex
	images                   map[string]map[string]any
	templates                []map[string]any
	creates, deletes, writes int
}

func newRenderFixture(t *testing.T) *renderFixture {
	t.Helper()
	s := &renderFixture{images: map[string]map[string]any{}}
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()
		id, key, ok := r.BasicAuth()
		if !ok || id != "test-id" || key != "test-key" {
			t.Error("missing API auth")
			w.WriteHeader(401)
			return
		}
		if !strings.Contains(r.UserAgent(), "HCTIPulumi/") {
			t.Error("missing provider user agent")
		}
		write := func(v any) { w.Header().Set("Content-Type", "application/json"); json.NewEncoder(w).Encode(v) }
		decode := func() map[string]any {
			var v map[string]any
			d := json.NewDecoder(r.Body)
			d.UseNumber()
			if e := d.Decode(&v); e != nil {
				t.Error(e)
			}
			return v
		}
		p := r.URL.Path
		switch {
		case strings.HasPrefix(p, "/v1/images/") && r.Method == "GET":
			v, ok := s.images[strings.TrimPrefix(p, "/v1/images/")]
			if !ok {
				w.WriteHeader(404)
				return
			}
			write(v)
		case (p == "/v1/image" || strings.HasPrefix(p, "/v1/image/t-")) && r.Method == "POST":
			v := decode()
			if v["dedupe_duration_s"] != json.Number("0") {
				t.Error("dedupe must be zero")
			}
			s.creates++
			id := fmt.Sprintf("image-%d", s.creates)
			kind := "html_css"
			if _, ok := v["url"]; ok {
				kind = "url"
			}
			if strings.HasPrefix(p, "/v1/image/t-") {
				kind = "templated"
				parts := strings.Split(strings.TrimPrefix(p, "/v1/image/"), "/")
				v["template_id"] = parts[0]
				version := "9007199254740993"
				if len(parts) > 1 {
					version = parts[1]
				}
				v["template_version"] = json.Number(version)
				if _, ok := v["device_scale"]; ok {
					t.Error("unexpected template body fields")
				}
			} else {
				if value, ok := v["device_scale"]; !ok || value != nil {
					t.Error("absent render input must be explicit null")
				}
			}
			format := v["format"]
			delete(v, "format")
			delete(v, "dedupe_duration_s")
			v["image_type"] = kind
			v["id"] = id
			v["created_at"] = "2026-09-16T00:00:00Z"
			v["storage_destination_hcti_storage_disabled"] = v["storage_destination_id"] == "storage-only"
			s.images[id] = v
			route := "image"
			if v["storage_destination_hcti_storage_disabled"] == true {
				route = "store"
			}
			suffix := ""
			if f, ok := format.(string); ok {
				suffix = "." + f
			}
			write(map[string]any{"id": id, "url": s.URL + "/v1/" + route + "/" + id + suffix})
		case strings.HasPrefix(p, "/v1/image/") && r.Method == "DELETE":
			delete(s.images, strings.TrimPrefix(p, "/v1/image/"))
			s.deletes++
			w.WriteHeader(202)
		case (p == "/v1/template" || p == "/v1/template/t-test") && r.Method == "POST":
			v := decode()
			s.writes++
			version := int64(9007199254740990) + int64(s.writes)
			v["id"] = "t-test"
			v["version"] = version
			v["template_type"] = "html_css"
			v["created_at"] = "2026-09-16T00:00:00Z"
			v["updated_at"] = "2026-09-16T00:00:00Z"
			s.templates = append([]map[string]any{v}, s.templates...)
			write(map[string]any{"template_id": "t-test", "template_version": version})
		case p == "/v1/template/t-test" && r.Method == "GET":
			write(map[string]any{"data": s.templates, "pagination": map[string]any{"next_page_start": nil}})
		case p == "/v1/template/t-test" && r.Method == "DELETE":
			s.templates = nil
			w.WriteHeader(204)
		default:
			t.Errorf("unexpected request (rendering must never run): %s %s", r.Method, p)
			w.WriteHeader(500)
		}
	}))
	t.Cleanup(s.Close)
	return s
}
