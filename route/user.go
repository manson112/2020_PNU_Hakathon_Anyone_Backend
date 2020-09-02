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
	ID         string `json:"id"`
	CategoryID string `json:"category_id"`
	Image      string `json:"image"`
	Category   string `json:"category"`
	StoreName  string `json:"store_name"`
	Address    string `json:"address"`
	CreatedAt  string `json:"created_at"`
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

	query := "SELECT distinct B.id, B.category_id, B.image, C.name cateogry, B.name store_name, B.address, A.created_at " +
		"FROM bookmark A " +
		"LEFT JOIN store_info B ON A.store_id = B.id " +
		"LEFT JOIN category C ON B.category_id = C.id " +
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
		err = results.Scan(&bookmark.ID, &bookmark.CategoryID, &bookmark.Image,
			&bookmark.Category, &bookmark.StoreName,
			&bookmark.Address, &bookmark.CreatedAt)
		if err != nil {
			log.Println(err)
			log.Println("Cannot get data")
			c.JSON(400, model.Get400Response(""))
			return
		}
		bookmarks = append(bookmarks, bookmark)
		log.Println(bookmark.Image)
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
		c.JSON(300, model.Get300Response(""))
	}
	log.Println(req.UserID)
	log.Println(req.StoreID)
	log.Println(req.Checked)
	query := ""
	if req.Checked == "true" {
		log.Println("True")
		query = "INSERT INTO bookmark (user_id, store_id, created_at) VALUES (" + req.UserID + ", " + req.StoreID + ", NOW());"
	} else if req.Checked == "false" {
		log.Println("False")
		query = "DELETE FROM bookmark WHERE user_id=" + req.UserID + " and store_id=" + req.StoreID + ";"
	}
	if query == "" {
		c.JSON(400, model.Get400Response(""))
		return
	}
	db := database.DB()
	insert, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, model.Get400Response(""))
		return
	}
	defer insert.Close()
	c.JSON(200, model.Get200Response(""))
}
