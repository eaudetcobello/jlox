package jlox

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
)

func TestIsTruthy(t *testing.T) {

	tests := []struct {
		input    any
		expected bool
	}{
		{
			input:    1.0,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#v", tt.input), func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(isTruthy(tt.input)).To(Equal(tt.expected))
		})
	}
}
