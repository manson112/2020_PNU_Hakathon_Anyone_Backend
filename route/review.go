package route

import (
	"anyone-server/database"
	"anyone-server/model"
	"log"

	"github.com/gin-gonic/gin"
)

// Review ::
type Review struct {
	StoreID     string `form:"store_id" binding:"required"`
	Noise       string `form:"noise" binding:"required"`
	Cleanliness string `form:"cleanliness" binding:"required"`
	Kindness    string `form:"kindness" binding:"required"`
	Wifi        string `form:"wifi" binding:"required"`
	UserID      string `form:"user_id" binding:"required"`
}

// PutReview :: [Post] /review/insert
func PutReview(c *gin.Context) {
	var review Review
	err := c.Bind(&review)
	if err != nil {
		log.Fatal(err)
		c.JSON(300, model.Get300Response(""))
		return
	}
	log.Println(review)

	query := "INSERT INTO review (store_id, noise, cleanliness, kindness, wifi, user_id, created_at) VALUES (" +
		review.StoreID + ", " + review.Noise + ", " + review.Cleanliness + ", " + review.Kindness + ", " + review.Wifi + ", " + review.UserID + ", NOW());"
	db := database.DB()
	insert, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, model.Get400Response(""))
	}
	defer insert.Close()

	c.JSON(200, model.Get200Response(""))
	return
}
