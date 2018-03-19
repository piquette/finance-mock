package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/piquette/finance-mock/fixture"
)

const invalidRoute = "Unrecognized request URL (%s: %s)."

var pathParameterPattern = regexp.MustCompile(`\{(\w+)\}`)

// StubServer handles incoming HTTP requests and responds to them appropriately
// based off the set of routes that it's been configured with.
type StubServer struct {
	Spec     *fixture.Spec
	Fixtures *fixture.Fixtures
	Verbose  bool
	Routes   map[fixture.HTTPVerb][]stubServerRoute
}

// stubServerRoute is a single route in a StubServer's routing table. It has a
// pattern to match an incoming path and a description of the method that would
// be executed in the event of a match.
type stubServerRoute struct {
	pattern   *regexp.Regexp
	operation *fixture.Operation
}

// HandleRequest handles an HTTP request directed at the API stub.
func (s *StubServer) HandleRequest(w http.ResponseWriter, r *http.Request) {

	start := time.Now()
	fmt.Printf("Request: %v %v\n", r.Method, r.URL.Path)
	w.Header().Set("Request-Id", "req_123")

	// pattern-match a route for the request.
	// -----------------------------------------
	route := s.routeRequest(r)
	if route == nil {
		description := fmt.Sprintf(invalidRoute, r.Method, r.URL.Path)
		apiError := createAPIError(errorCode, description)
		writeResponse(w, r, start, http.StatusNotFound, apiError)
		return
	}

	// Determine if the routing table has an appropriate response.
	// -----------------------------------------
	specResponse, ok := route.operation.Responses["200"]
	if !ok {
		fmt.Printf("Couldn't find 200 response in spec\n")
		writeResponse(w, r, start, http.StatusInternalServerError,
			createInternalServerError())
		return
	}

	// Parse query and build request data.
	// -----------------------------------------
	requestData, err := ParseFormString(r.URL.RawQuery)
	if err != nil {
		fmt.Printf("Couldn't parse url query: %v\n", err)
		writeResponse(w, r, start, http.StatusInternalServerError,
			createInternalServerError())
		return
	}

	// Perform request validation based on route.
	// -----------------------------------------
	if requestData["symbols"] == nil {
		fmt.Printf("Couldn't parse url query: %v\n", err)
		writeResponse(w, r, start, http.StatusInternalServerError,
			createInternalServerError())
		return
	}

	// Build the response data.
	// -----------------------------------------
	resourceID := specResponse.Content["resource"]
	responseData := s.Fixtures.Resources[resourceID]
	if responseData == nil {
		fmt.Printf("Couldn't find resource for id (%v) in spec\n", resourceID)
		writeResponse(w, r, start, http.StatusInternalServerError,
			createInternalServerError())
		return
	}

	// Log reponse data to console.
	// -----------------------------------------
	if s.Verbose {
		responseDataJSON, err := json.MarshalIndent(responseData, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Printf("Response data: %s\n", responseDataJSON)
	}

	// Write response.
	// Done.
	// -----------------------------------------
	writeResponse(w, r, start, http.StatusOK, responseData)
}

// InitRouter maps server routes to possible operations.
func (s *StubServer) InitRouter() error {
	var numEndpoints int
	var numPaths int

	s.Routes = make(map[fixture.HTTPVerb][]stubServerRoute)

	for path, verbs := range s.Spec.Paths {
		numPaths++

		pathPattern := compilePath(path)

		if s.Verbose {
			fmt.Printf("Compiled path: %v\n", pathPattern.String())
		}

		for verb, operation := range verbs {
			numEndpoints++

			route := stubServerRoute{
				pattern:   pathPattern,
				operation: operation,
			}

			// net/http will always give us verbs in uppercase, so build our
			// routing table this way too
			verb = fixture.HTTPVerb(strings.ToUpper(string(verb)))

			s.Routes[verb] = append(s.Routes[verb], route)
		}
	}

	fmt.Printf("Routing to %v path(s) and %v endpoint(s)\n",
		numPaths, numEndpoints)
	return nil
}

func (s *StubServer) routeRequest(r *http.Request) *stubServerRoute {
	verbRoutes := s.Routes[fixture.HTTPVerb(r.Method)]
	for _, route := range verbRoutes {
		if route.pattern.MatchString(r.URL.Path) {
			return &route
		}
	}
	return nil
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

	return regexp.MustCompile(pattern + `\z`)
}

func isCurl(userAgent string) bool {
	return strings.HasPrefix(userAgent, "curl/")
}

func writeResponse(w http.ResponseWriter, r *http.Request, start time.Time, status int, data interface{}) {
	if data == nil {
		data = http.StatusText(status)
	}

	var encodedData []byte
	var err error

	if !isCurl(r.Header.Get("User-Agent")) {
		encodedData, err = json.Marshal(&data)
	} else {
		encodedData, err = json.MarshalIndent(&data, "", "  ")
		encodedData = append(encodedData, '\n')
	}

	if err != nil {
		fmt.Printf("Error serializing response: %v\n", err)
		writeResponse(w, r, start, http.StatusInternalServerError, nil)
		return
	}

	w.WriteHeader(status)
	_, err = w.Write(encodedData)
	if err != nil {
		fmt.Printf("Error writing to client: %v\n", err)
	}
	fmt.Printf("Response: elapsed=%v status=%v\n", time.Now().Sub(start), status)
}
