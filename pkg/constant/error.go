package constant

// ApplicationError is the struct that defines error on http call (status code except 200 & 201)
type ApplicationError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status"`
	Code       string `json:"code"`
	Success    bool   `json:"success"`
}
