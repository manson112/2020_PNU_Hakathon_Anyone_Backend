package model

// Request ::
type Request struct {
	Data map[string]interface{} `form:"data" binding:"required"`
}

// Response ::
type Response struct {
	Code         int         `json:"code"`
	Message      string      `json:"message"`
	ResponseData interface{} `json:"response_data"`
}

// Get200Response :: Success
func Get200Response(jsonString interface{}) Response {
	var r Response
	r.Code = 200
	r.Message = "Success"
	r.ResponseData = jsonString
	return r
}

// Get300Response :: Failure, Data Not Satisfied
func Get300Response(jsonString interface{}) Response {
	var r Response
	r.Code = 300
	r.Message = "Failure : Data is not satisfied"
	r.ResponseData = jsonString
	return r
}

// Get400Response :: Failure, Databse Error
func Get400Response(jsonString interface{}) Response {
	var r Response
	r.Code = 400
	r.Message = "Failure : Database Error"
	r.ResponseData = jsonString
	return r
}

// Get500Response :: Failure
func Get500Response(jsonString interface{}) Response {
	var r Response
	r.Code = 500
	r.Message = "Failure"
	r.ResponseData = jsonString
	return r
}
