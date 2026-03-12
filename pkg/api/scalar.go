package api

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const scalarHTML = `<!doctype html>
<html>
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>E-Raport API Docs</title>
  </head>
  <body>
    <script id="api-reference" data-url="/openapi.yaml"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`

func ScalarDocs(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, scalarHTML)
}

func OpenAPIDoc(c *gin.Context) {
	execPath, _ := os.Executable()
	execDir := filepath.Dir(execPath)

	candidates := []string{
		"./docs/openapi.yaml",
		"../docs/openapi.yaml",
		filepath.Join(execDir, "docs", "openapi.yaml"),
	}

	for _, p := range candidates {
		if b, err := os.ReadFile(p); err == nil {
			c.Data(http.StatusOK, "application/yaml; charset=utf-8", b)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "openapi.yaml not found"})
}
