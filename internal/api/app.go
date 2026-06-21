package api

import (
	"io/fs"
	"net/http"
)

// NewAppHandler assembles the API routes, OpenAPI routes, and optional web UI
// into a single HTTP handler for the service.
func NewAppHandler(apiServer *Server, webRoot fs.FS, openAPISpec []byte) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/v1/", apiServer)
	mux.HandleFunc("/healthz", apiServer.ServeHTTP)

	// /metrics is deliberately unauthenticated so an internal Prometheus
	// scraper can reach it.
	mux.HandleFunc("GET /metrics", apiServer.handleMetrics)

	if len(openAPISpec) > 0 {
		mux.HandleFunc("GET /openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/yaml")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(openAPISpec)
		})
		mux.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(swaggerUIPage))
		})
	}

	if webRoot != nil {
		mux.Handle("/", http.FileServer(http.FS(webRoot)))
	}

	// Wrap everything in the observability middleware so every request is
	// logged in a structured form and counted/timed. apiServer.Mux is passed
	// so the middleware can recover the templated route pattern for labels.
	return observe(mux, apiServer.Mux, apiServer.Logger, apiServer.Metrics)
}

const swaggerUIPage = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Bonus Ledger API Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
  <style>
    body { margin: 0; background: #f5f5f2; }
    .topbar { display: none; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function () {
      window.ui = SwaggerUIBundle({
        url: '/openapi.yaml',
        dom_id: '#swagger-ui'
      });
    };
  </script>
</body>
</html>
`
