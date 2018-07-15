package server

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/piquette/finance-mock/fixture"
	"github.com/piquette/finance-mock/utils"
)

const invalidRoute = "Unrecognized request URL (%s: %s)."

var pathParameterPattern = regexp.MustCompile(`\{(\w+)\}`)

// Version is the mock server version number.
var Version string

// Market is the market sesion.
var Market = MarketStatePost

// Verbose controls log printing.
var Verbose = false

// Handler takes care of requests based on resource types.
type Handler interface {
	Handle(r *http.Request, rte *regexp.Regexp) (statusCode int, responseData interface{})
}

// StubServer handles incoming HTTP requests and responds to them appropriately
// based off the set of routes that it's been configured with.
type StubServer struct {
	Spec       *fixture.Spec
	Fixtures   *fixture.Fixtures
	handlerMap map[*regexp.Regexp]*Handler
}

// HandleRequest handles an HTTP request directed at the API stub.
func (s *StubServer) HandleRequest(w http.ResponseWriter, req *http.Request) {

	start := time.Now()
	utils.Log(Verbose, "Request: %v %v", req.Method, req.URL.Path)
	w.Header().Set("Request-Id", "req_123")

	// Reachability check.
	if req.URL.String() == "/" {
		s.writeResponse(w, req, start, http.StatusOK, nil)
		return
	}

	// pattern-match a handler for the request.
	hndlr, rte := s.routeRequest(req)
	if hndlr == nil {
		utils.Log(Verbose, "Couldn't find handler for url: %v", req.URL.String())
		s.writeResponse(w, req, start, http.StatusNotFound, nil)
		return
	}

	// Build the response data.
	h := *hndlr
	statusCode, responseData := h.Handle(req, rte)

	s.writeResponse(w, req, start, statusCode, responseData)
}

// InitRouter maps server routes to handlers.
func (s *StubServer) InitRouter() error {
	var numServices int
	var numRoutes int

	s.handlerMap = make(map[*regexp.Regexp]*Handler)

	for id, service := range s.Spec.Services {

		var h Handler
		switch id {
		case fixture.ServiceYFin:
			{
				h = &YFinService{
					Service:   service,
					Resources: s.Fixtures.Resources[id],
				}
			}
		default:
			continue
		}

		numServices++

		for path := range service.Paths {
			numRoutes++

			route := compilePath(path)
			utils.Log(Verbose, "Compiled route: %v", route.String())

			// Set the routes and operations.
			s.handlerMap[route] = &h
		}
	}

	utils.Log(Verbose, "Routing to %v service(s) and %v route(s)",
		numServices, numRoutes)

	return nil
}

func (s *StubServer) routeRequest(r *http.Request) (*Handler, *regexp.Regexp) {
	for rte, hdlr := range s.handlerMap {
		if rte.MatchString(r.URL.Path) {
			return hdlr, rte
		}
	}
	return nil, nil
}

func compilePath(path fixture.Path) *regexp.Regexp {
	pattern := `\A`
	parts := strings.Split(string(path), "/")

	for _, part := range parts {
		if part == "" {
			continue
		}

		submatches := pathParameterPattern.FindAllStringSubmatch(part, -1)
		if submatches == nil {
			pattern += `/` + part
		} else {
			pattern += `/(?P<` + submatches[0][1] + `>[\w-_.]+)`
		}
	}

	return regexp.MustCompile(pattern)
}

func isCurl(userAgent string) bool {
	return strings.HasPrefix(userAgent, "curl/")
}

func (s *StubServer) writeResponse(w http.ResponseWriter, r *http.Request, start time.Time, status int, data interface{}) {

	// Sanity check.
	if data == nil {
		data = http.StatusText(status)
	}

	var encodedData []byte
	var err error

	// Marshal response.
	// -----------------
	if !isCurl(r.Header.Get("User-Agent")) {
		encodedData, err = json.Marshal(&data)
	} else {
		encodedData, err = json.MarshalIndent(&data, "", "  ")
		encodedData = append(encodedData, '\n')
	}
	if err != nil {
		utils.Log(Verbose, "Error serializing response: %v", err)
		s.writeResponse(w, r, start, http.StatusInternalServerError, nil)
		return
	}

	// Set headers.
	w.Header().Set("Finance-Mock-Version", Version)
	w.WriteHeader(status)

	// Write response.
	_, err = w.Write(encodedData)
	if err != nil {
		utils.Log(Verbose, "Error writing to client: %v", err)
	}

	utils.Log(Verbose, "Response data: %s", encodedData)

	// Log result.
	utils.Log(Verbose, "Response: elapsed=%v status=%v", time.Now().Sub(start), status)
}
