// Package docs menyajikan dokumentasi API (Swagger UI) di /docs.
// openapi.yaml & index.html ikut ter-embed ke binary, jadi tidak tergantung working directory.
package docs

import (
	_ "embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed openapi.yaml
var openapiSpec []byte

//go:embed index.html
var indexHTML []byte

// Register: GET /docs (Swagger UI) dan GET /docs/openapi.yaml (spesifikasi).
func Register(r *gin.Engine) {
	r.GET("/docs", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
	})
	r.GET("/docs/openapi.yaml", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/yaml; charset=utf-8", openapiSpec)
	})
}
