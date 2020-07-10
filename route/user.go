package route

import (
	"anyone-server/database"
	"anyone-server/model"
	"log"

	"github.com/gin-gonic/gin"
)

// UserReq :: User request data for getting bookmarks
type UserReq struct {
	UserID string `form:"id" binding:"required"`
}

// Bookmark :: Selected bookmark data from database
type Bookmark struct {
	ID         string `json:"id"`
	CategoryID string `json:"category_id"`
	Category   string `json:"category"`
	StoreName  string `json:"store_name"`
	Address    string `json:"address"`
	CreatedAt  string `json:"created_at"`
}

// GetBookmarks :: [Post] /user/bookmark
func GetBookmarks(c *gin.Context) {
	var userReq UserReq
	err := c.Bind(&userReq)
	if err != nil {
		log.Fatal(err)
		c.JSON(300, model.Get300Response(""))
	}

	query := "SELECT distinct B.id, B.category_id, C.name cateogry, B.name store_name, B.address, A.created_at " +
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
		err = results.Scan(&bookmark.ID, &bookmark.CategoryID,
			&bookmark.Category, &bookmark.StoreName,
			&bookmark.Address, &bookmark.CreatedAt)
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
