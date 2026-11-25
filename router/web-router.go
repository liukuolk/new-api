package router

import (
	"embed"
	"net/http"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/controller"
	"github.com/QuantumNous/new-api/middleware"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

// WebAssets holds the embedded dashboard frontend assets.
type WebAssets struct {
	BuildFS   embed.FS
	IndexPage []byte
}

func SetWebRouter(router *gin.Engine, assets WebAssets, pluginDispatcher gin.HandlerFunc) {
	frontendFS := common.EmbedFolder(assets.BuildFS, "web/dist")
	compress := gzip.Gzip(gzip.DefaultCompression)

	router.NoRoute(
		pluginDispatcher,
		middleware.RouteTag("web"),
		func(c *gin.Context) {
			// ServeFrontendFiles sends build files already compressed.
			if !frontendFS.Exists("/", c.Request.URL.Path) {
				compress(c)
			}
		},
		middleware.AccessTokenAudit(),
		middleware.GlobalWebRateLimit(frontendFS),
		middleware.Cache(),
		middleware.ServeFrontendFiles(frontendFS),
		func(c *gin.Context) {
			// BasePathHandler has already stripped the configured prefix from
			// URL.Path and RequestURI, so gin's registered group prefixes
			// (/v1, /api, /mj, /assets, …) can be matched directly.
			path := c.Request.URL.Path
			if strings.HasPrefix(path, "/.well-known/") {
				// RFC 8615 well-known URIs (for example Chrome's DevTools
				// workspace probe) must not fall through to the SPA index page.
				c.Status(http.StatusNotFound)
				return
			}
			if strings.HasPrefix(path, "/v1") ||
				strings.HasPrefix(path, "/api") ||
				strings.HasPrefix(path, "/assets") ||
				strings.HasPrefix(path, "/mj") ||
				strings.HasPrefix(path, "/google") ||
				strings.HasPrefix(path, "/v1beta") {
				controller.RelayNotFound(c)
				return
			}
			c.Header("Cache-Control", "no-cache")
			c.Data(http.StatusOK, "text/html; charset=utf-8", assets.IndexPage)
		},
	)
}
