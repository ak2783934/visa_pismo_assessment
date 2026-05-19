package response

// ErrorResponse is the standard API error body.
type ErrorResponse struct {
	Error string `json:"error" example:"validation failed"`
}
