//go:build !solution

package httpgauge

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"sort"
	"sync"
)

type Gauge struct {
	mu          sync.Mutex
	endpointMap map[string]int
}

func New() *Gauge {
	return &Gauge{
		endpointMap: make(map[string]int),
	}
}

func (g *Gauge) Snapshot() map[string]int {
	g.mu.Lock()
	defer g.mu.Unlock()

	result := make(map[string]int, len(g.endpointMap))
	for k, v := range g.endpointMap {
		result[k] = v
	}
	return result
}

// ServeHTTP returns accumulated statistics in text format ordered by pattern.
//
// For example:
//
//	/a 10
//	/b 5
//	/c/{id} 7
func (g *Gauge) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.mu.Lock()
	defer g.mu.Unlock()

	var routes []string
	for route := range g.endpointMap {
		routes = append(routes, route)
	}

	sort.Strings(routes)

	for _, route := range routes {
		_, err := fmt.Fprintf(w, "%s %d\n", route, g.endpointMap[route])
		if err != nil {
			return
		}
	}
}

func (g *Gauge) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				routeContext := chi.RouteContext(r.Context())
				if routeContext != nil {
					pattern := routeContext.RoutePattern()
					if pattern != "" {
						g.mu.Lock()
						g.endpointMap[pattern]++
						g.mu.Unlock()
					}
				}
				panic(err)
			}
		}()

		next.ServeHTTP(w, r)

		routeContext := chi.RouteContext(r.Context())
		if routeContext != nil {
			pattern := routeContext.RoutePattern()
			if pattern != "" {
				g.mu.Lock()
				g.endpointMap[pattern]++
				g.mu.Unlock()
			}
		}
	})
}
