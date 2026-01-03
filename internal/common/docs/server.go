package docs

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
)

// Spec represents an OpenAPI specification with its configuration
type Spec struct {
	Name       string                      // Human-readable name for the API
	Path       string                      // Base path for the spec (e.g., "/auth")
	GetSwagger func() (*openapi3.T, error) // Function to get the OpenAPI spec
}

// Server handles serving OpenAPI documentation
type Server struct {
	specs      []Spec
	specCache  map[string]*openapi3.T
	cacheMutex sync.RWMutex
}

// New creates a new documentation server
func New(specs ...Spec) *Server {
	return &Server{
		specs:     specs,
		specCache: make(map[string]*openapi3.T),
	}
}

// Mount registers all documentation routes on the router
func (s *Server) Mount(r chi.Router) {
	// Combined docs index page
	r.Get("/docs", s.serveDocsIndex)

	// Individual spec endpoints
	for _, spec := range s.specs {
		spec := spec // capture range variable

		// JSON spec endpoint
		r.Get(spec.Path+"/doc/spec.json", s.serveSpec(spec))

		// Swagger UI endpoint
		r.Get(spec.Path+"/doc", s.serveSwaggerUI(spec))
	}
}

func (s *Server) getSpec(spec Spec) (*openapi3.T, error) {
	s.cacheMutex.RLock()
	if cached, ok := s.specCache[spec.Path]; ok {
		s.cacheMutex.RUnlock()
		return cached, nil
	}
	s.cacheMutex.RUnlock()

	swagger, err := spec.GetSwagger()
	if err != nil {
		return nil, err
	}

	s.cacheMutex.Lock()
	s.specCache[spec.Path] = swagger
	s.cacheMutex.Unlock()

	return swagger, nil
}

func (s *Server) serveSpec(spec Spec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		swagger, err := s.getSpec(spec)
		if err != nil {
			http.Error(w, "failed to load spec: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(swagger)
	}
}

func (s *Server) serveSwaggerUI(spec Spec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		html := generateSwaggerHTML(spec.Path+"/doc/spec.json", spec.Name+" API Documentation")
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}
}

func (s *Server) serveDocsIndex(w http.ResponseWriter, r *http.Request) {
	html := s.generateIndexHTML()
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func (s *Server) generateIndexHTML() string {
	specLinks := ""
	for _, spec := range s.specs {
		specLinks += `
			<a href="` + spec.Path + `/doc" class="spec-link">
				<div class="spec-card">
					<h2>` + spec.Name + `</h2>
					<p>View API Documentation</p>
					<span class="arrow">→</span>
				</div>
			</a>`
	}

	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>API Documentation</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
            background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
            min-height: 100vh;
            padding: 40px 20px;
        }
        .container {
            max-width: 800px;
            margin: 0 auto;
        }
        h1 {
            color: #fff;
            text-align: center;
            margin-bottom: 40px;
            font-size: 2.5rem;
            text-shadow: 0 2px 10px rgba(0,0,0,0.3);
        }
        .specs {
            display: grid;
            gap: 20px;
        }
        .spec-link {
            text-decoration: none;
        }
        .spec-card {
            background: rgba(255,255,255,0.1);
            backdrop-filter: blur(10px);
            border-radius: 16px;
            padding: 24px 32px;
            display: flex;
            align-items: center;
            justify-content: space-between;
            transition: all 0.3s ease;
            border: 1px solid rgba(255,255,255,0.1);
        }
        .spec-card:hover {
            background: rgba(255,255,255,0.15);
            transform: translateY(-2px);
            box-shadow: 0 10px 40px rgba(0,0,0,0.2);
        }
        .spec-card h2 {
            color: #fff;
            font-size: 1.5rem;
            margin-bottom: 4px;
        }
        .spec-card p {
            color: rgba(255,255,255,0.7);
            font-size: 0.9rem;
        }
        .arrow {
            color: #4ecdc4;
            font-size: 2rem;
            transition: transform 0.3s ease;
        }
        .spec-card:hover .arrow {
            transform: translateX(5px);
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>📚 API Documentation</h1>
        <div class="specs">` + specLinks + `
        </div>
    </div>
</body>
</html>`
}

func generateSwaggerHTML(specURL, title string) string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta name="description" content="` + title + `" />
    <title>` + title + `</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css" />
    <style>
        body { margin: 0; }
        .swagger-ui .topbar { display: none; }
        .swagger-ui .info { margin: 20px 0; }
    </style>
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js" crossorigin></script>
<script>
    window.onload = () => {
        window.ui = SwaggerUIBundle({
            url: "` + specURL + `",
            dom_id: '#swagger-ui',
            deepLinking: true,
            presets: [
                SwaggerUIBundle.presets.apis,
                SwaggerUIBundle.SwaggerUIStandalonePreset
            ],
            layout: "BaseLayout",
            defaultModelsExpandDepth: -1,
            docExpansion: "list",
            filter: true,
            showExtensions: true,
            showCommonExtensions: true
        });
    };
</script>
</body>
</html>`
}
