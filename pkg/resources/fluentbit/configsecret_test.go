package fluentbit

import (
	"testing"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/stretchr/testify/assert"
)

func TestInvalidFilterGrepConfig(t *testing.T) {
	invalidFilterGrep := &v1beta1.FilterGrep{
		Match:      "*",
		Regex:      []string{"regex", "reg2"},
		Exclude:    []string{"exclude"},
		Logical_Op: "AND",
	}

	_, err := toFluentdFilterGrep(invalidFilterGrep)

	assert.EqualError(t, err, "failed to parse grep filter for fluentbit, Logical_Op is set, it's not posible to set both Regex and Exclude")
}

func TestValidFilterGrepConfig(t *testing.T) {
	filterGrep := &v1beta1.FilterGrep{
		Match:      "*",
		Regex:      []string{"regex1", "regex2"},
		Logical_Op: "AND",
	}

	expectedFluentFilterGrep := &FluentdFilterGrep{
		Match:      "*",
		Regex:      []string{"regex1", "regex2"},
		Logical_Op: "AND",
	}

	parserFluentdfilterGrep, err := toFluentdFilterGrep(filterGrep)

	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, parserFluentdfilterGrep, expectedFluentFilterGrep)
}
