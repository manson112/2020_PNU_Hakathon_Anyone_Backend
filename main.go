package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"

	"anyone-server/database"
	db "anyone-server/database"
	fb "anyone-server/firebase"
	route "anyone-server/route"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func rJson() {
	plan, _ := ioutil.ReadFile("aa.json")
	var data map[string]interface{}
	json.Unmarshal(plan, &data)
	log.Println(len(data["Row"].([]interface{})))
	db := database.DB()
	i := 1
	for _, v := range data["Row"].([]interface{}) {
		m := v.(map[string]interface{})
		u := m["업종명"].(string)
		if u == "음식점" || u == "술집" || u == "카페" {
			cID := "1"
			if u == "음식점" {
				cID = "2"
			}
			if u == "술집" {
				cID = "3"
			}
			name, _ := m["업소명"].(string)
			address, ok := m["소재지_도로명_"].(string)
			if !ok {
				address = ""
			}
			phone, ok := m["소재지전화"].(string)
			if !ok {
				phone = ""
			}
			tSeat := fmt.Sprintf("%d", rand.Intn(30)+8)
			query := "INSERT INTO store_info (category_id, lat, lng, name, total_seat, current_seat, phone_number, address, created_at, updated_at) " +
				"VALUES (" + cID + ", '', '', '" + name + "', " + tSeat + ", 2, '" + phone + "', '" + address + "', NOW(), NOW());"
			_, err := db.Query(query)
			if err != nil {
				log.Println(err)
			}
			log.Println("[" + fmt.Sprintf("%d", i) + "/4120]")
			i++
		}
	}
}

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
	defer db.DB().Close()
	// rJson()

	// Creates a router without any middleware by default
	r := gin.New()
	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())
	r.LoadHTMLGlob("templates/*")
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	router := r.Group("/")
	{
		router.POST("/store/info", route.GetStoreInfo)
		router.POST("/store/home", route.GetStoreHome)
		router.POST("/review/insert", route.PutReview)
		router.POST("/user/bookmark", route.GetBookmarks)
		router.POST("/user/search/history", route.GetSearchHistory)
		router.POST("/store/home/near", route.GetStoreNearLocation)
		router.POST("/store/near", route.GetStoreNearLocation)
		router.GET("/input", route.InputLatLng)
		router.POST("/input", route.Input)

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
