package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/skye-lopez/go-index.prod/pg"
)

func Open() {
	port := os.Getenv("GIN_PORT")
	if port == "" {
		port = "8080"
	}

	devEnv := os.Getenv("DEV_ENV")
	if devEnv == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// NOTE: This will have to be configured later likely
	r.SetTrustedProxies(nil)

	_, err := pg.NewPG()
	if err != nil {
		log.Fatalf("Could not open API, issue connecting to DB\n%s", err)
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.Run(fmt.Sprintf(":%s", port))
}
