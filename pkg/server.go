package pkg

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func RunGin(port int) error {
	engine := gin.Default()

	engine.Any("/", func(c *gin.Context) {
		c.String(http.StatusOK, "HELLO")
	})

	engine.GET("/api", func(c *gin.Context) {
		source, ok := c.GetQuery("source")
		if !ok {
			log.Warn().Msgf("source is reuqired")
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "source is required",
			})
			return
		}

		dest := c.DefaultQuery("dest", DefaultDest)
		username := c.DefaultQuery("username", DefaultUsername)
		password := c.DefaultQuery("password", DefaultPassword)

		request, err := NewRequest(source, dest, username, password)
		if err != nil {
			log.Warn().Err(err).Send()
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": err.Error(),
			})
			return
		}

		if err := ImageCopy(c.Request.Context(), request); err != nil {
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, map[string]interface{}{
			"dest": dest,
		})
	})

	return engine.Run(fmt.Sprintf(":%d", port))
}
