package main

import (
	"log"
	"os"

	db "anyone-server/database"
	fb "anyone-server/firebase"
	route "anyone-server/route"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting...")
	// Load enviroment variables
	err := godotenv.Load(".env.development")
	if err != nil {
		log.Println("Main:: cannot load .env.development file")
	}

	// Firebase
	fb.InitFirebaseAdminSDK()

	db.InitDatabase()

	// Creates a router without any middleware by default
	r := gin.New()
	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	router := r.Group("/")
	{
		router.GET("/store/get/:store_id", route.GetStoreInfo)
		router.POST("/review/insert", route.PutReview)
	}

	// Authorization group
	authorized := r.Group("/auth")
	//authorized.Use(handlers.TokenAuthMiddleware())
	{
		authorized.GET("/ws", func(c *gin.Context) {
			//wsHandler(c.Writer, c.Request, c.ClientIP())
			// sockets.ServeWs(hub, c.Writer, c.Request)
		})
	}
	port := os.Getenv("PORT")
	r.Run(":" + port)
}
