package server

const (
	errorDescription = "An internal error occurred."
	errorCode        = "invalid-request"
)

// ErrorContainer contains a response error msg.
type ErrorContainer struct {
	Result    *interface{} `json:"result"`
	ErrorInfo struct {
		Code        string `json:"code"`
		Description string `json:"description"`
	} `json:"error"`
}

// ResponseError contains a response error msg.
type ResponseError struct {
	*ErrorContainer `json:"error"`
}

// Helper to create an internal server error for API issues.
func createInternalServerError() *ResponseError {
	return createAPIError(errorCode, errorDescription)
}

// This creates an error to return.
func createAPIError(code string, description string) *ResponseError {
	c := &ErrorContainer{
		Result: nil,
		ErrorInfo: struct {
			Code        string `json:"code"`
			Description string `json:"description"`
		}{
			Code:        code,
			Description: description,
		},
	}
	return &ResponseError{c}
}
