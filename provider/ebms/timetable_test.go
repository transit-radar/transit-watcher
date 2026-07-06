package ebms

import (
	"testing"

	ebmstest "codeberg.org/transit-radar/transit-watcher/test/data/apicms.ebms.vn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_ListTimetables(t *testing.T) {
	ts := ebmstest.NewServer()
	t.Cleanup(func() {
		ts.Close()
	})

	c, err := NewClient(WithHTTPClient(ts.Client()), WithDomain(ts.URL))
	require.NoError(t, err)

	routes, err := c.ListTimetables(t.Context(), "10")
	assert.NoError(t, err)
	assert.NotNil(t, routes)
}
