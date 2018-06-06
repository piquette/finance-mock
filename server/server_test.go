package server

import (
	"testing"

	"github.com/piquette/finance-mock/fixture"
	assert "github.com/stretchr/testify/require"
)

func TestCompilePath(t *testing.T) {
	assert.Equal(t, `\A/v7/quote`,
		compilePath(fixture.Path("/v7/quote")).String())
	assert.Equal(t, `\A/v7/quote/(?P<symbols>[\w-_.]+)`,
		compilePath(fixture.Path("/v7/quote/{symbols}")).String())
}

func TestIsCurl(t *testing.T) {
	testCases := []struct {
		userAgent string
		want      bool
	}{
		{"curl/7.51.0", true},

		// false because it's not something (to my knowledge) that cURL would
		// ever return
		{"curl", false},

		{"Mozilla", false},
		{"", false},
	}
	for _, tc := range testCases {
		t.Run(tc.userAgent, func(t *testing.T) {
			assert.Equal(t, tc.want, isCurl(tc.userAgent))
		})
	}
}
