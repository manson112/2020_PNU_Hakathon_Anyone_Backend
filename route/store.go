package route

import (
	"anyone-server/database"
	"log"

	"github.com/gin-gonic/gin"
)

// StoreInfo ::
type StoreInfo struct {
	ID          int    `json:"id"`
	CategoryID  int    `json:"category_id"`
	Category    string `json:"category"`
	StoreName   string `json:"store_name"`
	Tags        string `json:"tags"`
	TotalSeat   int    `json:"total_seat"`
	CurrentSeat int    `json:"current_seat"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Lat         string `json:"lag"`
	Lng         string `json:"lng"`
	Noise       int    `json:"noise"`
	Cleanliness int    `json:"cleanliness"`
	Kindness    int    `json:"kindness"`
	Wifi        int    `json:"wifi"`
}

// GetStoreInfo ::
func GetStoreInfo(c *gin.Context) {
	storeID := c.Param("store_id")

	var storeInfo StoreInfo
	query := "SELECT A.id, A.category_id, C.name category, A.name store_name, " +
		"GROUP_CONCAT(E.name) tags, A.total_seat, A.current_seat, " +
		"A.phone_number, A.address, A.lat, A.lng, " +
		"ROUND(AVG(B.noise)) noise, ROUND(AVG(B.cleanliness)) cleanliness, ROUND(AVG(B.kindness)) kindness, ROUND(AVG(B.wifi)) wifi " +
		"FROM store_info A " +
		"LEFT JOIN (SELECT store_id, " +
		"AVG(noise) noise, " +
		"AVG(cleanliness) cleanliness, " +
		"AVG(kindness) kindness, " +
		"AVG(wifi) wifi " +
		"FROM review " +
		"GROUP BY store_id ) B ON A.id = B.store_id " +
		"LEFT JOIN category C ON A.category_id = C.id " +
		"LEFT JOIN hashtaged_store D ON A.id = D.store_id " +
		"LEFT JOIN hashtags E ON D.hashtag_id = E.id " +
		"WHERE A.id=" + storeID + ";"

	db := database.DB()
	err := db.QueryRow(query).Scan(&storeInfo.ID, &storeInfo.CategoryID,
		&storeInfo.Category, &storeInfo.StoreName,
		&storeInfo.Tags, &storeInfo.TotalSeat,
		&storeInfo.CurrentSeat, &storeInfo.PhoneNumber,
		&storeInfo.Address, &storeInfo.Lat,
		&storeInfo.Lng, &storeInfo.Noise,
		&storeInfo.Cleanliness, &storeInfo.Kindness,
		&storeInfo.Wifi)
	if err != nil {
		log.Println(err)
		log.Println("Cannot exec query")
		return
	}

	c.JSON(200, storeInfo)
	return
}