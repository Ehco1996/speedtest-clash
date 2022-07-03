package speedtest

import (
	"context"
	"fmt"
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

func TestClient_Server_PingTest(t *testing.T) {
	ctx := context.TODO()
	c := NewClient(http.DefaultClient)
	user, err := c.FetchUserInfo(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	println(user.String())

	serverList, err := c.FetchServerList(ctx)
	require.NoError(t, err)
	for idx := range serverList {
		s := serverList[idx]
		println(s.String())
		require.NoError(t, s.GetPingLatency(ctx, c.inner))
	}
}

func TestClient_Server_DownLoad(t *testing.T) {
	ctx := context.TODO()
	c := NewClient(http.DefaultClient)
	user, err := c.FetchUserInfo(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	println(user.String())

	serverList, err := c.FetchServerList(ctx)
	require.NoError(t, err)
	require.Greater(t, len(serverList), 0)
	s := serverList[0]

	require.NoError(t, s.DownLoadTest(ctx, c.GetInnerClient()))
	require.Greater(t, s.DLSpeed, float64(0))
	fmt.Printf("download speed is %.2f mpbs", s.DLSpeed)
}
