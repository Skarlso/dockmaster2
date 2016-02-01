package main

//Message Basic Message format for greeting
type Message struct {
	Message string `json:"message"`
}

//ErrorResponse A wrapped error response with proper message
type ErrorResponse struct {
	ErrorMessage string `json:"error"`
}
