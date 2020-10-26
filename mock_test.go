package bonusly

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMock(t *testing.T) {
	require.Implements(t, (*Client)(nil), &MockClient{})
}
