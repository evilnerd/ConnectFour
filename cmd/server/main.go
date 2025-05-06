package main

import (
	"connectfour/docs"
	"connectfour/internal/handlers"
	"fmt"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"github.com/swaggo/http-swagger"
	"net/http"
	"os"
)

// @title ConnectFour API
// @version 1.0
// @description API for the ConnectFour online game service
// @termsOfService http://swagger.io/terms/

// @contact.name ConnectFour Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8443
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	log.SetLevel(log.DebugLevel)
	log.Println("ConnectFour Server")

	// Programmatically set swagger info
	docs.SwaggerInfo.Title = "ConnectFour API"
	docs.SwaggerInfo.Description = "API for the ConnectFour online game service"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:" + port()
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r := chi.NewRouter()
	handlers.SetupMiddlewares(r)
	handlers.SetupRoutes(r)

	// Add Swagger UI route
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	port := port()
	log.Printf("Starting on port %s...\n", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		fmt.Printf("Error while running the api: %v", err)
	}
}

func port() string {
	port, ok := os.LookupEnv("CONNECT_FOUR_SERVER_PORT")
	if !ok {
		port = "8443"
	}
	return port
}
