package server

import (
	"fmt"
	"github.com/whyxn/easynas/backend/pkg/log"
	"github.com/whyxn/easynas/backend/pkg/server/router"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()

	r.Use(router.TokenAuthMiddleware())

	// Setup CORS Config
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.MaxAge = time.Second * 3600
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Content-Type")
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "User-Type")
	corsConfig.ExposeHeaders = append(corsConfig.ExposeHeaders, "Content-Length")
	r.Use(cors.New(corsConfig))

	// Setting API Base Path for HTTP APIs
	httpRouter := r.Group("/")

	// Setting up all Http Routes
	router.AddApiRoutes(httpRouter)

	log.Logger.Infof("Starting Web Server in port %s", "8080")
	err := r.Run(fmt.Sprintf(":%s", "8080")) // listen and serve on 0.0.0.0:PORT
	if err != nil {
		log.Logger.Errorw("Failed to start Web Server", "err", err.Error())
	}
}
