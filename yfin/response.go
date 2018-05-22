package yfin

import "net/http"

const (
	internalErrorDescription = "An internal error occurred."
	internalErrorInfo        = "invalid-request"

	symbolsErrorDescription = "Missing value for the \"symbols\" argument"
	symbolsErrorInfo        = "argument-error"
)

// Error internal error information structure.
type Error struct {
	Info        string `json:"code"`
	Description string `json:"description"`
}

// Response contains a response msg.
type Response struct {
	Result interface{} `json:"result"`
	Error  *Error      `json:"error"`
}

// ErrorResponse contains a response error msg.
type ErrorResponse struct {
	*Response `json:"error"`
}

// QuoteResponse contains a quote response msg.
type QuoteResponse struct {
	*Response `json:"quoteResponse"`
}

// CreateQuote creates an missing argument error for API issues.
func CreateQuote(quotes []interface{}) (int, *QuoteResponse) {

	c := &Response{
		Result: quotes,
		Error:  nil,
	}
	return http.StatusOK, &QuoteResponse{c}
}

// CreateMissingSymbolsError creates an missing argument error for API issues.
func CreateMissingSymbolsError() (int, *ErrorResponse) {
	return http.StatusBadRequest, createAPIError(symbolsErrorInfo, symbolsErrorDescription)
}

// CreateInternalServerError creates an internal server error for API issues.
func CreateInternalServerError() (int, *ErrorResponse) {
	return http.StatusInternalServerError, createAPIError(internalErrorInfo, internalErrorDescription)
}

// This creates an error to return.
func createAPIError(info string, description string) *ErrorResponse {
	c := &Response{
		Result: nil,
		Error: &Error{
			Info:        info,
			Description: description,
		},
	}
	return &ErrorResponse{c}
}
