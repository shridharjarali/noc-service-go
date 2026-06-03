package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"digit-oss/noc-services/internal/domain"
)

// NewRouter creates a gin.Engine with all NOC routes registered.
// Maps every @PostMapping from NOCController.java:
//
//	POST /v1/noc/_create  → handler.Create
//	POST /v1/noc/_update  → handler.Update
//	POST /v1/noc/_search  → handler.Search
func NewRouter(service domain.NOCService, basePath string) *gin.Engine {
	r := gin.Default()

	handler := &NOCHandler{Service: service}

	// Ensure basePath starts with a slash and does not end with one
	if basePath == "" {
		basePath = ""
	} else {
		if basePath[0] != '/' {
			basePath = "/" + basePath
		}
		// trim trailing slash
		if len(basePath) > 1 && basePath[len(basePath)-1] == '/' {
			basePath = basePath[:len(basePath)-1]
		}
	}

	v1 := r.Group(basePath + "/v1/noc")
	{
		v1.POST("/_create", handler.Create)
		v1.POST("/_update", handler.Update)
		v1.POST("/_search", handler.Search)
	}

	// Health-check endpoint under the same basePath
	r.GET(basePath+"/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	return r
}
