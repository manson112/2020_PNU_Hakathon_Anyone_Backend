package route

import (
	"anyone-server/database"
	"anyone-server/model"
	"log"

	"github.com/gin-gonic/gin"
)

// UserReq :: User request data for getting bookmarks
type UserReq struct {
	UserID string `form:"userID" binding:"required"`
}

// Bookmark :: Selected bookmark data from database
type Bookmark struct {
	ID          int    `json:"id"`
	CategoryID  string `json:"category_id"`
	Category    string `json:"category"`
	StoreName   string `json:"name"`
	TotalSeat   int    `json:"total_seat"`
	CurrentSeat int    `json:"current_seat"`
	Address     string `json:"address"`
	ImageURL    string `json:"image"`
	Lat         string `json:"lat"`
	Lng         string `json:"lng"`
	Noise       string `json:"noise"`
	Cleanliness string `json:"cleanliness"`
	Kindness    string `json:"kindness"`
	Wifi        string `json:"wifi"`
	CreatedAt   string `json:"created_at"`
}

// SearchHistory :: Selected search history data from database
type SearchHistory struct {
	SearchQuery string `json:"search_query"`
	CreatedAt   string `json:"created_at"`
}

// GetBookmarks :: [Post] /user/bookmark
func GetBookmarks(c *gin.Context) {
	var userReq UserReq
	err := c.Bind(&userReq)
	if err != nil {
		log.Fatal(err)
		c.JSON(300, model.Get300Response(""))
	}

	query := "SELECT distinct " +
		"B.id, B.category_id, B.image,  " +
		"C.name category, B.name store_name,  " +
		"B.address, B.total_seat, B.current_seat, " +
		"B.lat, B.lng, " +
		"IFNULL(D.noise, 0) noise,  " +
		"IFNULL(D.cleanliness, 0) cleanliness,  " +
		"IFNULL(D.kindness, 0) kindness,  " +
		"IFNULL(D.wifi, 0) wifi,  " +
		"A.created_at " +
		"FROM bookmark A " +
		"LEFT JOIN store_info B ON A.store_id = B.id " +
		"LEFT JOIN category C ON B.category_id = C.id " +
		"LEFT JOIN (SELECT store_id, " +
		"AVG(noise) noise, " +
		"AVG(cleanliness) cleanliness, " +
		"AVG(kindness) kindness, " +
		"AVG(wifi) wifi " +
		"FROM review " +
		"GROUP BY store_id ) D ON A.store_id = D.store_id " +
		"WHERE A.user_id= " + userReq.UserID + " " +
		"ORDER BY A.created_at DESC;"

	db := database.DB()
	results, err := db.Query(query)
	if err != nil {
		log.Println(err)
		log.Println("Cannot exec query")
		c.JSON(400, model.Get400Response(""))
		return
	}
	var bookmarks []Bookmark
	for results.Next() {
		var bookmark Bookmark
		err = results.Scan(&bookmark.ID,
			&bookmark.CategoryID,
			&bookmark.ImageURL,
			&bookmark.Category,
			&bookmark.StoreName,
			&bookmark.Address,
			&bookmark.TotalSeat,
			&bookmark.CurrentSeat,
			&bookmark.Lat,
			&bookmark.Lng,
			&bookmark.Noise,
			&bookmark.Cleanliness,
			&bookmark.Kindness,
			&bookmark.Wifi,
			&bookmark.CreatedAt)
		if err != nil {
			log.Println(err)
			log.Println("Cannot get data")
			c.JSON(400, model.Get400Response(""))
			return
		}
		bookmarks = append(bookmarks, bookmark)
	}
	c.JSON(200, model.Get200Response(bookmarks))
}

// GetSearchHistory :: [Post] /user/search/history
func GetSearchHistory(c *gin.Context) {
	var userReq UserReq
	err := c.Bind(&userReq)
	if err != nil {
		log.Fatal(err)
		c.JSON(300, model.Get300Response(""))
	}

	query := "SELECT search_query, created_at " +
		"FROM search_history " +
		"WHERE user_id=" + userReq.UserID + " " +
		"ORDER BY created_at DESC;"

	db := database.DB()
	results, err := db.Query(query)
	if err != nil {
		log.Println(err)
		log.Println("Cannot exec query")
		c.JSON(400, model.Get400Response(""))
		return
	}
	var shs []SearchHistory
	for results.Next() {
		var sh SearchHistory
		err = results.Scan(&sh.SearchQuery, &sh.CreatedAt)
		if err != nil {
			log.Println(err)
			log.Println("Cannot get data")
			c.JSON(400, model.Get400Response(""))
			return
		}
		shs = append(shs, sh)
	}
	c.JSON(200, model.Get200Response(shs))
}

type PutBookmarkReq struct {
	UserID  string `form:"userID" binding:"required"`
	StoreID string `form:"storeID" binding:"required"`
	Checked string `form:"checked" binding:"required"`
}

// PutBookmark ::
func PutBookmark(c *gin.Context) {
	var req PutBookmarkReq
	err := c.Bind(&req)
	if err != nil {
		log.Fatal(err)
		c.JSON(300, model.Get300Response(nil))
	}
	log.Println(req.UserID)
	log.Println(req.StoreID)
	log.Println(req.Checked)
	query := ""
	if req.Checked == "true" {
		log.Println("True")
		query = "INSERT IGNORE INTO bookmark (user_id, store_id, created_at) VALUES (" + req.UserID + ", " + req.StoreID + ", NOW());"
	} else if req.Checked == "false" {
		log.Println("False")
		query = "DELETE FROM bookmark WHERE user_id=" + req.UserID + " and store_id=" + req.StoreID + ";"
	}
	if query == "" {
		c.JSON(400, model.Get400Response(nil))
		return
	}
	db := database.DB()
	insert, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, model.Get400Response(nil))
		return
	}
	defer insert.Close()
	c.JSON(200, model.Get200Response(nil))
}
