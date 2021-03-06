package rest

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	bookApi "github.com/khanhpdt/bookmark-api/internal/app/rest/book"
	tagApi "github.com/khanhpdt/bookmark-api/internal/app/rest/tag"
)

// Init initializes REST APIs.
func Init() {
	log.Println("Setting up REST APIs...")

	r := gin.Default()

	// - No origin allowed by default
	// - GET,POST, PUT, HEAD methods
	// - Credentials share disabled
	// - Preflight requests cached for 12 hours
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // for development
	r.Use(cors.New(config))

	r.MaxMultipartMemory = 8 << 20 // 8 MB (default is 32 MB)

	setupApis(r)

	err := r.Run(":8081") // listen and serve on 0.0.0.0:8081

	log.Fatalf("error starting gin on port 8081 %s", err)
}

func setupApis(r *gin.Engine) {
	bookApi.Setup(r)
	tagApi.Setup(r)
}
