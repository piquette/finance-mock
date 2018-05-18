package server

import (
	"net/http"
	"strings"

	"github.com/piquette/finance-mock/fixture"
	"github.com/piquette/finance-mock/utils"
	"github.com/piquette/finance-mock/yfin"
)

const ()

// YFinService is a service that manages yahoo finance requests.
type YFinService struct {
	Service   *fixture.Service
	Resources fixture.Resources
}

// Handle validates a request and returns a response.
func (y *YFinService) Handle(r *http.Request) (statusCode int, responseData interface{}) {

	// Parse query and build request data.
	// -----------------------------------------
	requestData, err := utils.ParseFormString(r.URL.RawQuery)
	if err != nil {
		utils.Log(Verbose, "Couldn't parse url query: %v", err)
		return yfin.CreateMissingSymbolsError()
	}

	// Determine which YFin resource is requested.
	p := strings.Split(r.URL.Path, "/")
	destination := p[len(p)-1]

	switch fixture.ResourceID(destination) {
	case fixture.YFinQuotes:
		return y.quote(requestData)
	}

	return yfin.CreateInternalServerError()
}

func (y *YFinService) quote(requestData map[string]interface{}) (statusCode int, responseData interface{}) {

	s := requestData["symbols"]

	if s == nil || s == "" {
		return yfin.CreateMissingSymbolsError()
	}

	symbolList := strings.Split(s.(string), ",")
	resourceTree := y.Resources[fixture.YFinQuotes].(map[string]interface{})

	var quotes []interface{}
	for _, s := range symbolList {
		resourceTree = resourceTree[s].(map[string]interface{})
		q := resourceTree[strings.ToUpper(string(Market))]
		quotes = append(quotes, q)
	}

	return yfin.CreateQuote(quotes)
}
