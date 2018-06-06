package server

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/piquette/finance-mock/fixture"
	"github.com/piquette/finance-mock/utils"
	"github.com/piquette/finance-mock/yfin"
)

// YFinService is a service that manages yahoo finance requests.
type YFinService struct {
	Service   *fixture.Service
	Resources fixture.Resources
}

// Handle validates a request and returns a response.
func (y *YFinService) Handle(req *http.Request, rte *regexp.Regexp) (statusCode int, responseData interface{}) {

	// Parse query and build request data.
	// -----------------------------------------
	requestData, err := utils.ParseFormString(req.URL.RawQuery)
	if err != nil {
		utils.Log(Verbose, "Couldn't parse url query: %v", err)
		return yfin.CreateInternalServerError()
	}

	// Determine which YFin resource is requested.
	for p, op := range y.Service.Paths {
		// Match path.
		if rte.MatchString(string(p)) {
			// TODO: validate params properly.
			switch op.ResourceID {
			case fixture.YFinQuotes:
				return y.quote(requestData)
			case fixture.YFinChart:
				symbol, err := url.PathUnescape(path.Base(req.URL.Path))
				if err != nil {
					utils.Log(Verbose, "Couldn't parse chart symbol")
					break
				}
				return y.chart(symbol, requestData)
			}
		}
	}

	utils.Log(Verbose, "Couldn't figure out what yfin resource was requested")
	return yfin.CreateInternalServerError()
}

func (y *YFinService) quote(requestData map[string]interface{}) (statusCode int, responseData interface{}) {
	utils.Log(Verbose, "Retrieving quote resource.")

	s := requestData["symbols"]

	if s == nil {
		return yfin.CreateMissingSymbolsError()
	}

	symbolList := strings.Split(s.(string), ",")
	resourceTree := y.Resources[fixture.YFinQuotes].(map[string]interface{})

	quotes := []interface{}{}
	for _, symbol := range symbolList {

		r := resourceTree[symbol]
		if r == nil {
			continue
		}

		quoteMap := r.(map[string]interface{})
		q := quoteMap[strings.ToUpper(string(Market))]
		quotes = append(quotes, q)
	}

	return yfin.CreateQuote(quotes)
}

func (y *YFinService) chart(symbol string, requestData map[string]interface{}) (statusCode int, responseData interface{}) {

	utils.Log(Verbose, "Retrieving chart resource.")

	// TODO: validate properties...

	resourceTree := y.Resources[fixture.YFinChart].(map[string]interface{})
	fmt.Println(symbol)
	r := resourceTree[symbol]
	if r == nil {
		r = resourceTree["error"]
	}
	chartMap := r.(map[string]interface{})

	return yfin.CreateChart(chartMap)
}
