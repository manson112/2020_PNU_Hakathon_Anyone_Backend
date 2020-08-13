package model

type FirebaseCredential struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
}

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
