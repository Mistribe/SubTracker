package endpoints

type HttpErrorResponse struct {
	Message string `json:"message" binding:"required"`
}
