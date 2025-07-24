package endpoints

type httpError struct {
	Message string `json:"message" binding:"required"`
}
