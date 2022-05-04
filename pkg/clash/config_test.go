package clash

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig("../../test/proxies.yaml")
	require.NoError(t, err)
	require.Equal(t, 2, len(cfg.Proxies)) // direct + 1 proxy node
}
