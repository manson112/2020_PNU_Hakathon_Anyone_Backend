package route

import (
	"anyone-server/database"
	"anyone-server/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// StoreReq ::
type StoreReq struct {
	StoreID string `form:"store_id"`
}

// StoreHomeReq ::
type StoreHomeReq struct {
	CategoryID string `form:"category_id" binding:"required"`
}

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

// StoreHomeHash ::
type StoreHomeHash struct {
	ID          int    `json:"id"`
	StoreName   string `json:"name"`
	TotalSeat   int    `json:"total_seat"`
	CurrentSeat int    `json:"current_seat"`
	Address     string `json:"address"`
	ImageURL    string `json:"image"`
	Noise       string `json:"noise"`
	Cleanliness string `json:"cleanliness"`
	Kindness    string `json:"kindness"`
	Wifi        string `json:"wifi"`
}

// GetStoreInfo :: [Post] /store/info
func GetStoreInfo(c *gin.Context) {
	var storeReq StoreReq
	err := c.Bind(&storeReq)
	if err != nil {
		log.Fatal(err)
		c.JSON(300, model.Get300Response(""))
	}
	log.Println(storeReq.StoreID)

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
		"WHERE A.id=" + storeReq.StoreID + ";"

	db := database.DB()
	err = db.QueryRow(query).Scan(&storeInfo.ID, &storeInfo.CategoryID,
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
		c.JSON(400, model.Get400Response(""))
		return
	}

	var res []StoreInfo
	res = append(res, storeInfo)

	c.JSON(200, model.Get200Response(res))
	return
}

// GetStoreHome :: [Post] /store/home
func GetStoreHome(c *gin.Context) {
	var storeHomeReq StoreHomeReq
	err := c.Bind(&storeHomeReq)
	if err != nil {
		log.Fatal(err)
		c.JSON(300, model.Get300Response(""))
	}
	log.Println(storeHomeReq.CategoryID)

	query := "SELECT A.id, A.name, A.total_seat, A.current_seat, A.address, A.image, B.noise, B.cleanliness, B.kindness, B.wifi FROM store_info A " +
		"LEFT JOIN (SELECT store_id, " +
		"AVG(noise) noise, " +
		"AVG(cleanliness) cleanliness, " +
		"AVG(kindness) kindness, " +
		"AVG(wifi) wifi " +
		"FROM review " +
		"GROUP BY store_id ) B ON A.id = B.store_id " +
		"WHERE A.category_id=" + storeHomeReq.CategoryID

	db := database.DB()
	results, err := db.Query(query)
	if err != nil {
		log.Println(err)
		log.Println("Cannot exec query")
		c.JSON(400, model.Get400Response(""))
		return
	}

	var hashItems []StoreHomeHash
	for results.Next() {
		var hashItem StoreHomeHash
		err = results.Scan(&hashItem.ID, &hashItem.StoreName,
			&hashItem.TotalSeat, &hashItem.CurrentSeat, &hashItem.Address, &hashItem.ImageURL,
			&hashItem.Noise, &hashItem.Cleanliness, &hashItem.Kindness, &hashItem.Wifi)
		if err != nil {
			log.Println(err)
			log.Println("Cannot get data")
			c.JSON(400, model.Get400Response(""))
			return
		}
		hashItems = append(hashItems, hashItem)
	}
	c.JSON(200, model.Get200Response(hashItems))
}

// StoreNearLocReq ::
type StoreNearLocReq struct {
	CategoryID string `form:"categoryID" binding:"required"`
	Latitude   string `form:"lat" binding:"required"`
	Longitude  string `form:"lng" binding:"required"`
}

// StoreNearLoc ::
type StoreNearLoc struct {
	ID          int    `json:"id"`
	CategoryID  string `json:"category_id"`
	StoreName   string `json:"name"`
	TotalSeat   int    `json:"total_seat"`
	CurrentSeat int    `json:"current_seat"`
	Address     string `json:"address"`
	ImageURL    string `json:"image"`
	Lat         string `json:"lat"`
	Lng         string `json:"lng"`
	Distance    string `json:"distance"`
}

// GetStoreNearLocation ::
func GetStoreNearLocation(c *gin.Context) {
	var storeNearLocReq StoreNearLocReq
	err := c.Bind(&storeNearLocReq)
	if err != nil {
		log.Fatal(err)
		c.JSON(300, model.Get300Response(""))
	}
	log.Println(storeNearLocReq.CategoryID)
	log.Println(storeNearLocReq.Latitude)
	log.Println(storeNearLocReq.Longitude)

	query := "SELECT A.id, A.category_id, A.image as image, A.name, A.address, A.total_seat, A.current_seat, A.lat as latitude, A.lng as longitude, ( 6371000 * acos( cos( radians(" + storeNearLocReq.Latitude + ") ) * cos( radians( A.lat ) ) * cos( radians( A.lng ) - radians(" + storeNearLocReq.Longitude + ") ) + sin( radians(" + storeNearLocReq.Latitude + ") ) * sin(radians(A.lat)) ) ) AS distance FROM store_info A " +
		"HAVING distance < 500 and A.category_id=" + storeNearLocReq.CategoryID + " order by distance limit 20"

	db := database.DB()
	results, err := db.Query(query)
	if err != nil {
		log.Println(err)
		log.Println("Cannot exec query")
		c.JSON(400, model.Get400Response(""))
		return
	}

	var items []StoreNearLoc
	for results.Next() {
		var item StoreNearLoc
		err = results.Scan(&item.ID, &item.CategoryID, &item.ImageURL, &item.StoreName, &item.Address, &item.TotalSeat, &item.CurrentSeat, &item.Lat, &item.Lng, &item.Distance)
		if err != nil {
			log.Println(err)
			log.Println("Cannot get data")
			c.JSON(400, model.Get400Response(""))
			return
		}
		log.Println(item.ID)
		log.Println(item.StoreName)
		log.Println(item.Lat, " ", item.Lng)
		items = append(items, item)
	}
	log.Println(items)
	c.JSON(200, model.Get200Response(items))
}

// // InputLatLng ::
// func InputLatLng(c *gin.Context) {
// 	db := database.DB()
// 	query := "SELECT id, name, address FROM store_info WHERE lat='' and lng='' and (category_id=1 or category_id=2) order by name;"
// 	results, err := db.Query(query)
// 	if err != nil {
// 		log.Println(err)
// 		log.Println("Cannot exec query")
// 		c.HTML(http.StatusOK, "index.tmpl", gin.H{"title": "ERROR"})
// 	}
// 	type R struct {
// 		ID      string `json:"id"`
// 		Name    string `json:"name"`
// 		Address string `json:"address"`
// 	}

// 	var list []R
// 	for results.Next() {
// 		var res R
// 		results.Scan(&res.ID, &res.Name, &res.Address)
// 		list = append(list, res)
// 	}

// 	c.HTML(http.StatusOK, "index.tmpl", gin.H{"title": "이미지", "dataList": list})
// }

// InputLatLng ::
func InputLatLng(c *gin.Context) {
	db := database.DB()
	query := "SELECT id, name, address FROM store_info WHERE image is null and (category_id=1 or category_id=2) order by name;"
	results, err := db.Query(query)
	if err != nil {
		log.Println(err)
		log.Println("Cannot exec query")
		c.HTML(http.StatusOK, "index.tmpl", gin.H{"title": "ERROR"})
	}
	type R struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Address string `json:"address"`
	}

	var list []R
	for results.Next() {
		var res R
		results.Scan(&res.ID, &res.Name, &res.Address)
		list = append(list, res)
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{"title": "이미지", "dataList": list})
}

type ReqLatLng struct {
	ID  string `form:"id" binding:"required"`
	Lat string `form:"lat" binding:"required"`
	Lng string `form:"lng" binding:"required"`
}
type ReqImage struct {
	ID    string `form:"id" binding:"required"`
	Image string `form:"image" binding:"required"`
}

// Input ::
func Input(c *gin.Context) {
	var req ReqImage
	err := c.Bind(&req)
	if err != nil {
		log.Fatal(err)
		c.JSON(300, model.Get300Response(""))
	}
	// query := "UPDATE store_info SET lat=" + req.Lat + ", lng=" + req.Lng + " WHERE id=" + req.ID + ";"
	query := "UPDATE store_info SET image='" + req.Image + "' WHERE id=" + req.ID + ";"
	db := database.DB()
	insert, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, model.Get400Response(""))
	}
	defer insert.Close()

	c.Redirect(http.StatusMovedPermanently, "/input")
}

// Input ::
// func Input(c *gin.Context) {
// 	var req ReqLatLng
// 	err := c.Bind(&req)
// 	if err != nil {
// 		log.Fatal(err)
// 		c.JSON(300, model.Get300Response(""))
// 	}
// 	// query := "UPDATE store_info SET lat=" + req.Lat + ", lng=" + req.Lng + " WHERE id=" + req.ID + ";"
// 	query := "UPDATE store_info SET lat=" + req.Lat + ", lng=" + req.Lng + " WHERE id=" + req.ID + ";"
// 	db := database.DB()
// 	insert, err := db.Query(query)
// 	if err != nil {
// 		log.Fatal(err)
// 		c.JSON(400, model.Get400Response(""))
// 	}
// 	defer insert.Close()
// 	c.JSON(http.StatusOK, "")
// 	// c.Redirect(http.StatusMovedPermanently, "/input")
// }
