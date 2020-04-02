package blastertest

import (
    "github.com/stretchr/testify/suite"
    "testing"
)

// TestBlaster function is needed to actually run the tests.
func TestBlaster(t *testing.T) {
    suite.Run(t, new(blasterTestSuite))
}

