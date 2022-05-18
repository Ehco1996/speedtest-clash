package speedtest

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_FetchUserInfo(t *testing.T) {
	c := NewClient(http.DefaultClient)
	user, err := c.FetchUserInfo(context.TODO())
	require.NoError(t, err)
	require.NotEmpty(t, user)
	println(user.String())
}

func TestClient_FetchServerList(t *testing.T) {
	ctx := context.TODO()
	c := NewClient(http.DefaultClient)
	user, err := c.FetchUserInfo(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	println(user.String())

	serverList, err := c.FetchServerList(ctx)
	require.NoError(t, err)
	for idx := range serverList {
		println(serverList[idx].String())
	}
}
