package config

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func checkError(t *testing.T, expected interface{}, actual error, msgAndArgs ...interface{}) {
	t.Helper()
	switch expected := expected.(type) {
	case nil:
		require.NoError(t, actual, msgAndArgs...)
	case bool:
		if expected {
			require.Error(t, actual, msgAndArgs...)
		} else {
			require.NoError(t, actual, msgAndArgs...)
		}
	case func(error) bool:
		require.True(t, expected(actual), msgAndArgs...)
	default:
		require.Equal(t, expected, actual, msgAndArgs...)
	}
}

var leadingTabs = regexp.MustCompile("(?m:^\t+)")

func untab(s string) string {
	return leadingTabs.ReplaceAllStringFunc(s, func(match string) string {
		return strings.Repeat("    ", len(match))
	})
}
