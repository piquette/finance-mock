package yfin

import "net/http"

const (
	internalErrorDescription = "An internal error occurred."
	internalErrorInfo        = "invalid-request"

	symbolsErrorDescription = "Missing value for the \"symbols\" argument"
	symbolsErrorInfo        = "argument-error"

	chartErrorDescription = "No data found, symbol may be delisted"
	chartErrorInfo        = "Not Found"
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

// ChartResponse contains a chart response msg.
type ChartResponse struct {
	*Response `json:"chart"`
}

// CreateQuote creates a valid quote response.
func CreateQuote(quotes []interface{}) (int, *QuoteResponse) {
	c := &Response{
		Result: quotes,
		Error:  nil,
	}
	return http.StatusOK, &QuoteResponse{c}
}

// CreateChart creates a valid chart response.
func CreateChart(chart interface{}) (int, *ChartResponse) {
	c := &Response{
		Result: []interface{}{chart},
		Error:  nil,
	}
	return http.StatusOK, &ChartResponse{c}
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
