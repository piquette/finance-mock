package server

import (
	"net/http"
	"time"

	"github.com/piquette/finance-mock/utils"
)

const (
	// MarketStatePre genius.
	MarketStatePre MarketState = "pre"

	// MarketStateRegular genius.
	MarketStateRegular MarketState = "regular"

	// MarketStatePost genius.
	MarketStatePost MarketState = "post"
)

// MarketState is a market session.
type MarketState string

// HandleConfigRequest handles an HTTP port directed at the API config stub.
func (s *StubServer) HandleConfigRequest(w http.ResponseWriter, r *http.Request) {

	start := time.Now()
	validStates := []string{string(MarketStatePre), string(MarketStateRegular), string(MarketStatePost)}

	newState := r.PostFormValue("state")
	if newState == "" || !utils.Contains(validStates, newState) {
		utils.Log(Verbose, "Couldn't parse config url")
		s.writeResponse(w, r, start, http.StatusBadRequest, nil)
		return
	}

	// Set market state.
	utils.Log(Verbose, "Changed market state from %v to %v", Market, newState)
	Market = MarketState(newState)

	// Write response.
	s.writeResponse(w, r, start, http.StatusOK, nil)
}
